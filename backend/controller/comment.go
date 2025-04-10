package controller

import (
	"TalkSphere/models"
	"TalkSphere/pkg/mysql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	PostID   int64  `json:"post_id" binding:"required" example:"1"`      // 帖子ID
	Content  string `json:"content" binding:"required" example:"这是一条评论"` // 评论内容
	ParentID *int64 `json:"parent_id" example:"0"`                       // 父评论ID，顶级评论为null
}

// @Summary 创建评论
// @Description 创建一条新评论，支持回复其他评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param request body CreateCommentRequest true "评论信息"
// @Success 200 {object} Response{data=models.Comment} "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 401 {object} Response "未授权"
// @Failure 403 {object} Response "无权限"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /comments [post]
func CreateComment(c *gin.Context) {
	var req CreateCommentRequest
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
		ResponseError(c, CodeServerBusy)
		return
	}
	// 开启事务
	tx := mysql.DB.Begin()

	// 检查帖子是否存在且未删除
	var postExists int64
	if err := tx.Model(&models.Post{}).
		Where("id = ? AND status != -1", req.PostID).
		Count(&postExists).Error; err != nil {
		tx.Rollback()
		ResponseError(c, CodeServerBusy)
		return
	}
	if postExists == 0 {
		tx.Rollback()
		ResponseError(c, CodePostNotExist)
		return
	}

	comment := &models.Comment{
		PostID:   req.PostID,
		UserID:   userID,
		Content:  req.Content,
		ParentID: req.ParentID,
		Status:   1,
	}

	// 如果是回复其他评论
	if req.ParentID != nil {
		var parentComment models.Comment
		if err := tx.Where("id = ? AND status = 1", *req.ParentID).
			First(&parentComment).Error; err != nil {
			tx.Rollback()
			if err == gorm.ErrRecordNotFound {
				ResponseError(c, CodeCommentNotExist)
				return
			}
			ResponseError(c, CodeServerBusy)
			return
		}

		// 设置根评论ID
		if parentComment.RootID == nil {
			// 如果父评论是顶级评论，则根评论ID为父评论ID
			comment.RootID = req.ParentID
		} else {
			// 否则继承父评论的根评论ID
			comment.RootID = parentComment.RootID
		}

		// 如果是回复其他评论，更新父评论的回复数
		if err := tx.Model(&models.Comment{}).
			Where("id = ?", *req.ParentID).
			UpdateColumn("reply_count", gorm.Expr("reply_count + ?", 1)).
			Error; err != nil {
			tx.Rollback()
			zap.L().Error("更新父评论回复数失败",
				zap.Error(err),
				zap.Int64("parent_id", *req.ParentID))
			ResponseError(c, CodeServerBusy)
			return
		}
	}

	// 创建评论
	if err := tx.Create(comment).Error; err != nil {
		tx.Rollback()
		ResponseError(c, CodeServerBusy)
		return
	}

	// 更新帖子评论数
	if err := tx.Model(&models.Post{}).
		Where("id = ?", req.PostID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).
		Error; err != nil {
		tx.Rollback()
		ResponseError(c, CodeServerBusy)
		return
	}

	tx.Commit()
	ResponseSuccess(c, comment)
}

