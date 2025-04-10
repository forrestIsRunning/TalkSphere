package middleware

import (
	"TalkSphere/pkg/rbac"
	"net/http"

	"github.com/gin-gonic/gin"
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

		// 获取请求的路径和方法
		obj := c.Request.URL.Path
		act := c.Request.Method

		// 使用 Casbin 检查权限
		if rbac.CheckPermission(role, obj, act) {
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "Permission denied",
			})
			c.Abort()
		}
	}
}
