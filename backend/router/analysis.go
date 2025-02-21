package router

import (
	"TalkSphere/controller"
	"TalkSphere/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAnalysisRouter(r *gin.RouterGroup) {
	authGroup := r.Group("/analysis")
	authGroup.Use(middleware.JWTAuthMiddleware(), middleware.AdminRequired())
	{
		authGroup.GET("/users/active", controller.GetActiveUsers)
		authGroup.GET("/users/growth", controller.GetUsersGrowth)

		authGroup.GET("/posts/active", controller.GetActivePosts)
		authGroup.GET("/posts/growth", controller.GetPostsGrowth)
		authGroup.GET("posts/wordcloud", controller.GetPostsWordCloud)
	}
}
