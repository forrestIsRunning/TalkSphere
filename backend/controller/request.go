package controller

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const CtxtUserID = "userID"
const CtxUserName = "userName"

var ErrorUserNotLogin = errors.New("用户未登录")

// getCurrentUserID 获取当前登录的用户 ID
func getCurrentUserID(c *gin.Context) (userID string, err error) {
	uid, ok := c.Get(CtxtUserID)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(string)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

// getPageInfo 获取分页参数
func getPageInfo(c *gin.Context) (int64, int64) {
	pageNumStr := c.Query("page")
	pageSizeStr := c.Query("size")
	var (
		page int64
		size int64
		err  error
	)
	page, err = strconv.ParseInt(pageNumStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return page, size
}

// convertUserIDToInt64 将字符串类型的用户ID转换为int64
func convertUserIDToInt64(userIDStr string) (int64, error) {
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID format: %v", err)
	}
	return userID, nil
}

// getCurrentUserIDInt64 获取当前登录用户的ID（int64类型）
func getCurrentUserIDInt64(c *gin.Context) (int64, error) {
	userIDStr, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("getCurrentUserIDInt64", zap.Error(err))
		return 0, err
	}
	return convertUserIDToInt64(userIDStr)
}
