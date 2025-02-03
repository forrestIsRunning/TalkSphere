package controller

import (
	"TalkSphere/dao/mysql"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"TalkSphere/models"

	"github.com/gin-gonic/gin"
)

// CreateBoard 创建板块
func CreateBoard(c *gin.Context) {
	var board models.Board
	if err := c.ShouldBindJSON(&board); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	userID, exists := c.Get(CtxtUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	board.CreatorID = userID.(int64)

	if err := mysql.DB.Create(&board).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ResponseSuccess(c, board)
}

// DeleteBoard 删除板块
func DeleteBoard(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	if err := mysql.DB.Delete(&models.Board{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ResponseSuccess(c, gin.H{"message": "删除成功"})
}

// UpdateBoard 修改板块
func UpdateBoard(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	var board models.Board
	if err := c.ShouldBindJSON(&board); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	board.ID = id
	if err := mysql.DB.Model(&board).Updates(board).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ResponseSuccess(c, board)
}

// GetAllBoards 查询所有板块
func GetAllBoards(c *gin.Context) {
	var boards []models.Board
	if err := mysql.DB.Find(&boards).Error; err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	zap.L().Info("boards", zap.Any("boards", boards))
	ResponseSuccess(c, boards)
}
