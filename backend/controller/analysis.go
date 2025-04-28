package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/TalkSphere/backend/pkg/mysql"
	"github.com/TalkSphere/backend/pkg/oss"
	"github.com/TalkSphere/backend/setting"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// 定义一个通用的结构体类型
type GrowthData struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// WordCloudItem 词云数据结构
type WordCloudItem struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

// GetActiveUsers 活跃用户
func GetActiveUsers(c *gin.Context) {
	// 定义返回的数据结构
	type ActiveUser struct {
		UserID        int64   `json:"user_id"`
		Username      string  `json:"username"`
		AvatarURL     string  `json:"avatar_url"`
		PostCount     int     `json:"post_count"`
		LikeCount     int     `json:"like_received_count"`
		FavoriteCount int     `json:"favorite_received_count"`
		LastLoginAt   string  `json:"last_login_at"`
		ActivityScore float64 `json:"activity_score"`
	}

	// 获取时间范围参数
	timeRange := c.Query("time_range")
	if timeRange == "" {
		timeRange = "daily"
	}

	// 构建查询SQL
	query := `
        WITH user_stats AS (
            SELECT 
                u.id,
                u.username,
                u.avatar_url,
                u.last_login_at,
                COUNT(DISTINCT p.id) as post_count,
                COALESCE(SUM(p.like_count), 0) as total_likes,
                COALESCE(SUM(p.favorite_count), 0) as total_favorites,
                (
                    0.3 * (CASE 
                        WHEN u.last_login_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR) THEN 100
                        WHEN u.last_login_at >= DATE_SUB(NOW(), INTERVAL 72 HOUR) THEN 70
                        WHEN u.last_login_at >= DATE_SUB(NOW(), INTERVAL 7 DAY) THEN 40
                        ELSE 10
                    END) +
                    0.3 * COUNT(DISTINCT p.id) * 10 +
                    0.2 * COALESCE(SUM(p.like_count), 0) +
                    0.2 * COALESCE(SUM(p.favorite_count), 0)
                ) as activity_score
            FROM users u
            LEFT JOIN posts p ON u.id = p.author_id AND p.status = 1
            WHERE u.status = 1
    `

	// 根据时间范围添加条件
	switch timeRange {
	case "daily":
		query += " AND (p.created_at IS NULL OR DATE(p.created_at) = CURDATE())"
	case "weekly":
		query += " AND (p.created_at IS NULL OR p.created_at >= DATE_SUB(CURDATE(), INTERVAL 7 DAY))"
	case "monthly":
		query += " AND (p.created_at IS NULL OR p.created_at >= DATE_SUB(CURDATE(), INTERVAL 30 DAY))"
	default:
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 完成查询SQL
	query += `
            GROUP BY u.id, u.username, u.avatar_url, u.last_login_at
        )
        SELECT *
        FROM user_stats
        WHERE activity_score > 0
        ORDER BY activity_score DESC
        LIMIT 10
    `

	// 执行查询
	var activeUsers []ActiveUser
	err := mysql.DB.Raw(query).Scan(&activeUsers).Error

	if err != nil {
		zap.L().Error("get active users failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 添加调试日志
	zap.L().Debug("active users query result",
		zap.String("time_range", timeRange),
		zap.Any("users", activeUsers),
		zap.String("query", query))

	// 格式化时间
	for i := range activeUsers {
		if activeUsers[i].LastLoginAt != "" {
			t, err := time.Parse(time.RFC3339, activeUsers[i].LastLoginAt)
			if err == nil {
				activeUsers[i].LastLoginAt = t.Format("2006-01-02 15:04:05")
			}
		}
	}

	ResponseSuccess(c, gin.H{
		"active_users": activeUsers,
		"time_range":   timeRange,
	})
}

// GetUsersGrowth 用户增长量
func GetUsersGrowth(c *gin.Context) {
	// 获取最近7天的完整日期列表
	now := time.Now()
	dailyGrowth := make([]GrowthData, 7)
	weeklyGrowth := make([]GrowthData, 7)
	monthlyGrowth := make([]GrowthData, 6)

	// 初始化日期和默认值
	for i := 0; i < 7; i++ {
		// 填充每日数据
		date := now.AddDate(0, 0, -i)
		dailyGrowth[i] = GrowthData{
			Date:  date.Format("2006-01-02"),
			Count: 0,
		}

		// 填充每周数据
		if i < 7 {
			weekDate := now.AddDate(0, 0, -i*7)
			weekYear, weekNum := weekDate.ISOWeek()
			weeklyGrowth[i] = GrowthData{
				Date:  fmt.Sprintf("%d-W%02d", weekYear, weekNum),
				Count: 0,
			}
		}

		// 填充每月数据
		if i < 6 {
			monthDate := now.AddDate(0, -i, 0)
			monthlyGrowth[i] = GrowthData{
				Date:  monthDate.Format("2006-01"),
				Count: 0,
			}
		}
	}

	// 查询每日数据
	var dailyResults []GrowthData
	err := mysql.DB.Raw(`
        SELECT 
            DATE(created_at) as date,
            COUNT(*) as count
        FROM users
        WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)
        GROUP BY DATE(created_at)
        ORDER BY date DESC
    `).Find(&dailyResults).Error

	if err != nil {
		zap.L().Error("failed to query daily growth", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 修改周数据查询
	var weeklyResults []GrowthData
	err = mysql.DB.Raw(`
        SELECT 
            CONCAT(
                YEAR(MIN(created_at)),
                '-W',
                LPAD(WEEK(MIN(created_at), 1), 2, '0')
            ) as date,
            COUNT(*) as count
        FROM users
        WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 7 WEEK)
        GROUP BY YEAR(created_at), WEEK(created_at, 1)
        ORDER BY YEAR(created_at) DESC, WEEK(created_at, 1) DESC
    `).Find(&weeklyResults).Error

	if err != nil {
		zap.L().Error("failed to query weekly growth", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 修改月数据查询
	var monthlyResults []GrowthData
	err = mysql.DB.Raw(`
        SELECT 
            DATE_FORMAT(MIN(created_at), '%Y-%m') as date,
            COUNT(*) as count
        FROM users
        WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 MONTH)
        GROUP BY YEAR(created_at), MONTH(created_at)
        ORDER BY YEAR(created_at) DESC, MONTH(created_at) DESC
    `).Find(&monthlyResults).Error

	if err != nil {
		zap.L().Error("failed to query monthly growth", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 更新日增长数据
	for _, result := range dailyResults {
		resultDate := result.Date[:10]
		for i := range dailyGrowth {
			if dailyGrowth[i].Date == resultDate {
				dailyGrowth[i].Count = result.Count
			}
		}
	}

	// 更新周增长数据
	for _, result := range weeklyResults {
		for i := range weeklyGrowth {
			if weeklyGrowth[i].Date == result.Date {
				weeklyGrowth[i].Count = result.Count
			}
		}
	}

	// 更新月增长数据
	for _, result := range monthlyResults {
		for i := range monthlyGrowth {
			if monthlyGrowth[i].Date == result.Date {
				monthlyGrowth[i].Count = result.Count
			}
		}
	}

	ResponseSuccess(c, gin.H{
		"daily_growth":   dailyGrowth,
		"weekly_growth":  weeklyGrowth,
		"monthly_growth": monthlyGrowth,
	})
}

// GetActivePosts 活跃帖子
func GetActivePosts(c *gin.Context) {
	// 定义返回的数据结构
	type ActivePost struct {
		PostID        int64   `json:"post_id"`
		Title         string  `json:"title"`
		Content       string  `json:"content"`
		AuthorID      int64   `json:"author_id"`
		AuthorName    string  `json:"author_name"`
		AuthorAvatar  string  `json:"author_avatar"`
		CreatedAt     string  `json:"created_at"`
		LikeCount     int     `json:"like_count"`
		FavoriteCount int     `json:"favorite_count"`
		CommentCount  int     `json:"comment_count"`
		ActivityScore float64 `json:"activity_score"`
	}

	// 获取时间范围参数
	timeRange := c.Query("time_range")
	if timeRange == "" {
		timeRange = "daily"
	}

	// 获取最新帖子的时间作为参考点
	var latestPostTime time.Time
	err := mysql.DB.Raw("SELECT MAX(created_at) FROM posts WHERE status = 1").Scan(&latestPostTime).Error
	if err != nil {
		zap.L().Error("get latest post time failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 构建查询SQL
	query := `
        WITH post_stats AS (
            SELECT 
                p.id as post_id,
                p.title,
                LEFT(p.content, 200) as content,
                p.author_id,
                u.username as author_name,
                u.avatar_url as author_avatar,
                p.created_at,
                p.like_count,
                p.favorite_count,
                p.comment_count,
                (
                    10 + /* 基础分 */
                    0.4 * (CASE 
                        WHEN p.created_at >= DATE_SUB(?, INTERVAL 24 HOUR) THEN 100
                        WHEN p.created_at >= DATE_SUB(?, INTERVAL 72 HOUR) THEN 70
                        WHEN p.created_at >= DATE_SUB(?, INTERVAL 7 DAY) THEN 40
                        ELSE 10
                    END) +
                    0.3 * COALESCE(p.like_count, 0) +
                    0.2 * COALESCE(p.favorite_count, 0) +
                    0.1 * COALESCE(p.comment_count, 0)
                ) as activity_score
            FROM posts p
            LEFT JOIN users u ON p.author_id = u.id
            WHERE p.status = 1
    `

	// 根据时间范围添加条件
	var timeCondition string
	switch timeRange {
	case "daily":
		timeCondition = "DATE(p.created_at) = DATE(?)"
	case "weekly":
		timeCondition = "p.created_at >= DATE_SUB(DATE(?), INTERVAL 7 DAY)"
	case "monthly":
		timeCondition = "p.created_at >= DATE_SUB(DATE(?), INTERVAL 30 DAY)"
	default:
		ResponseError(c, CodeInvalidParam)
		return
	}

	query += " AND " + timeCondition

	// 完成查询SQL
	query += `
        )
        SELECT *
        FROM post_stats
        ORDER BY activity_score DESC
        LIMIT 10
    `

	// 执行查询，使用最新帖子时间作为参考点
	var activePosts []ActivePost
	err = mysql.DB.Raw(query,
		latestPostTime, // 用于24小时判断
		latestPostTime, // 用于72小时判断
		latestPostTime, // 用于7天判断
		latestPostTime, // 用于时间范围条件
	).Scan(&activePosts).Error

	if err != nil {
		zap.L().Error("get active posts failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 记录查询结果
	zap.L().Info("query result",
		zap.Int("found_posts", len(activePosts)),
		zap.String("time_range", timeRange),
		zap.Time("reference_time", latestPostTime))

	// 格式化时间
	for i := range activePosts {
		if activePosts[i].CreatedAt != "" {
			t, err := time.Parse(time.RFC3339, activePosts[i].CreatedAt)
			if err == nil {
				activePosts[i].CreatedAt = t.Format("2006-01-02 15:04:05")
			}
		}
		// 截断内容
		if len(activePosts[i].Content) > 200 {
			activePosts[i].Content = activePosts[i].Content[:200] + "..."
		}
	}

	ResponseSuccess(c, gin.H{
		"active_posts": activePosts,
		"time_range":   timeRange,
	})
}

// GetPostsGrowth 帖子增长量
func GetPostsGrowth(c *gin.Context) {
	// 获取最近7天的完整日期列表
	now := time.Now()
	dailyGrowth := make([]GrowthData, 7)
	weeklyGrowth := make([]GrowthData, 7)
	monthlyGrowth := make([]GrowthData, 6)

	// 初始化日期和默认值
	for i := 0; i < 7; i++ {
		// 填充每日数据
		date := now.AddDate(0, 0, -i)
		dailyGrowth[i] = GrowthData{
			Date:  date.Format("2006-01-02"),
			Count: 0,
		}

		// 填充每周数据
		if i < 7 {
			weekDate := now.AddDate(0, 0, -i*7)
			weekYear, weekNum := weekDate.ISOWeek()
			weeklyGrowth[i] = GrowthData{
				Date:  fmt.Sprintf("%d-W%02d", weekYear, weekNum),
				Count: 0,
			}
		}

		// 填充每月数据
		if i < 6 {
			monthDate := now.AddDate(0, -i, 0)
			monthlyGrowth[i] = GrowthData{
				Date:  monthDate.Format("2006-01"),
				Count: 0,
			}
		}
	}

	// 查询每日数据
	var dailyResults []GrowthData
	err := mysql.DB.Raw(`
        SELECT 
            DATE(created_at) as date,
            COUNT(*) as count
        FROM posts
        WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)
        AND status = 1
        GROUP BY DATE(created_at)
        ORDER BY date DESC
    `).Find(&dailyResults).Error

	if err != nil {
		zap.L().Error("failed to query daily growth", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 查询每周数据
	var weeklyResults []GrowthData
	err = mysql.DB.Raw(`
        SELECT 
            CONCAT(
                YEAR(MIN(created_at)),
                '-W',
                LPAD(WEEK(MIN(created_at), 1), 2, '0')
            ) as date,
            COUNT(*) as count
        FROM posts
        WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 7 WEEK)
        AND status = 1
        GROUP BY YEAR(created_at), WEEK(created_at, 1)
        ORDER BY YEAR(created_at) DESC, WEEK(created_at, 1) DESC
    `).Find(&weeklyResults).Error

	if err != nil {
		zap.L().Error("failed to query weekly growth", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 查询每月数据
	var monthlyResults []GrowthData
	err = mysql.DB.Raw(`
        SELECT 
            DATE_FORMAT(MIN(created_at), '%Y-%m') as date,
            COUNT(*) as count
        FROM posts
        WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 MONTH)
        AND status = 1
        GROUP BY YEAR(created_at), MONTH(created_at)
        ORDER BY YEAR(created_at) DESC, MONTH(created_at) DESC
    `).Find(&monthlyResults).Error

	if err != nil {
		zap.L().Error("failed to query monthly growth", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 更新日增长数据
	for _, result := range dailyResults {
		resultDate := result.Date[:10]
		for i := range dailyGrowth {
			if dailyGrowth[i].Date == resultDate {
				dailyGrowth[i].Count = result.Count
			}
		}
	}

	// 更新周增长数据
	for _, result := range weeklyResults {
		for i := range weeklyGrowth {
			if weeklyGrowth[i].Date == result.Date {
				weeklyGrowth[i].Count = result.Count
			}
		}
	}

	// 更新月增长数据
	for _, result := range monthlyResults {
		for i := range monthlyGrowth {
			if monthlyGrowth[i].Date == result.Date {
				monthlyGrowth[i].Count = result.Count
			}
		}
	}

	ResponseSuccess(c, gin.H{
		"daily_growth":   dailyGrowth,
		"weekly_growth":  weeklyGrowth,
		"monthly_growth": monthlyGrowth,
	})
}

// GetPostsWordCloud 词云图
func GetPostsWordCloud(c *gin.Context) {
	// 查询所有有效帖子的内容
	var posts []struct {
		Content string
	}

	if err := mysql.DB.Table("posts").
		Select("content").
		Where("status = 1").
		Find(&posts).Error; err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 合并所有内容
	var allContent string
	for _, post := range posts {
		allContent += " " + post.Content
	}

	// 准备请求数据
	requestData := map[string]interface{}{
		"text": allContent,
	}

	// 将请求数据转换为JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		zap.L().Error("JSON编码失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 发送请求到Python服务
	resp, err := http.Post(
		"http://127.0.0.1:5008/generate_wordcloud",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		zap.L().Error("请求Python服务失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	defer resp.Body.Close()

	// 解析响应
	var result struct {
		Success   bool   `json:"success"`
		ImagePath string `json:"image_path"`
		FileSize  int64  `json:"file_size"`
		Error     string `json:"error,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		zap.L().Error("解析Python服务响应失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	if !result.Success {
		zap.L().Error("Python服务生成词云图失败", zap.String("error", result.Error))
		ResponseError(c, CodeServerBusy)
		return
	}

	_, err = os.ReadFile(result.ImagePath)
	if err != nil {
		zap.L().Error("读取图片文件失败",
			zap.Error(err),
			zap.String("path", result.ImagePath))
		ResponseError(c, CodeServerBusy)
		return
	}

	fd, err := os.Open(result.ImagePath)
	if err != nil {
		zap.L().Error("读取图片文件失败",
			zap.Error(err),
			zap.String("path", result.ImagePath))
		ResponseError(c, CodeServerBusy)
	}
	defer fd.Close()
	if _, err = oss.Client.Object.Put(c, result.ImagePath, fd, nil); err != nil {
		zap.L().Error("上传词云图失败",
			zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 生成访问URL
	url := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s",
		setting.Conf.OSSConfig.BucketName,
		setting.Conf.OSSConfig.Region,
		result.ImagePath)
	defer os.Remove(result.ImagePath)
	ResponseSuccess(c, gin.H{
		"word_cloud_url": url,
	})
}
