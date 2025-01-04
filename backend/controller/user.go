package controller

import (
	"TalkSphere/dao/mysql"
	"TalkSphere/models"
	"TalkSphere/pkg/encrypt"
	"TalkSphere/pkg/jwt"
	"TalkSphere/pkg/oss"
	"TalkSphere/pkg/snowflake"
	"TalkSphere/pkg/upload"
	"TalkSphere/setting"
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterParams 注册请求参数
type RegisterParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
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

// UpdateAvatarParams 更新头像的请求参数
type UpdateAvatarParams struct {
	Avatar string `json:"avatar" binding:"required"`
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
	user = models.User{
		UserID:   userID,
		Username: params.Username,
		Password: encrypt.EncryptPassword(params.Password),
		Email:    params.Email,
		//Avatar:   params.Avatar,
		//Bio:      params.Bio,
	}

	if err := mysql.DB.Create(&user).Error; err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
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
	if encrypt.EncryptPassword(params.Password) != user.Password {
		ResponseError(c, CodeInvalidPassword)
		return
	}

	// 4. 生成Token
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, gin.H{
		"token":    token,
		"user_id":  user.UserID,
		"username": user.Username,
	})
}

// UpdateUserBio 修改用户bio
func UpdateUserBio(c *gin.Context) {
	// 获取参数
	var params UpdateProfileParams
	if err := c.ShouldBindJSON(&params); err != nil { // 注意这里要传入指针
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取当前用户ID
	userID, ok := c.Get("userID")
	if !ok {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 检查用户是否存在
	var user models.User
	result := mysql.DB.Where("user_id = ?", userID.(string)).First(&user)
	if result.Error != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	if result.RowsAffected == 0 {
		ResponseError(c, CodeUserNotExist)
		return
	}

	// 更新bio
	result = mysql.DB.Model(&user).Where("user_id = ?", userID.(string)).Update("bio", params.Bio)
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

	// 获取当前用户ID
	userID, ok := c.Get("userID")
	if !ok {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 先保存到临时目录
	tmpPath, err := upload.SaveAvatar(file, userID.(int64))
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 上传到腾讯云 OSS
	objectKey := fmt.Sprintf("avatars/%d_%d%s", userID.(int64), time.Now().Unix(), path.Ext(file.Filename))

	// 打开临时文件
	f, err := os.Open(tmpPath)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	defer f.Close()

	// 上传到 COS
	_, err = oss.Client.Object.Put(context.Background(), objectKey, f, nil)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 生成访问URL
	avatarURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", setting.Conf.OSSConfig.BucketName, setting.Conf.OSSConfig.Region, objectKey)

	// 更新用户头像URL
	var user models.User
	result := mysql.DB.Model(&user).Where("user_id = ?", userID.(int64)).Update("avatar", avatarURL)
	if result.Error != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 删除临时文件
	os.Remove(tmpPath)

	ResponseSuccess(c, gin.H{
		"avatar_url": avatarURL,
	})
}

// GetUserProfile 获取用户详情
func GetUserProfile(c *gin.Context) {
	// 获取用户ID参数
	userIDStr := c.Param("id") // 从URL参数获取
	var userID interface{}

	if userIDStr == "" {
		// 如果URL中没有ID，则获取当前登录用户的ID
		var ok bool
		userID, ok = c.Get("userID")
		if !ok {
			ResponseError(c, CodeNeedLogin)
			return
		}
	} else {
		userID = userIDStr
	}

	// 查询用户信息
	var user models.User
	result := mysql.DB.Where("user_id = ?", userID).First(&user)
	if result.Error != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	if result.RowsAffected == 0 {
		ResponseError(c, CodeUserNotExist)
		return
	}

	// 构造响应数据
	response := ProfileResponse{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Bio:      user.Bio,
	}

	ResponseSuccess(c, response)
}
