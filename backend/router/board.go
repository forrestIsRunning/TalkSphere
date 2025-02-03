package router

import (
	"TalkSphere/controller"
	"TalkSphere/middleware"

	"github.com/gin-gonic/gin"
)

func InitBoardRouter(r *gin.RouterGroup) {
	// 公开接口 - 不需要登录
	publicGroup := r.Group("/boards")
	{
		publicGroup.GET("", controller.GetAllBoards) // 获取板块列表不需要登录
	}

	// 需要登录的接口
	authGroup := r.Group("/boards")
	authGroup.Use(middleware.JWTAuthMiddleware(), middleware.AdminRequired())
	{
		authGroup.POST("", controller.CreateBoard)       // 创建板块需要登录
		authGroup.PUT("/:id", controller.UpdateBoard)    // 更新板块需要登录
		authGroup.DELETE("/:id", controller.DeleteBoard) // 删除板块需要登录
	}
}
