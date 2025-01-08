package router

import (
	"TalkSphere/controller"
	"TalkSphere/middleware"

	"github.com/gin-gonic/gin"
)

func InitPostRouter(r *gin.Engine) {
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
	// 获取用户的帖子列表
	postGroup.GET("/user/:user_id", controller.GetUserPosts)
	// 获取板块下的帖子列表
	postGroup.GET("/board/:board_id", controller.GetBoardPosts)

	//照片上传
	postGroup.POST("/image", controller.UploadPostImage)
}
