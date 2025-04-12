package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/TalkSphere/backend/controller"
	"github.com/TalkSphere/backend/middleware"
	"github.com/TalkSphere/backend/pkg/logger"
	"github.com/TalkSphere/backend/setting"
)

func Setup() *gin.Engine {
	gin.SetMode(setting.Conf.GinConfig.Mode)

	r := gin.Default()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	//r.POST("/auth/check", controller.CheckPermission)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 添加 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 公开路由组 - 不需要认证
	publicGroup := r.Group("/api")
	{
		// 认证相关
		publicGroup.POST("/login", controller.LoginHandler)
		publicGroup.POST("/register", controller.RegisterHandler)

		// 注册各个模块的公开路由
		RegisterPublicRoutes(publicGroup)
	}

	// 需要认证和权限验证的路由组
	authGroup := r.Group("/api")
	authGroup.Use(middleware.JWTAuthMiddleware(), middleware.RBACMiddleware())
	{
		// 注册各个模块的需认证路由
		RegisterAuthRoutes(authGroup)
	}

	return r
}

// RegisterPublicRoutes 注册所有公开路由
func RegisterPublicRoutes(r *gin.RouterGroup) {
	// 板块相关
	r.GET("/boards", controller.GetAllBoards)

	// 帖子相关
	r.GET("/posts/:id", controller.GetPostDetail)
	r.GET("/posts/board/:board_id", controller.GetBoardPosts)

	// 评论相关
	r.GET("/comments/post/:post_id", controller.GetPostComments)
}

// RegisterAuthRoutes 注册所有需要认证的路由
func RegisterAuthRoutes(r *gin.RouterGroup) {
	// 用户相关
	r.GET("/profile", controller.GetUserProfile)
	r.POST("/bio", controller.UpdateUserBio)
	r.POST("/avatar", controller.UpdateUserAvatar)
	r.GET("/users", controller.GetUserLists)

	// 板块管理
	r.POST("/boards", controller.CreateBoard)
	r.PUT("/boards/:id", controller.UpdateBoard)
	r.DELETE("/boards/:id", controller.DeleteBoard)

	// 帖子相关
	r.POST("/posts", controller.CreatePost)
	r.PUT("/posts/:id", controller.UpdatePost)
	r.DELETE("/posts/:id", controller.DeletePost)
	r.GET("/posts/user", controller.GetUserPosts)
	r.GET("/posts/user/likes", controller.GetUserLikedPosts)
	r.GET("/posts/user/favorites", controller.GetUserFavoritePosts)
	r.POST("/posts/image", controller.UploadPostImage)

	// 互动相关
	r.POST("/comments", controller.CreateComment)
	r.DELETE("/comments/:id", controller.DeleteComment)
	r.POST("/likes", controller.CreateLike)
	r.GET("/likes/status", controller.GetLikeStatus)
	r.POST("/favorites/post/:post_id", controller.CreateFavorite)
	r.GET("/favorites", controller.GetUserFavorites)

	// 分析相关
	r.GET("/analysis/users/active", controller.GetActiveUsers)
	r.GET("/analysis/users/growth", controller.GetUsersGrowth)
	r.GET("/analysis/posts/active", controller.GetActivePosts)
	r.GET("/analysis/posts/growth", controller.GetPostsGrowth)
	r.GET("/analysis/posts/wordcloud", controller.GetPostsWordCloud)

	// 统计相关
	r.GET("/admin/stats", controller.GetSystemStats)

	// 权限管理相关
	//获取用户的所有权限
	r.GET("/permission/user/:user_id", controller.GetUserPermissions)
	//修改用户权限
	r.POST("/permission/user/:user_id", controller.UpdateUserPermissions)
	//校验用户权限
	r.GET("/permission/check/:user_id", controller.CheckPermission)
	//获取角色所拥有的权限
	r.GET("/permission/role", controller.GetRolePermissions)
	//获取用户的角色
	r.GET("/permission/user/role/:user_id", controller.GetUserRole)
}
