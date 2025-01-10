package router

import (
	"TalkSphere/controller"
	"TalkSphere/middleware"

	"github.com/gin-gonic/gin"
)

func InitInteractionRoutes(r *gin.Engine) {
	// 评论相关路由
	commentGroup := r.Group("/comments")
	{
		// 创建评论
		commentGroup.POST("", middleware.JWTAuthMiddleware(), controller.CreateComment)
		// 获取帖子的评论列表
		commentGroup.GET("/post/:post_id", controller.GetPostComments)
		// 删除评论
		commentGroup.DELETE("/:id", middleware.JWTAuthMiddleware(), controller.DeleteComment)
	}

	// 点赞相关路由
	likeGroup := r.Group("/likes")
	{
		likeGroup.POST("", middleware.JWTAuthMiddleware(), controller.CreateLike)
		//GetLikeStatus
		likeGroup.GET("/status", middleware.JWTAuthMiddleware(), controller.GetLikeStatus)
	}

	// 收藏相关路由
	favoriteGroup := r.Group("/favorites")
	{
		// 收藏/取消收藏帖子
		favoriteGroup.POST("/post/:post_id", middleware.JWTAuthMiddleware(), controller.CreateFavorite)
		// 获取用户收藏列表
		favoriteGroup.GET("", middleware.JWTAuthMiddleware(), controller.GetUserFavorites)
	}
}
