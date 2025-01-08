package router

import (
	"TalkSphere/controller"
	"TalkSphere/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	//用户登陆注册
	r.POST("/register", controller.RegisterHandler)
	r.POST("/login", controller.LoginHandler)

	//用户信息
	auth := r.Use(middleware.JWTAuthMiddleware())
	auth.GET("/profile", controller.GetUserProfile)
	auth.GET("/profile/:id", controller.GetUserProfile)
	auth.POST("/bio", controller.UpdateUserBio)
	auth.POST("/avatar", controller.UpdateUserAvatar)
}
