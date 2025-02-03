package middleware

import (
	"TalkSphere/controller"
	"TalkSphere/pkg/rbac"
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// CheckAdmin 检查是否为管理员
func CheckAdmin(c *gin.Context) bool {
	userID, exists := c.Get(controller.CtxtUserID)
	if !exists {
		return false
	}

	// 从上下文获取 enforcer
	enforcer, exists := c.Get("enforcer")
	if !exists {
		return false
	}

	e, ok := enforcer.(*casbin.Enforcer)
	if !ok {
		return false
	}

	// 检查用户是否是管理员
	isAdmin, err := e.HasRoleForUser(fmt.Sprintf("%d", userID.(int64)), "admin")
	if err != nil {
		return false
	}

	return isAdmin
}

// AdminRequired 管理员权限要求中间件
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("enforcer", rbac.Enforcer)

		if !CheckAdmin(c) {
			c.AbortWithStatusJSON(403, gin.H{
				"code":    1005,
				"message": "需要管理员权限",
			})
			return
		}
		fmt.Println("--AdminRequired---")
		c.Next()
	}
}
