package router

import (
	"TalkSphere/controller"
	"TalkSphere/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterStatsRoutes(r *gin.RouterGroup) {
	adminGroup := r.Group("")
	adminGroup.Use(middleware.JWTAuthMiddleware(), middleware.AdminRequired())
	{
		adminGroup.GET("/admin/stats", controller.GetSystemStats)
	}
}
