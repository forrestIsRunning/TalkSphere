package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/TalkSphere/backend/models"
	"github.com/TalkSphere/backend/pkg/mysql"
	"github.com/TalkSphere/backend/pkg/upload"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"gorm.io/gorm"
)

// CreatePostRequest 创建帖子请求参数
type CreatePostRequest struct {
	Title    string   `json:"title" binding:"required,min=3,max=100"`
	Content  string   `json:"content" binding:"required,min=10"`
	BoardID  int64    `json:"board_id" binding:"required"`
	Tags     []string `json:"tags" binding:"omitempty,dive,max=20"`
	ImageIDs []int64  `json:"image_ids" binding:"omitempty,dive,min=1"`
}

// UpdatePostRequest 更新帖子请求参数
type UpdatePostRequest struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	BoardID  int64    `json:"board_id"`
	Tags     []string `json:"tags"`
	ImageIDs []int64  `json:"image_ids"`
}

// CreatePost 创建帖子
func CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	userID, err := getCurrentUserIDInt64(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 对内容进行 XSS 清理
	p := bluemonday.UGCPolicy()
	// 允许常用的富文本标签和属性
	p.AllowStandardURLs()
	p.AllowStandardAttributes()
	p.AllowImages()
	p.AllowLists()
	p.AllowTables()
	p.AllowStyles("text-align", "color", "background-color", "font-size", "margin", "padding")
	p.AllowAttrs("class").Globally()

	sanitizedContent := p.Sanitize(req.Content)

	// 生成摘要
	div := strings.NewReader(sanitizedContent)
	doc, err := html.Parse(div)
	if err != nil {
		zap.L().Error("parse html failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	var text string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			text += n.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	// 清理文本并生成摘要
	text = strings.TrimSpace(text)
	runeText := []rune(text)
	const maxExcerptLength = 100 // 减小长度以确保不超过数据库限制

	var excerpt string
	if len(runeText) > maxExcerptLength {
		excerpt = string(runeText[:maxExcerptLength]) + "..."
	} else {
		excerpt = text
	}

	// 检查是否包含图片
	if strings.Contains(sanitizedContent, "<img") {
		// 确保添加图片标记后不超过数据库字段长度限制
		imgText := " [图片]"
		if len([]rune(excerpt))+len([]rune(imgText)) <= 250 { // 留一些余量
			excerpt += imgText
		}
	}

	// 创建帖子
	post := &models.Post{
		Title:    req.Title,
		Content:  sanitizedContent,
		Excerpt:  excerpt,
		BoardID:  &req.BoardID,
		AuthorID: &userID,
	}

	// 开启事务
	tx := mysql.DB.Begin()
	if err := tx.Create(post).Error; err != nil {
		tx.Rollback()
		zap.L().Error("create post failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 处理标签
	if len(req.Tags) > 0 {
		var tags []models.Tag
		for _, tagName := range req.Tags {
			var tag models.Tag
			// 查找或创建标签
			if err := tx.Where("name = ?", tagName).FirstOrCreate(&tag, models.Tag{Name: tagName}).Error; err != nil {
				tx.Rollback()
				zap.L().Error("create post failed", zap.Error(err))
				ResponseError(c, CodeServerBusy)
				return
			}
			tags = append(tags, tag)
		}
		if err := tx.Model(post).Association("Tags").Replace(tags); err != nil {
			tx.Rollback()
			zap.L().Error("create post failed", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}
	}

	// 处理图片
	if len(req.ImageIDs) > 0 {
		// 验证图片所属权并更新关联
		var images []models.PostImage
		if err := tx.Where("id IN ? AND user_id = ? AND status = 1 AND (post_id = 0 OR post_id = ?)",
			req.ImageIDs, userID, post.ID).Find(&images).Error; err != nil {
			tx.Rollback()
			zap.L().Error("find images failed", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}

		// 确保所有图片都存在且属于当前用户
		if len(images) != len(req.ImageIDs) {
			tx.Rollback()
			ResponseError(c, CodeInvalidParam)
			return
		}

		// 更新图片关联
		for i, img := range images {
			if err := tx.Model(&img).Updates(map[string]interface{}{
				"post_id":    post.ID,
				"sort_order": i,
			}).Error; err != nil {
				tx.Rollback()
				zap.L().Error("update image failed", zap.Error(err))
				ResponseError(c, CodeServerBusy)
				return
			}
		}
	}

	tx.Commit()
	ResponseSuccess(c, gin.H{
		"post_id": post.ID,
	})
}

// GetPostDetail 获取帖子详情
func GetPostDetail(c *gin.Context) {
	// 1. 参数验证
	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		zap.L().Error("parse post id failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 查询帖子
	var post models.Post
	result := mysql.DB.Preload("Tags").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", 1).Order("sort_order")
		}).
		Where("status != ?", -1). // 不查询已删除的帖子
		First(&post, postID)

	// 3. 错误处理
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 帖子不存在
			zap.L().Info("post not found", zap.Int64("post_id", postID))
			ResponseError(c, CodePostNotExist)
			return
		}
		// 数据库错误
		zap.L().Error("query post failed",
			zap.Int64("post_id", postID),
			zap.Error(result.Error))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 4. 构建响应数据
	type PostResponse struct {
		ID            int64        `json:"id"`
		Title         string       `json:"title"`
		Content       string       `json:"content"`
		BoardID       *int64       `json:"board_id"`
		AuthorID      *int64       `json:"author_id"`
		ViewCount     int          `json:"view_count"`
		LikeCount     int          `json:"like_count"`
		FavoriteCount int          `json:"favorite_count"`
		CommentCount  int          `json:"comment_count"`
		CreatedAt     time.Time    `json:"created_at"`
		UpdatedAt     time.Time    `json:"updated_at"`
		Tags          []models.Tag `json:"tags"`
		ImageURLs     []string     `json:"image_urls"`
	}
	response := PostResponse{
		ID:            post.ID,
		Title:         post.Title,
		Content:       post.Content,
		BoardID:       post.BoardID,
		AuthorID:      post.AuthorID,
		ViewCount:     post.ViewCount,
		LikeCount:     post.LikeCount,
		FavoriteCount: post.FavoriteCount,
		CommentCount:  post.CommentCount,
		CreatedAt:     post.CreatedAt,
		UpdatedAt:     post.UpdatedAt,
		Tags:          post.Tags,
		ImageURLs:     make([]string, 0, len(post.Images)),
	}

	// 5. 提取图片URL
	for _, img := range post.Images {
		response.ImageURLs = append(response.ImageURLs, img.ImageURL)
	}

	// 6. 异步更新浏览量
	go func() {
		if err := mysql.DB.Model(&post).
			UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).
			Error; err != nil {
			zap.L().Error("update view count failed",
				zap.Int64("post_id", postID),
				zap.Error(err))
		}
	}()

	fmt.Println(response)

	ResponseSuccess(c, response)
}

// DeletePost 删除帖子
func DeletePost(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	userID, err := getCurrentUserIDInt64(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	var post models.Post
	if err := mysql.DB.First(&post, postID).Error; err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 检查是否是帖子作者
	if *post.AuthorID != userID {
		ResponseError(c, CodeNoPermision)
		return
	}

	// 软删除帖子
	if err := mysql.DB.Model(&post).Update("status", -1).Error; err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}

// UpdatePost 更新帖子
func UpdatePost(c *gin.Context) {
	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	userID, err := getCurrentUserIDInt64(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 获取帖子信息
	var post models.Post
	if err := mysql.DB.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ResponseError(c, CodeInvalidParam)
		} else {
			zap.L().Error("get post failed", zap.Error(err))
			ResponseError(c, CodeServerBusy)
		}
		return
	}

	// 检查是否是作者
	if *post.AuthorID != userID {
		ResponseError(c, CodeNoPermision)
		return
	}

	// 对内容进行 XSS 清理
	p := bluemonday.UGCPolicy()
	p.AllowStandardURLs()
	p.AllowStandardAttributes()
	p.AllowImages()
	p.AllowLists()
	p.AllowTables()
	p.AllowStyles("text-align", "color", "background-color", "font-size", "margin", "padding")
	p.AllowAttrs("class").Globally()

	// 更新帖子内容
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		sanitizedContent := p.Sanitize(req.Content)
		updates["content"] = sanitizedContent

		// 生成摘要
		div := strings.NewReader(sanitizedContent)
		doc, err := html.Parse(div)
		if err != nil {
			zap.L().Error("parse html failed", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}

		var text string
		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.Type == html.TextNode {
				text += n.Data
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)

		// 清理文本并生成摘要
		text = strings.TrimSpace(text)
		runeText := []rune(text)
		const maxExcerptLength = 100

		var excerpt string
		if len(runeText) > maxExcerptLength {
			excerpt = string(runeText[:maxExcerptLength]) + "..."
		} else {
			excerpt = text
		}

		// 检查是否包含图片
		if strings.Contains(sanitizedContent, "<img") {
			imgText := " [图片]"
			if len([]rune(excerpt))+len([]rune(imgText)) <= 250 {
				excerpt += imgText
			}
		}

		updates["excerpt"] = excerpt
	}
	if req.BoardID != 0 {
		updates["board_id"] = req.BoardID
	}

	// 开启事务
	tx := mysql.DB.Begin()

	// 更新帖子基本信息
	if err := tx.Model(&post).Updates(updates).Error; err != nil {
		tx.Rollback()
		zap.L().Error("update post failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 更新标签
	if req.Tags != nil {
		// 删除原有标签关联
		if err := tx.Where("post_id = ?", postID).Delete(&models.PostTag{}).Error; err != nil {
			tx.Rollback()
			zap.L().Error("delete post tags failed", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}

		// 添加新标签
		for _, tagName := range req.Tags {
			// 查找或创建标签
			var tag models.Tag
			if err := tx.Where("name = ?", tagName).FirstOrCreate(&tag, models.Tag{Name: tagName}).Error; err != nil {
				tx.Rollback()
				zap.L().Error("create tag failed", zap.Error(err))
				ResponseError(c, CodeServerBusy)
				return
			}

			// 创建帖子-标签关联
			if err := tx.Create(&models.PostTag{PostID: postID, TagID: tag.ID}).Error; err != nil {
				tx.Rollback()
				zap.L().Error("create post tag failed", zap.Error(err))
				ResponseError(c, CodeServerBusy)
				return
			}
		}
	}

	// 更新图片关联
	if req.ImageIDs != nil {
		// 删除原有图片关联
		if err := tx.Where("post_id = ?", postID).Delete(&models.PostImage{}).Error; err != nil {
			tx.Rollback()
			zap.L().Error("delete post images failed", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}

		// 添加新图片关联
		for _, imageID := range req.ImageIDs {
			if err := tx.Create(&models.PostImage{PostID: postID, ID: imageID}).Error; err != nil {
				tx.Rollback()
				zap.L().Error("create post image failed", zap.Error(err))
				ResponseError(c, CodeServerBusy)
				return
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		zap.L().Error("commit transaction failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}

// GetBoardPosts 获取板块帖子列表
func GetBoardPosts(c *gin.Context) {
	// 获取板块ID
	boardID := c.Param("board_id")
	zap.L().Info("开始获取板块帖子列表",
		zap.String("board_id", boardID),
		zap.String("search_query", c.Query("search_query")),
		zap.String("search_type", c.Query("search_type")))

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// 获取搜索参数
	searchQuery := c.Query("search_query")
	searchType := c.Query("search_type")

	// 构建查询条件
	query := mysql.DB.Model(&models.Post{}).Where("posts.board_id = ? AND posts.status != -1", boardID)

	// 根据搜索类型添加搜索条件
	if searchQuery != "" {
		switch searchType {
		case "username":
			// 关联用户表进行用户名搜索
			query = query.Joins("JOIN users ON posts.author_id = users.id").
				Where("users.username LIKE ?", "%"+searchQuery+"%")
		case "content":
			// 搜索帖子内容
			query = query.Where("posts.content LIKE ?", "%"+searchQuery+"%")
		case "all":
			// 同时搜索标题和内容
			query = query.Where("(posts.title LIKE ? OR posts.content LIKE ?)",
				"%"+searchQuery+"%", "%"+searchQuery+"%")
		}
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		zap.L().Error("获取帖子总数失败",
			zap.Error(err),
			zap.String("board_id", boardID),
			zap.String("search_query", searchQuery),
			zap.String("search_type", searchType))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 分页查询帖子
	var posts []models.Post
	if err := query.
		Preload("Author"). // 预加载作者信息
		Preload("Images"). // 预加载图片信息
		Offset((page - 1) * size).
		Limit(size).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		zap.L().Error("查询帖子列表失败",
			zap.Error(err),
			zap.String("board_id", boardID),
			zap.String("search_query", searchQuery),
			zap.String("search_type", searchType),
			zap.Int("page", page),
			zap.Int("size", size))
		ResponseError(c, CodeServerBusy)
		return
	}

	zap.L().Info("成功查询帖子列表",
		zap.Int("post_count", len(posts)),
		zap.Int64("total", total))

	// 构建响应数据
	var postList []map[string]interface{}
	for _, post := range posts {
		postData := map[string]interface{}{
			"id":             post.ID,
			"title":          post.Title,
			"content":        post.Content,
			"author_id":      post.AuthorID,
			"view_count":     post.ViewCount,
			"like_count":     post.LikeCount,
			"comment_count":  post.CommentCount,
			"created_at":     post.CreatedAt,
			"updated_at":     post.UpdatedAt,
			"favorite_count": post.FavoriteCount,
		}

		// 处理作者信息
		if post.Author != nil {
			postData["author"] = map[string]interface{}{
				"id":         post.Author.ID,
				"username":   post.Author.Username,
				"avatar_url": post.Author.AvatarURL,
			}
		}

		// 处理图片
		var images []map[string]interface{}
		if len(post.Images) > 0 {
			for _, img := range post.Images {
				images = append(images, map[string]interface{}{
					"id":         img.ID,
					"url":        img.ImageURL,
					"post_id":    img.PostID,
					"created_at": img.CreatedAt,
				})
			}
		}
		postData["images"] = images

		postList = append(postList, postData)
	}

	ResponseSuccess(c, gin.H{
		"posts":        postList,
		"total":        total,
		"current_page": page,
		"page_size":    size,
	})
}

// UploadPostImage 上传帖子图片
func UploadPostImage(c *gin.Context) {
	zap.L().Info("开始处理图片上传请求")

	// 打印所有接收到的表单字段名
	form, _ := c.MultipartForm()
	if form != nil {
		zap.L().Info("收到的表单字段",
			zap.Any("form_fields", form.File))
	}

	file, err := c.FormFile("image")
	if err != nil {
		zap.L().Error("获取上传文件失败",
			zap.Error(err),
			zap.String("error_type", "form_file_error"),
			zap.String("expected_field", "image"))
		ResponseError(c, CodeInvalidParam)
		return
	}
	zap.L().Info("成功获取上传文件",
		zap.String("filename", file.Filename),
		zap.Int64("size", file.Size))

	userID, err := getCurrentUserIDInt64(c)
	if err != nil {
		zap.L().Error("未找到用户ID")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	zap.L().Info("获取到用户ID", zap.Int64("user_id", userID))

	imageURL, err := upload.SaveImageToOSS(file, "post_images", userID)
	if err != nil {
		zap.L().Error("保存图片到OSS失败",
			zap.Error(err),
			zap.String("filename", file.Filename),
			zap.Int64("user_id", userID))
		if err.Error() == "文件大小超过限制" || err.Error() == "不支持的文件类型" {
			ResponseError(c, CodeInvalidParam)
		} else {
			ResponseError(c, CodeServerBusy)
		}
		return
	}
	zap.L().Info("图片成功上传到OSS",
		zap.String("image_url", imageURL))

	postImage := &models.PostImage{
		UserID:    userID,
		ImageURL:  imageURL,
		Status:    1,
		SortOrder: 0,
	}

	if err := mysql.DB.Create(postImage).Error; err != nil {
		zap.L().Error("保存图片记录到数据库失败",
			zap.Error(err),
			zap.Any("post_image", postImage))
		ResponseError(c, CodeServerBusy)
		return
	}
	zap.L().Info("图片记录成功保存到数据库",
		zap.Int64("image_id", postImage.ID),
		zap.String("image_url", imageURL))

	ResponseSuccess(c, gin.H{
		"image_id":  postImage.ID,
		"image_url": imageURL,
	})
}

// GetUserPosts 获取用户的帖子列表
func GetUserPosts(c *gin.Context) {
	userID, err := getCurrentUserIDInt64(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	zap.L().Info("开始获取用户帖子列表")

	// 获取分页参数
	page, size := getPageInfo(c)
	zap.L().Info("获取分页参数",
		zap.Int64("page", page),
		zap.Int64("size", size))

	var posts []models.Post
	var total int64

	// 构建查询
	db := mysql.DB.Model(&models.Post{}).Where("author_id = ? AND status != -1", userID)

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		zap.L().Error("获取帖子总数失败",
			zap.Error(err),
			zap.Int64("user_id", userID))
		ResponseError(c, CodeServerBusy)
		return
	}
	zap.L().Info("获取到帖子总数", zap.Int64("total", total))

	// 查询帖子数据
	if err := db.Preload("Author").
		Preload("Tags").
		Preload("Images").
		Order("created_at DESC").
		Offset(int((page - 1) * size)).
		Limit(int(size)).
		Find(&posts).Error; err != nil {
		zap.L().Error("查询用户帖子失败",
			zap.Error(err),
			zap.Int64("user_id", userID),
			zap.Int64("page", page),
			zap.Int64("size", size))
		ResponseError(c, CodeServerBusy)
		return
	}

	zap.L().Info("成功查询用户帖子",
		zap.Int("post_count", len(posts)),
		zap.Int64("user_id", userID))

	// 记录每个帖子的基本信息
	for i, post := range posts {
		zap.L().Debug("帖子详情",
			zap.Int("index", i),
			zap.Int64("post_id", post.ID),
			zap.String("title", post.Title),
			zap.Int64("author_id", *post.AuthorID),
			zap.Int("tag_count", len(post.Tags)),
			zap.Int("image_count", len(post.Images)))
	}

	ResponseSuccess(c, gin.H{
		"posts":        posts,
		"total":        total,
		"current_page": page,
		"page_size":    size,
	})
}

// GetUserLikedPosts 获取用户点赞的帖子
func GetUserLikedPosts(c *gin.Context) {
	userID, err := getCurrentUserIDInt64(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	zap.L().Info("开始获取用户点赞的帖子")

	// 获取分页参数
	page, size := getPageInfo(c)
	zap.L().Info("获取分页参数",
		zap.Int64("page", page),
		zap.Int64("size", size))

	// 获取帖子和总数
	posts, total, err := getUserLikedPosts(userID, page, size)
	if err != nil {
		zap.L().Error("获取用户点赞帖子失败",
			zap.Error(err),
			zap.Int64("user_id", userID),
			zap.Int64("page", page),
			zap.Int64("size", size))
		ResponseError(c, CodeServerBusy)
		return
	}
	zap.L().Info("成功获取用户点赞帖子",
		zap.Int64("total", total),
		zap.Int("post_count", len(posts)))

	ResponseSuccess(c, gin.H{
		"posts":        posts,
		"total":        total,
		"current_page": page,
		"page_size":    size,
	})
}

func getUserLikedPosts(userID int64, page, size int64) ([]models.Post, int64, error) {
	zap.L().Info("开始查询用户点赞帖子",
		zap.Int64("user_id", userID),
		zap.Int64("page", page),
		zap.Int64("size", size))

	var posts []models.Post
	var total int64

	// 构建基础查询
	db := mysql.DB.Table("posts").
		Joins("JOIN likes ON posts.id = likes.target_id").
		Where("likes.user_id = ? AND likes.target_type = 1 AND posts.status != -1", userID)

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		zap.L().Error("获取点赞帖子总数失败",
			zap.Error(err),
			zap.Int64("user_id", userID))
		return nil, 0, err
	}
	zap.L().Info("获取到点赞帖子总数", zap.Int64("total", total))

	// 查询帖子数据
	err := db.Preload("Author").
		Preload("Tags").
		Preload("Images").
		Order("likes.created_at DESC").
		Offset(int((page - 1) * size)).
		Limit(int(size)).
		Find(&posts).Error

	if err != nil {
		zap.L().Error("查询点赞帖子详情失败",
			zap.Error(err),
			zap.Int64("user_id", userID),
			zap.Int64("page", page),
			zap.Int64("size", size))
		return nil, 0, err
	}

	zap.L().Info("成功查询点赞帖子详情",
		zap.Int("post_count", len(posts)),
		zap.Int64("user_id", userID))

	// 记录每个帖子的基本信息
	for i, post := range posts {
		zap.L().Debug("帖子详情",
			zap.Int("index", i),
			zap.Int64("post_id", post.ID),
			zap.String("title", post.Title),
			zap.Int64("author_id", *post.AuthorID),
			zap.Int("tag_count", len(post.Tags)),
			zap.Int("image_count", len(post.Images)))
	}

	return posts, total, nil
}

// GetUserFavoritePosts 获取用户收藏的帖子
func GetUserFavoritePosts(c *gin.Context) {
	userID, err := getCurrentUserIDInt64(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	zap.L().Info("开始获取用户收藏的帖子")

	// 获取分页参数
	page, size := getPageInfo(c)
	zap.L().Info("获取分页参数",
		zap.Int64("page", page),
		zap.Int64("size", size))

	// 获取帖子和总数
	posts, total, err := getUserFavoritePosts(userID, page, size)
	if err != nil {
		zap.L().Error("获取用户收藏帖子失败",
			zap.Error(err),
			zap.Int64("user_id", userID),
			zap.Int64("page", page),
			zap.Int64("size", size))
		ResponseError(c, CodeServerBusy)
		return
	}
	zap.L().Info("成功获取用户收藏帖子",
		zap.Int64("total", total),
		zap.Int("post_count", len(posts)))

	ResponseSuccess(c, gin.H{
		"posts":        posts,
		"total":        total,
		"current_page": page,
		"page_size":    size,
	})
}

func getUserFavoritePosts(userID int64, page, size int64) ([]models.Post, int64, error) {
	zap.L().Info("开始查询用户收藏帖子",
		zap.Int64("user_id", userID),
		zap.Int64("page", page),
		zap.Int64("size", size))

	var posts []models.Post
	var total int64

	// 构建基础查询
	db := mysql.DB.Table("posts").
		Joins("JOIN favorites ON posts.id = favorites.post_id").
		Where("favorites.user_id = ? AND posts.status != -1", userID)

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		zap.L().Error("获取收藏帖子总数失败",
			zap.Error(err),
			zap.Int64("user_id", userID))
		return nil, 0, err
	}
	zap.L().Info("获取到收藏帖子总数", zap.Int64("total", total))

	// 查询帖子数据
	err := db.Preload("Author").
		Preload("Tags").
		Preload("Images").
		Order("favorites.created_at DESC").
		Offset(int((page - 1) * size)).
		Limit(int(size)).
		Find(&posts).Error

	if err != nil {
		zap.L().Error("查询收藏帖子详情失败",
			zap.Error(err),
			zap.Int64("user_id", userID),
			zap.Int64("page", page),
			zap.Int64("size", size))
		return nil, 0, err
	}

	zap.L().Info("成功查询收藏帖子详情",
		zap.Int("post_count", len(posts)),
		zap.Int64("user_id", userID))

	// 记录每个帖子的基本信息
	for i, post := range posts {
		zap.L().Debug("帖子详情",
			zap.Int("index", i),
			zap.Int64("post_id", post.ID),
			zap.String("title", post.Title),
			zap.Int64("author_id", *post.AuthorID),
			zap.Int("tag_count", len(post.Tags)),
			zap.Int("image_count", len(post.Images)))
	}

	return posts, total, nil
}

// GetUserCommentedPosts 获取用户评论过的帖子
func GetUserCommentedPosts(c *gin.Context) {
	userID, err := getCurrentUserIDInt64(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	zap.L().Info("开始获取用户评论过的帖子")

	// 获取分页参数
	page, size := getPageInfo(c)
	zap.L().Info("获取分页参数",
		zap.Int64("page", page),
		zap.Int64("size", size))

	// 获取帖子和总数
	posts, total, err := getUserCommentedPosts(userID, page, size)
	if err != nil {
		zap.L().Error("获取用户评论过的帖子失败",
			zap.Error(err),
			zap.Int64("user_id", userID),
			zap.Int64("page", page),
			zap.Int64("size", size))
		ResponseError(c, CodeServerBusy)
		return
	}
	zap.L().Info("成功获取用户评论过的帖子",
		zap.Int64("total", total),
		zap.Int("post_count", len(posts)))

	ResponseSuccess(c, gin.H{
		"posts":        posts,
		"total":        total,
		"current_page": page,
		"page_size":    size,
	})
}

func getUserCommentedPosts(userID int64, page, size int64) ([]models.Post, int64, error) {
	zap.L().Info("开始查询用户评论过的帖子",
		zap.Int64("user_id", userID),
		zap.Int64("page", page),
		zap.Int64("size", size))

	var posts []models.Post
	var total int64

	// 构建基础查询
	db := mysql.DB.Table("posts").
		Joins("JOIN comments ON posts.id = comments.post_id").
		Where("comments.user_id = ? AND comments.status = 1 AND posts.status != -1", userID).
		Group("posts.id") // 使用 GROUP BY 去重，因为一个用户可能对同一个帖子评论多次

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		zap.L().Error("获取评论过的帖子总数失败",
			zap.Error(err),
			zap.Int64("user_id", userID))
		return nil, 0, err
	}
	zap.L().Info("获取到评论过的帖子总数", zap.Int64("total", total))

	// 查询帖子数据
	err := db.Preload("Author").
		Preload("Tags").
		Preload("Images").
		Order("MAX(comments.created_at) DESC"). // 按最新评论时间排序
		Offset(int((page - 1) * size)).
		Limit(int(size)).
		Find(&posts).Error

	if err != nil {
		zap.L().Error("查询评论过的帖子详情失败",
			zap.Error(err),
			zap.Int64("user_id", userID),
			zap.Int64("page", page),
			zap.Int64("size", size))
		return nil, 0, err
	}

	zap.L().Info("成功查询评论过的帖子详情",
		zap.Int("post_count", len(posts)),
		zap.Int64("user_id", userID))

	// 记录每个帖子的基本信息
	for i, post := range posts {
		zap.L().Debug("帖子详情",
			zap.Int("index", i),
			zap.Int64("post_id", post.ID),
			zap.String("title", post.Title),
			zap.Int64("author_id", *post.AuthorID),
			zap.Int("tag_count", len(post.Tags)),
			zap.Int("image_count", len(post.Images)))
	}

	return posts, total, nil
}
