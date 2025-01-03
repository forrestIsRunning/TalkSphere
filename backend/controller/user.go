package controller

import (
	"TalkSphere/dao/mysql"
	"TalkSphere/models"
	"TalkSphere/pkg/encrypt"
	"TalkSphere/pkg/jwt"
	"TalkSphere/pkg/snowflake"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	var params models.RegisterParams
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

// LoginHandler 处理登录请求
func LoginHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	var params models.LoginParams
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
