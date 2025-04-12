package middleware

import (
	"net/http"

	"github.com/TalkSphere/backend/controller"
	"github.com/TalkSphere/backend/pkg/rbac"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const ROLE = "role"

// RBACMiddleware 统一的权限检查中间件
func RBACMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户角色
		role := c.GetString(ROLE)
		if role == "" {
			role = "guest"
		}

		// 超级管理员直接放行
		if role == "super_admin" {
			c.Next()
			return
		}

		// 获取请求的路径和方法
		obj := c.Request.URL.Path
		if obj == "" {
			obj = "/"
		}
		act := c.Request.Method

		// 获取用户ID
		userID, exists := c.Get(controller.CtxtUserID)
		if !exists {
			zap.L().Error("user ID not found in context",
				zap.String("path", obj),
				zap.String("method", act))
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未授权",
			})
			c.Abort()
			return
		}

		// 确保 userID 是字符串类型
		userIDStr, ok := userID.(string)
		if !ok {
			zap.L().Error("user ID is not string type",
				zap.String("path", obj),
				zap.String("method", act),
				zap.Any("userID", userID))
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
			})
			c.Abort()
			return
		}

		// 使用 CheckUserPermission 检查权限（考虑角色继承）
		if rbac.CheckUserPermission(userIDStr, obj, act) {
			zap.L().Info("permission granted",
				zap.String("userID", userIDStr),
				zap.String("role", role),
				zap.String("path", obj),
				zap.String("method", act))
			c.Next()
		} else {
			zap.L().Warn("permission denied",
				zap.String("userID", userIDStr),
				zap.String("role", role),
				zap.String("path", obj),
				zap.String("method", act))
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "没有访问权限",
				"data": gin.H{
					"path":   obj,
					"method": act,
					"role":   role,
				},
			})
			c.Abort()
		}
	}
}
