package controller

import (
	"fmt"
	"github.com/TalkSphere/backend/models"
	"github.com/TalkSphere/backend/pkg/encrypt"
	"github.com/TalkSphere/backend/pkg/jwt"
	"github.com/TalkSphere/backend/pkg/mysql"
	"github.com/TalkSphere/backend/pkg/rbac"
	"github.com/TalkSphere/backend/pkg/snowflake"
	"github.com/TalkSphere/backend/pkg/upload"
	"github.com/TalkSphere/backend/setting"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RegisterParams 注册请求参数
type RegisterParams struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatar_url"`
}

// LoginParams 登录请求参数
type LoginParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileParams 更新个人资料的请求参数
type UpdateProfileParams struct {
	Bio string `json:"bio" binding:"max=200"`
}

// ProfileResponse 获取个人资料的响应
type ProfileResponse struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}

// TODO 未来对接邮箱注册

func RegisterHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	var params RegisterParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 检查用户名是否已存在
	var user models.User
	if result := mysql.DB.Where("username = ?", params.Username).First(&user); result.RowsAffected > 0 {
		ResponseError(c, CodeUserExist)
		return
	}

	// 3. 检查邮箱是否已存在
	if result := mysql.DB.Where("email = ?", params.Email).First(&user); result.RowsAffected > 0 {
		ResponseError(c, CodeEmailExist)
		return
	}

	// 4. 创建用户
	userID := snowflake.GenID()
	if params.Bio == "" {
		params.Bio = "no bio"
	}
	if params.AvatarUrl == "" {
		params.AvatarUrl = setting.Conf.DefaultAvatar.AvatarURL
	}
	user = models.User{
		ID:           userID,
		Username:     params.Username,
		PasswordHash: encrypt.EncryptPassword(params.Password),
		Email:        params.Email,
		AvatarURL:    params.AvatarUrl,
		Bio:          params.Bio,
		Status:       1,
		LastLoginAt:  nil,
	}

	if err := mysql.DB.Create(&user).Error; err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 获取 enforcer
	enforcer, exists := c.Get("enforcer")
	if !exists {
		enforcer = rbac.Enforcer
	}
	e := enforcer.(*casbin.Enforcer)

	// 设置用户角色
	var role string
	if params.Username == "admin" {
		role = "admin"
	} else {
		role = "user" // 默认角色为普通用户
	}

	// 为用户添加角色
	_, err := e.AddRoleForUser(fmt.Sprintf("%d", user.ID), role)
	if err != nil {
		zap.L().Error("设置用户角色失败", zap.Error(err))
		// 不要因为设置角色失败就中断注册流程
	} else {
		// 保存策略到数据库
		if err := e.SavePolicy(); err != nil {
			zap.L().Error("保存策略失败", zap.Error(err))
		}
	}

	ResponseSuccess(c, user)
}

func LoginHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	var params LoginParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 查询用户是否存在
	var user models.User
	result := mysql.DB.Where("username = ?", params.Username).First(&user)
	if result.RowsAffected == 0 {
		ResponseError(c, CodeUserNotExist)
		return
	}

	// 3. 验证密码
	if encrypt.EncryptPassword(params.Password) != user.PasswordHash {
		ResponseError(c, CodeInvalidPassword)
		return
	}

	// 4. 生成Token
	token, err := jwt.GenToken(user.ID, user.Username)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, gin.H{
		"token":    "Bearer " + token,
		"userID":   user.ID,
		"username": user.Username,
	})
}

// UpdateUserBio 修改用户bio
func UpdateUserBio(c *gin.Context) {
	// 获取参数
	var params UpdateProfileParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取当前用户ID
	userID, ok := c.Get(CtxtUserID)
	if !ok {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 检查用户是否存在
	var user models.User
	result := mysql.DB.Where("id = ?", userID.(int64)).First(&user)
	if result.Error != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	if result.RowsAffected == 0 {
		ResponseError(c, CodeUserNotExist)
		return
	}

	// 更新bio
	result = mysql.DB.Model(&user).Where("id = ?", userID.(int64)).Update("bio", params.Bio)
	if result.Error != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 返回成功响应
	ResponseSuccess(c, nil)
}

// UpdateUserAvatar 修改用户头像
func UpdateUserAvatar(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	userID, ok := c.Get(CtxtUserID)
	if !ok {
		ResponseError(c, CodeNeedLogin)
		return
	}

	avatarURL, err := upload.SaveImageToOSS(file, "avatar", userID.(int64))
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 更新用户头像URL
	var user models.User
	result := mysql.DB.Model(&user).Where("id = ?", userID.(int64)).Update("avatar_url", avatarURL)
	if result.Error != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, gin.H{
		"avatar_url": avatarURL,
	})
}

// GetUserProfile 获取用户详情
func GetUserProfile(c *gin.Context) {
	userID, ok := c.Get(CtxtUserID)
	if !ok {
		ResponseError(c, CodeNeedLogin)
		return
	}
	userID = userID.(int64)

	// 查询用户信息
	var user models.User
	result := mysql.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	if result.RowsAffected == 0 {
		ResponseError(c, CodeUserNotExist)
		return
	}

	response := ProfileResponse{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.AvatarURL,
		Bio:      user.Bio,
	}

	ResponseSuccess(c, response)
}

func CheckAdminPermission(c *gin.Context) {
	// 从URL参数获取用户ID
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	enforcer, exists := c.Get("enforcer")
	if !exists {
		enforcer = rbac.Enforcer
	}

	e := enforcer.(*casbin.Enforcer)

	// 检查用户是否有 admin 角色
	isAdmin, err := e.HasRoleForUser(fmt.Sprintf("%d", userID), "admin")
	if err != nil {
		zap.L().Error("检查管理员权限失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, gin.H{"is_admin": isAdmin})
}

func GetUserLists(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)

	// 获取搜索关键词
	keyword := c.Query("keyword")

	// 查询用户列表
	var users []models.User
	var total int64

	// 构建查询条件
	query := mysql.DB.Model(&models.User{}).Where("status = ?", 1)

	// 如果有搜索关键词,添加模糊搜索条件
	if keyword != "" {
		query = query.Where(
			"username LIKE ? OR id LIKE ? OR email LIKE ? OR bio LIKE ?",
			"%"+keyword+"%",
			"%"+keyword+"%",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		zap.L().Error("获取用户总数失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 分页查询用户列表
	if err := query.Offset(int((page - 1) * size)).
		Limit(int(size)).
		Find(&users).Error; err != nil {
		zap.L().Error("获取用户列表失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 构造返回数据
	userList := make([]gin.H, 0, len(users))
	for _, user := range users {
		userList = append(userList, gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"avatar":     user.AvatarURL,
			"bio":        user.Bio,
			"created_at": user.CreatedAt,
		})
	}

	ResponseSuccess(c, gin.H{
		"total": total,
		"users": userList,
	})
}

// GetUserRole 获取当前用户的角色
func GetUserRole(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get(CtxtUserID)
	if !exists {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 转换为字符串（因为 casbin 中存储的是字符串）
	userIDStr := fmt.Sprintf("%v", userID)

	// 获取用户角色
	role, err := rbac.GetUserRole(userIDStr)
	if err != nil {
		zap.L().Error("get user role failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, gin.H{
		"user_id": userID,
		"role":    role,
	})
}
