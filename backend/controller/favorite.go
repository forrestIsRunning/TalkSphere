package controller

import (
	"strconv"

	"github.com/TalkSphere/backend/models"
	"github.com/TalkSphere/backend/pkg/mysql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary 收藏/取消收藏帖子
// @Description 收藏或取消收藏指定帖子
// @Tags 收藏
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param post_id path int true "帖子ID"
// @Success 200 {object} Response{data=map[string]string} "成功，返回favorited或unfavorited"
// @Failure 400 {object} Response "参数错误"
// @Failure 401 {object} Response "未授权"
// @Failure 404 {object} Response "帖子不存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /favorites/post/{post_id} [post]
func CreateFavorite(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	userID, err := getCurrentUserIDInt64(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 开启事务
	tx := mysql.DB.Begin()

	// 检查帖子是否存在
	var post models.Post
	if err := tx.First(&post, postID).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			ResponseError(c, CodePostNotExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 检查是否已收藏
	var favorite models.Favorite
	err = tx.Where("user_id = ? AND post_id = ?", userID, postID).First(&favorite).Error

	if err == nil {
		// 已收藏，取消收藏
		if err := tx.Delete(&favorite).Error; err != nil {
			tx.Rollback()
			ResponseError(c, CodeServerBusy)
			return
		}

		// 更新帖子收藏数
		if err := tx.Model(&post).UpdateColumn("favorite_count",
			gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
			tx.Rollback()
			ResponseError(c, CodeServerBusy)
			return
		}

		tx.Commit()
		ResponseSuccess(c, gin.H{"status": "unfavorited"})
		return
	}

	if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		ResponseError(c, CodeServerBusy)
		return
	}

	// 创建收藏
	favorite = models.Favorite{
		UserID: userID,
		PostID: postID,
	}

	if err := tx.Create(&favorite).Error; err != nil {
		tx.Rollback()
		ResponseError(c, CodeServerBusy)
		return
	}

	// 更新帖子收藏数
	if err := tx.Model(&post).UpdateColumn("favorite_count",
		gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
		tx.Rollback()
		ResponseError(c, CodeServerBusy)
		return
	}

	tx.Commit()
	ResponseSuccess(c, gin.H{"status": "favorited"})
}

// @Summary 获取用户收藏列表
// @Description 获取当前用户的收藏帖子列表
// @Tags 收藏
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int false "页码，默认1" default(1)
// @Param size query int false "每页数量，默认10" default(10)
// @Success 200 {object} Response{data=[]models.Favorite} "成功"
// @Failure 401 {object} Response "未授权"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /favorites [get]
func GetUserFavorites(c *gin.Context) {
	userID, err := getCurrentUserIDInt64(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	page, size := getPageInfo(c)

	var favorites []models.Favorite
	var total int64

	db := mysql.DB.Model(&models.Favorite{}).Where("user_id = ?", userID)
	db.Count(&total)

	if err := db.Preload("Post").Preload("Post.Tags").
		Offset(int((page - 1) * size)).
		Limit(int(size)).
		Find(&favorites).Error; err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, gin.H{
		"favorites": favorites,
		"total":     total,
	})
}
