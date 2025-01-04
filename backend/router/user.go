package router

import (
	"TalkSphere/controller"
	"github.com/gin-gonic/gin"
)

func UserInit(r *gin.Engine) {
	auth := r.Group("/api")
	//用户登陆注册
	auth.POST("/register", controller.RegisterHandler)
	auth.POST("/login", controller.LoginHandler)

	//用户信息
	auth.GET("/profile", controller.GetUserProfile)
	auth.POST("/profile", controller.UpdateUserBio)
	auth.POST("/avatar", controller.UpdateUserAvatar)
}
