package controller

import (
	"net/http"
	"strconv"

	"github.com/TalkSphere/backend/pkg/mysql"
	"go.uber.org/zap"

	"github.com/TalkSphere/backend/models"

	"github.com/gin-gonic/gin"
)

// CreateBoard 创建板块
func CreateBoard(c *gin.Context) {
	var board models.Board
	if err := c.ShouldBindJSON(&board); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	userID, err := getCurrentUserIDInt64(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	board.CreatorID = userID

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
	zap.L().Info("开始获取所有板块")

	var boards []models.Board
	result := mysql.DB.Find(&boards)

	if result.Error != nil {
		zap.L().Error("查询板块失败",
			zap.Error(result.Error))
		ResponseError(c, CodeServerBusy)
		return
	}

	zap.L().Info("成功获取板块列表",
		zap.Int("total_boards", len(boards)))

	ResponseSuccess(c, boards)
}
