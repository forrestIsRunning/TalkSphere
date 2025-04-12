package controller

import (
	"fmt"
	"strconv"

	"github.com/TalkSphere/backend/pkg/rbac"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UpdatePermissionsRequest struct {
	UserID      string                   `json:"user_id" binding:"required"`
	Permissions []map[string]interface{} `json:"permissions" binding:"required"`
}

// GetUserPermissions 获取用户权限
func GetUserPermissions(c *gin.Context) {
	// 获取目标用户ID
	targetUserID := c.Param("user_id")
	if targetUserID == "" {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 检查当前用户是否有权限查看目标用户的权限
	currentUserID, exists := c.Get(CtxtUserID)
	if !exists {
		ResponseError(c, CodeNeedLogin)
		return
	}

	currentUserIDStr, err := convertUserIDToString(currentUserID)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取当前用户角色
	currentRole, err := rbac.GetUserRole(currentUserIDStr)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 只有超级管理员和管理员可以查看用户权限
	if currentRole != "super_admin" && currentRole != "admin" {
		ResponseError(c, CodeNoPermision)
		return
	}

	permissions := rbac.GetUserAllPermissions(targetUserID)
	ResponseSuccess(c, permissions)
}

// UpdateUserPermissions 更新用户权限
func UpdateUserPermissions(c *gin.Context) {
	// 获取目标用户ID
	targetUserID := c.Param("user_id")
	if targetUserID == "" {
		ResponseError(c, CodeInvalidParam)
		return
	}

	var req UpdatePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 检查当前用户是否有权限修改目标用户的权限
	currentUserID, exists := c.Get(CtxtUserID)
	if !exists {
		ResponseError(c, CodeNeedLogin)
		return
	}

	currentUserIDStr, err := convertUserIDToString(currentUserID)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取当前用户角色
	currentRole, err := rbac.GetUserRole(currentUserIDStr)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 只有超级管理员可以修改用户权限
	if currentRole != "super_admin" {
		ResponseError(c, CodeNoPermision)
		return
	}

	err = rbac.UpdateUserPermissions(targetUserID, req.Permissions)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}

// CheckPermission 检查权限
func CheckPermission(c *gin.Context) {
	// 获取目标用户ID
	targetUserID := c.Param("user_id")
	if targetUserID == "" {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 检查当前用户是否有权限查看目标用户的权限
	currentUserID, exists := c.Get(CtxtUserID)
	if !exists {
		ResponseError(c, CodeNeedLogin)
		return
	}

	currentUserIDStr, err := convertUserIDToString(currentUserID)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取当前用户角色
	currentRole, err := rbac.GetUserRole(currentUserIDStr)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 只有超级管理员和管理员可以查看用户权限
	if currentRole != "super_admin" && currentRole != "admin" {
		ResponseError(c, CodeNoPermision)
		return
	}

	obj := c.Query("obj")
	act := c.Query("act")

	if obj == "" || act == "" {
		ResponseError(c, CodeInvalidParam)
		return
	}

	hasPermission := rbac.CheckUserPermission(targetUserID, obj, act)
	ResponseSuccess(c, gin.H{
		"has_permission": hasPermission,
	})
}

// GetRolePermissions 获取角色权限
func GetRolePermissions(c *gin.Context) {
	role := c.Query("role")
	if role == "" {
		ResponseError(c, CodeInvalidParam)
		return
	}

	permissions, _ := rbac.GetRolePermissions(role)
	ResponseSuccess(c, permissions)
}

// GetUserRole 获取用户角色
func GetUserRole(c *gin.Context) {
	// 获取目标用户ID
	targetUserID := c.Param("user_id")
	if targetUserID == "" {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 检查当前用户是否有权限查看目标用户的角色
	currentUserID, exists := c.Get(CtxtUserID)
	if !exists {
		ResponseError(c, CodeNeedLogin)
		return
	}

	currentUserIDStr, err := convertUserIDToString(currentUserID)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取当前用户角色
	currentRole, err := rbac.GetUserRole(currentUserIDStr)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 只有超级管理员和管理员可以查看用户角色
	if currentRole != "super_admin" && currentRole != "admin" {
		ResponseError(c, CodeNoPermision)
		return
	}

	// 获取目标用户角色
	role, err := rbac.GetUserRole(targetUserID)
	if err != nil {
		zap.L().Error("get user role failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, gin.H{
		"user_id": targetUserID,
		"role":    role,
	})
}

// convertUserIDToString 将用户ID转换为字符串
func convertUserIDToString(userIDInterface interface{}) (string, error) {
	switch v := userIDInterface.(type) {
	case int64:
		return strconv.FormatInt(v, 10), nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("invalid user ID type: %T", v)
	}
}

// UpdateUserRoleRequest 更新用户角色的请求
type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

// UpdateUserRole 更新用户角色
func UpdateUserRole(c *gin.Context) {
	// 获取目标用户ID
	targetUserID := c.Param("user_id")
	if targetUserID == "" {
		ResponseError(c, CodeInvalidParam)
		return
	}

	zap.L().Info("开始更新用户角色",
		zap.String("target_user_id", targetUserID))

	// 检查当前用户是否有权限修改目标用户的角色
	currentUserID, exists := c.Get(CtxtUserID)
	if !exists {
		zap.L().Error("当前用户未登录")
		ResponseError(c, CodeNeedLogin)
		return
	}

	currentUserIDStr, err := convertUserIDToString(currentUserID)
	if err != nil {
		zap.L().Error("当前用户ID转换失败",
			zap.Any("current_user_id", currentUserID),
			zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取当前用户角色
	currentRole, err := rbac.GetUserRole(currentUserIDStr)
	if err != nil {
		zap.L().Error("获取当前用户角色失败",
			zap.String("current_user_id", currentUserIDStr),
			zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	zap.L().Info("当前用户信息",
		zap.String("current_user_id", currentUserIDStr),
		zap.String("current_role", currentRole))

	// 只有超级管理员可以修改用户角色
	if currentRole != "super_admin" {
		zap.L().Warn("非超级管理员尝试修改用户角色",
			zap.String("current_user_id", currentUserIDStr),
			zap.String("current_role", currentRole))
		ResponseError(c, CodeNoPermision)
		return
	}

	// 解析请求体
	var req UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("请求参数解析失败", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	zap.L().Info("请求的角色更新信息",
		zap.String("target_user_id", targetUserID),
		zap.String("new_role", req.Role))

	// 验证角色是否有效
	validRoles := map[string]bool{
		"user":        true,
		"admin":       true,
		"super_admin": true,
	}
	if !validRoles[req.Role] {
		zap.L().Error("无效的角色",
			zap.String("role", req.Role))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取用户当前角色
	currentUserRole, err := rbac.GetUserRole(targetUserID)
	if err != nil {
		zap.L().Error("获取目标用户当前角色失败",
			zap.String("target_user_id", targetUserID),
			zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	zap.L().Info("目标用户当前角色",
		zap.String("target_user_id", targetUserID),
		zap.String("current_role", currentUserRole))

	// 如果角色没有变化，直接返回成功
	if currentUserRole == req.Role {
		zap.L().Info("用户角色未发生变化，无需更新",
			zap.String("target_user_id", targetUserID),
			zap.String("role", currentUserRole))
		ResponseSuccess(c, nil)
		return
	}

	// 移除用户的当前角色（如果是 guest，不需要移除）
	if currentUserRole != "" && currentUserRole != "guest" {
		zap.L().Info("开始移除用户当前角色",
			zap.String("target_user_id", targetUserID),
			zap.String("current_role", currentUserRole))

		if !rbac.RemoveRole(targetUserID, currentUserRole) {
			zap.L().Error("移除用户当前角色失败",
				zap.String("target_user_id", targetUserID),
				zap.String("current_role", currentUserRole))
			ResponseError(c, CodeServerBusy)
			return
		}

		zap.L().Info("成功移除用户当前角色",
			zap.String("target_user_id", targetUserID),
			zap.String("removed_role", currentUserRole))
	}

	// 添加新角色
	zap.L().Info("开始添加新角色",
		zap.String("target_user_id", targetUserID),
		zap.String("new_role", req.Role))

	if !rbac.AddRole(targetUserID, req.Role) {
		zap.L().Error("添加新角色失败",
			zap.String("target_user_id", targetUserID),
			zap.String("new_role", req.Role))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 验证角色是否设置成功
	newRole, err := rbac.GetUserRole(targetUserID)
	if err != nil {
		zap.L().Error("验证新角色设置失败",
			zap.String("target_user_id", targetUserID),
			zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	if newRole != req.Role {
		zap.L().Error("角色设置验证失败",
			zap.String("target_user_id", targetUserID),
			zap.String("expected_role", req.Role),
			zap.String("actual_role", newRole))
		ResponseError(c, CodeServerBusy)
		return
	}

	zap.L().Info("用户角色更新成功",
		zap.String("target_user_id", targetUserID),
		zap.String("old_role", currentUserRole),
		zap.String("new_role", newRole))

	ResponseSuccess(c, nil)
}
