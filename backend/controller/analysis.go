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
	now := time.Now()

	// 准备时间范围
	sevenDaysAgo := now.AddDate(0, 0, -6) // 修改为-6，因为要包含今天
	sevenWeeksAgo := now.AddDate(0, 0, -7*7)
	sixMonthsAgo := now.AddDate(0, -6, 0)

	// 准备结果数组
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
			weekYear, weekNum := now.AddDate(0, 0, -i*7).ISOWeek()
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

	// 查询实际数据并更新
	var actualDaily, actualWeekly, actualMonthly []GrowthData

	// 查询每天的用户增长量
	if err := mysql.DB.Table("users").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("DATE(created_at) >= DATE(?)", sevenDaysAgo).
		Group("DATE(created_at)").
		Order("date DESC").
		Find(&actualDaily).Error; err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 查询每周的用户增长量
	if err := mysql.DB.Table("users").
		Select("YEARWEEK(created_at, 1) as date, COUNT(*) as count").
		Where("created_at >= ?", sevenWeeksAgo).
		Group("YEARWEEK(created_at, 1)").
		Order("date DESC").
		Find(&actualWeekly).Error; err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 将周数格式转换为更易读的格式
	for i := range actualWeekly {
		yearWeek := actualWeekly[i].Date
		year := yearWeek[:4]
		week := yearWeek[4:]
		actualWeekly[i].Date = fmt.Sprintf("%s-W%s", year, week)
	}

	// 查询每月的用户增长量
	if err := mysql.DB.Table("users").
		Select("DATE_FORMAT(created_at, '%Y-%m') as date, COUNT(*) as count").
		Where("created_at >= ?", sixMonthsAgo).
		Group("DATE_FORMAT(created_at, '%Y-%m')").
		Order("date DESC").
		Find(&actualMonthly).Error; err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 更新实际数据
	for _, actual := range actualDaily {
		for i, growth := range dailyGrowth {
			if growth.Date == actual.Date {
				dailyGrowth[i].Count = actual.Count
			}
		}
	}

	for _, actual := range actualWeekly {
		for i, growth := range weeklyGrowth {
			if growth.Date == actual.Date {
				weeklyGrowth[i].Count = actual.Count
			}
		}
	}

	for _, actual := range actualMonthly {
		for i, growth := range monthlyGrowth {
			if growth.Date == actual.Date {
				monthlyGrowth[i].Count = actual.Count
			}
		}
	}

	// 返回结果
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
	now := time.Now()

	// 准备时间范围
	sevenDaysAgo := now.AddDate(0, 0, -6) // 包含今天
	sevenWeeksAgo := now.AddDate(0, 0, -7*7)
	sixMonthsAgo := now.AddDate(0, -6, 0)

	// 准备结果数组
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
			weekYear, weekNum := now.AddDate(0, 0, -i*7).ISOWeek()
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

	// 查询实际数据并更新
	var actualDaily, actualWeekly, actualMonthly []GrowthData

	// 查询每天的帖子增长量
	if err := mysql.DB.Table("posts").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("DATE(created_at) >= DATE(?) AND status = 1", sevenDaysAgo).
		Group("DATE(created_at)").
		Order("date DESC").
		Find(&actualDaily).Error; err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 查询每周的帖子增长量
	if err := mysql.DB.Table("posts").
		Select("YEARWEEK(created_at, 1) as date, COUNT(*) as count").
		Where("created_at >= ? AND status = 1", sevenWeeksAgo).
		Group("YEARWEEK(created_at, 1)").
		Order("date DESC").
		Find(&actualWeekly).Error; err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 将周数格式转换为更易读的格式
	for i := range actualWeekly {
		yearWeek := actualWeekly[i].Date
		year := yearWeek[:4]
		week := yearWeek[4:]
		actualWeekly[i].Date = fmt.Sprintf("%s-W%s", year, week)
	}

	// 查询每月的帖子增长量
	if err := mysql.DB.Table("posts").
		Select("DATE_FORMAT(created_at, '%Y-%m') as date, COUNT(*) as count").
		Where("created_at >= ? AND status = 1", sixMonthsAgo).
		Group("DATE_FORMAT(created_at, '%Y-%m')").
		Order("date DESC").
		Find(&actualMonthly).Error; err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	// 更新实际数据
	for _, actual := range actualDaily {
		for i, growth := range dailyGrowth {
			if growth.Date == actual.Date {
				dailyGrowth[i].Count = actual.Count
			}
		}
	}

	for _, actual := range actualWeekly {
		for i, growth := range weeklyGrowth {
			if growth.Date == actual.Date {
				weeklyGrowth[i].Count = actual.Count
			}
		}
	}

	for _, actual := range actualMonthly {
		for i, growth := range monthlyGrowth {
			if growth.Date == actual.Date {
				monthlyGrowth[i].Count = actual.Count
			}
		}
	}

	// 返回结果
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
