package router

import (
	"TalkSphere/controller"
	"TalkSphere/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
	// 公开路由 - 不需要登录
	r.POST("/register", controller.RegisterHandler)
	r.POST("/login", controller.LoginHandler)

	// 需要登录的路由 - 使用路由组
	authGroup := r.Group("")
	authGroup.Use(middleware.JWTAuthMiddleware())
	{
		authGroup.GET("/profile", controller.GetUserProfile)
		authGroup.GET("/profile/:id", controller.GetUserProfile)
		authGroup.POST("/bio", controller.UpdateUserBio)
		authGroup.POST("/avatar", controller.UpdateUserAvatar)
	}
}
