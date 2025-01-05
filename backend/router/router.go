package router

import (
	"TalkSphere/pkg/logger"
	"TalkSphere/setting"
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	// 设置 gin 框架日志输出模式
	gin.SetMode(setting.Conf.GinConfig.Mode)

	r := gin.Default()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	RegisterUserRoutes(r)
	RegisterBoardRoutes(r)

	return r
}
