package controller

import (
	"TalkSphere/dao/mysql"
	"TalkSphere/pkg/oss"
	"TalkSphere/setting"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

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

// GetActiveUsers 活跃用户 top10
func GetActiveUsers(context *gin.Context) {

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

// GetActivePosts 帖子 top10
func GetActivePosts(context *gin.Context) {

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
		"http://127.0.0.1:5000/generate_wordcloud",
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
