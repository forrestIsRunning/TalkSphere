package controller

import (
	"TalkSphere/dao/mysql"
	"TalkSphere/models"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LikeRequest 点赞请求
type LikeRequest struct {
	TargetID   int64 `json:"target_id" binding:"required"`
	TargetType int8  `json:"target_type" binding:"required,oneof=1 2"` // 1:post, 2:comment
}

// @Summary 点赞/取消点赞
// @Description 对帖子或评论进行点赞/取消点赞操作
// @Tags 点赞
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param request body LikeRequest true "点赞信息"
// @Success 200 {object} Response{data=map[string]string} "成功，返回liked或unliked"
// @Failure 400 {object} Response "参数错误"
// @Failure 401 {object} Response "未授权"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /likes [post]
func CreateLike(c *gin.Context) {
	var req LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	userIDInterface, exists := c.Get(CtxtUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 正确的类型断言
	userID, ok := userIDInterface.(int64)
	if !ok {
		zap.L().Error("类型断言失败")
		ResponseError(c, CodeServerBusy)
		return
	}

	// 开启事务
	tx := mysql.DB.Begin()

	// 检查是否已点赞
	var like models.Like
	err := tx.Where("user_id = ? AND target_id = ? AND target_type = ?",
		userID, req.TargetID, req.TargetType).First(&like).Error

	if err == nil {
		// 已点赞，取消点赞
		if err := tx.Delete(&like).Error; err != nil {
			tx.Rollback()
			zap.L().Error("取消点赞失败", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}

		// 更新目标点赞数
		if req.TargetType == 1 {
			if err := tx.Model(&models.Post{}).Where("id = ?", req.TargetID).
				UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error; err != nil {
				tx.Rollback()
				zap.L().Error("更新评论点赞数失败", zap.Error(err))
				ResponseError(c, CodeServerBusy)
				return
			}
		} else {
			if err := tx.Model(&models.Comment{}).Where("id = ?", req.TargetID).
				UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error; err != nil {
				tx.Rollback()
				zap.L().Error("更新评论点赞数失败", zap.Error(err))
				ResponseError(c, CodeServerBusy)
				return
			}
		}

		tx.Commit()
		ResponseSuccess(c, gin.H{"status": "unliked"})
		return
	}

	if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		zap.L().Error("查询点赞失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 创建点赞
	like = models.Like{
		UserID:     userID,
		TargetID:   req.TargetID,
		TargetType: req.TargetType,
	}

	if err := tx.Create(&like).Error; err != nil {
		tx.Rollback()
		zap.L().Error("创建点赞失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 更新目标点赞数
	if req.TargetType == 1 {
		if err := tx.Model(&models.Post{}).Where("id = ?", req.TargetID).
			UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
			tx.Rollback()
			zap.L().Error("更新点赞数失败", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}
	} else {
		if err := tx.Model(&models.Comment{}).Where("id = ?", req.TargetID).
			UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
			tx.Rollback()
			zap.L().Error("更新点赞数失败", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}
	}

	tx.Commit()
	ResponseSuccess(c, gin.H{"status": "liked"})
}
