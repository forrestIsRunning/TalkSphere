package router

import (
	"TalkSphere/controller"
	"github.com/gin-gonic/gin"
)

func RegisterBoardRoutes(r *gin.Engine) {
	r.GET("/boards", controller.GetAllBoards)
	r.POST("/boards", controller.CreateBoard)
	r.DELETE("/boards/:id", controller.DeleteBoard)
	r.PUT("/boards/:id", controller.UpdateBoard)
}
