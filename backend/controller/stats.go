package controller

import (
	"TalkSphere/models"
	"TalkSphere/pkg/mysql"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetSystemStats(c *gin.Context) {
	// 获取用户总数
	userCount, err := getUserCount()
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 获取帖子总数
	postCount, err := getPostCount()
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 获取板块总数
	boardCount, err := getBoardCount()
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, gin.H{
		"userCount":  userCount,
		"postCount":  postCount,
		"boardCount": boardCount,
	})
}

func getUserCount() (int64, error) {
	var count int64
	result := mysql.DB.Model(&models.User{}).Where("status = ?", 1).Count(&count)
	if result.Error != nil {
		zap.L().Error("get user count failed", zap.Error(result.Error))
		return 0, result.Error
	}
	return count, nil
}

func getPostCount() (int64, error) {
	var count int64
	result := mysql.DB.Model(&models.Post{}).Where("status = ?", 1).Count(&count)
	if result.Error != nil {
		zap.L().Error("get post count failed", zap.Error(result.Error))
		return 0, result.Error
	}
	return count, nil
}

func getBoardCount() (int64, error) {
	var count int64
	result := mysql.DB.Model(&models.Board{}).Count(&count)
	if result.Error != nil {
		zap.L().Error("get board count failed", zap.Error(result.Error))
		return 0, result.Error
	}
	return count, nil
}
