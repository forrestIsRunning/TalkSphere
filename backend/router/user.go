package router

import (
	"TalkSphere/controller"
	"TalkSphere/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
	authGroup := r.Group("")
	authGroup.Use(middleware.JWTAuthMiddleware())
	{
		authGroup.GET("/profile", controller.GetUserProfile)
		authGroup.POST("/bio", controller.UpdateUserBio)
		authGroup.POST("/avatar", controller.UpdateUserAvatar)
		authGroup.GET("/users", middleware.AdminRequired(), controller.GetUserLists)
	}
}