// @Summary 获取帖子评论列表
// @Description 获取指定帖子的评论列表，支持分页和排序
// @Tags 评论
// @Accept json
// @Produce json
// @Param post_id path int true "帖子ID"
// @Param page query int false "页码，默认1" default(1)
// @Param size query int false "每页数量，默认10" default(10)
// @Param sort query string false "排序方式：hot(热门)、new(最新)、top(最佳)" default(hot)
// @Success 200 {object} Response{data=[]models.Comment} "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /comments/post/{post_id} [get]
func GetPostComments(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		fmt.Print(postID, err)
		ResponseError(c, CodeInvalidParam)
		return
	}

	sortBy := c.DefaultQuery("sort", "hot") // hot, new, top

	// 1. 首先获取所有顶级评论（分页）
	page, size := getPageInfo(c)
	var rootComments []models.Comment
	var total int64

	query := mysql.DB.Model(&models.Comment{}).
		Where("post_id = ? AND status = 1 AND parent_id IS NULL", postID)

	// 计算总数
	query.Count(&total)

	// 根据排序方式查询
	switch sortBy {
	case "hot":
		query = query.Order("score DESC, created_at DESC")
	case "new":
		query = query.Order("created_at DESC")
	case "top":
		query = query.Order("score DESC")
	}

	// 获取分页后的顶级评论
	if err := query.Offset(int((page - 1) * size)).
		Limit(int(size)).
		Find(&rootComments).Error; err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	if len(rootComments) == 0 {
		ResponseSuccess(c, gin.H{
			"comments": []models.Comment{},
			"total":    0,
		})
		return
	}

	// 2. 获取这些顶级评论的所有子评论
	var rootIDs []int64
	for _, comment := range rootComments {
		rootIDs = append(rootIDs, comment.ID)
	}

	var childComments []models.Comment
	if err := mysql.DB.Where("root_id IN ? AND status = 1", rootIDs).
		Order("created_at ASC").
		Find(&childComments).Error; err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 获取所有相关用户信息
	var userIDs []int64
	userIDMap := make(map[int64]struct{})
	for _, comment := range append(rootComments, childComments...) {
		if _, exists := userIDMap[comment.UserID]; !exists {
			userIDs = append(userIDs, comment.UserID)
			userIDMap[comment.UserID] = struct{}{}
		}
	}

	var users []models.User
	if err := mysql.DB.Where("id IN ?", userIDs).
		Select("id, username, avatar_url").
		Find(&users).Error; err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 构建用户映射
	userMap := make(map[int64]*models.User)
	for i := range users {
		userMap[users[i].ID] = &users[i]
	}

	// 4. 构建评论树
	commentMap := make(map[int64]*models.Comment)

	// 处理顶级评论
	for i := range rootComments {
		rootComments[i].User = userMap[rootComments[i].UserID]
		commentMap[rootComments[i].ID] = &rootComments[i]
	}

	// 处理子评论
	for _, comment := range childComments {
		comment.User = userMap[comment.UserID]
		if parent, exists := commentMap[*comment.ParentID]; exists {
			parent.Children = append(parent.Children, comment)
		}
	}

	ResponseSuccess(c, gin.H{
		"comments": rootComments,
		"total":    total,
	})
}

// @Summary 删除评论
// @Description 删除指定评论及其子评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "评论ID"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 401 {object} Response "未授权"
// @Failure 403 {object} Response "无权限"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /comments/{id} [delete]
func DeleteComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
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
		ResponseError(c, CodeServerBusy)
		return
	}

	tx := mysql.DB.Begin()

	// 获取评论信息
	var comment models.Comment
	if err := tx.First(&comment, commentID).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			ResponseError(c, CodeCommentNotExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 检查权限
	if comment.UserID != userID {
		tx.Rollback()
		ResponseError(c, CodeNoPermision)
		return
	}

	// 软删除评论及其所有子评论
	if err := tx.Model(&models.Comment{}).
		Where("id = ? OR root_id = ?", commentID, commentID).
		Update("status", -1).Error; err != nil {
		tx.Rollback()
		ResponseError(c, CodeServerBusy)
		return
	}

	// 更新帖子评论数
	if err := tx.Model(&models.Post{}).
		Where("id = ?", comment.PostID).
		UpdateColumn("comment_count", gorm.Expr("comment_count - 1")).
		Error; err != nil {
		tx.Rollback()
		ResponseError(c, CodeServerBusy)
		return
	}

	// 更新父评论回复数
	if comment.ParentID != nil {
		if err := tx.Model(&models.Comment{}).
			Where("id = ?", *comment.ParentID).
			UpdateColumn("reply_count", gorm.Expr("reply_count - ?", 1)).
			Error; err != nil {
			tx.Rollback()
			zap.L().Error("更新父评论回复数失败", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}
	}

	tx.Commit()
	ResponseSuccess(c, nil)
}
