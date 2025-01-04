package middlewares

import (
	"TalkSphere/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 对于API请求返回JSON错误
			if strings.HasPrefix(c.Request.URL.Path, "/api/") {
				c.JSON(200, gin.H{
					"code": 2003,
					"msg":  "需要登录",
				})
				c.Abort()
				return
			}
			// 对于页面请求重定向到登录页
			c.Redirect(302, "/")
			c.Abort()
			return
		}

		// 检查 Authorization 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			if strings.HasPrefix(c.Request.URL.Path, "/api/") {
				c.JSON(200, gin.H{
					"code": 2003,
					"msg":  "无效的认证格式",
				})
				c.Abort()
				return
			}
			c.Redirect(302, "/")
			c.Abort()
			return
		}

		// 解析 Token
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			if strings.HasPrefix(c.Request.URL.Path, "/api/") {
				c.JSON(200, gin.H{
					"code": 2003,
					"msg":  "无效的Token",
				})
				c.Abort()
				return
			}
			c.Redirect(302, "/")
			c.Abort()
			return
		}

		// 将当前请求的 userID 信息保存到请求的上下文c上
		c.Set("userID", mc.UserID)
		c.Set("username", mc.Username)
		c.Next()
	}
}
