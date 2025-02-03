package router

import (
	"TalkSphere/controller"
	"TalkSphere/pkg/logger"
	"TalkSphere/setting"
	"time"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	// 设置 gin 框架日志输出模式
	gin.SetMode(setting.Conf.GinConfig.Mode)

	r := gin.Default()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	//r.POST("/auth/check", controller.CheckPermission)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 添加 CORS 中间件到所有路由
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/login", controller.LoginHandler)
	r.POST("/register", controller.RegisterHandler)
	r.GET("/user/check/:id", controller.CheckAdminPermission)

	apiGroup := r.Group("/api")
	RegisterUserRoutes(apiGroup)
	InitBoardRouter(apiGroup)
	InitPostRouter(apiGroup)
	InitInteractionRoutes(apiGroup)
	return r
}
