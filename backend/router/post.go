package router

import (
	"TalkSphere/controller"
	"TalkSphere/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterPostRouter(r *gin.RouterGroup) {
	postGroup := r.Group("/posts")
	postGroup.Use(middleware.JWTAuthMiddleware()) // 需要登录才能访问

	// 创建帖子
	postGroup.POST("", controller.CreatePost)
	// 获取帖子详情
	postGroup.GET("/:id", controller.GetPostDetail)
	// 删除帖子
	postGroup.DELETE("/:id", controller.DeletePost)
	// 更新帖子
	postGroup.PUT("/:id", controller.UpdatePost)

	// 获取板块下的帖子列表
	postGroup.GET("/board/:board_id", controller.GetBoardPosts)

	// 获取用户发表的帖子列表
	postGroup.GET("/user", controller.GetUserPosts)
	//获取用户点赞帖子
	postGroup.GET("/user/likes", controller.GetUserLikedPosts)
	//获取用户收藏帖子
	postGroup.GET("/user/favorites", controller.GetUserFavoritePosts)

	//照片上传
	postGroup.POST("/image", controller.UploadPostImage)
}
