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
