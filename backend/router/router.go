package router

import (
	"TalkSphere/middlewares"
	"TalkSphere/pkg/logger"
	"TalkSphere/setting"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	// 设置 gin 框架日志输出模式
	gin.SetMode(setting.Conf.GinConfig.Mode)

	r := gin.Default()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 静态文件服务
	r.Static("/static", "../frontend/static")

	// HTML模板
	r.LoadHTMLGlob("../frontend/templates/*")

	// 首页和登录页
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "auth.html", nil)
	})

	// 论坛主页（需要登录）
	r.GET("/forum", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "forum.html", nil)
	})

	// 个人资料页面（需要登录）
	r.GET("/profile", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "profile.html", nil)
	})

	UserInit(r)

	return r
}
