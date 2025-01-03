package router

import (
	"TalkSphere/controller"
	"TalkSphere/logger"
	"TalkSphere/setting"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	// 设置 gin 框架日志输出模式
	gin.SetMode(setting.Conf.GinConfig.Mode)

	r := gin.Default()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	// 添加静态文件服务
	r.Static("/static", "../frontend/static")

	// 添加HTML模板
	r.LoadHTMLGlob("../frontend/templates/*")

	// 添加首页路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "auth.html", nil)
	})

	r.POST("/api/register", controller.RegisterHandler)
	r.POST("/api/login", controller.LoginHandler)

	return r
}
