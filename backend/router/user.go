package router

import (
	"TalkSphere/controller"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	//用户登陆注册
	r.POST("/register", controller.RegisterHandler)
	r.POST("/login", controller.LoginHandler)

	//用户信息
	r.GET("/profile/:id", controller.GetUserProfile)
	r.POST("/bio", controller.UpdateUserBio)
	r.POST("/avatar", controller.UpdateUserAvatar)
}
