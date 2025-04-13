package middleware

import (
	"strings"

	"github.com/TalkSphere/backend/controller"
	"github.com/TalkSphere/backend/pkg/jwt"
	"github.com/TalkSphere/backend/pkg/rbac"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带 Token 的三种方式：1. 放在请求头 2. 放在请求体 3. 放在 URI
		// 这里假设 Token 放在 Header 的 Authorization 中，并使用 Bearer 开头
		// Authorization: Bearer xxx.xxx.xxx
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		// 按空格切割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		// parts[1] 是获取的 tokenString
		token := parts[1]
		// 如果是 guest_token_anonymous，设置默认的 guest 用户信息
		if token == "guest_token_anonymous" {
			c.Set(controller.CtxtUserID, "0")
			c.Set(controller.CtxUserName, "guest")
			c.Set(ROLE, "guest")
			c.Next()
			return
		}
		// 否则正常解析 JWT token
		mc, err := jwt.ParseToken(token)
		if err != nil {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		c.Set(controller.CtxtUserID, mc.UserID)
		c.Set(controller.CtxUserName, mc.Username)

		// 获取并设置用户角色
		role, err := rbac.GetUserRole(mc.UserID)
		if err != nil {
			controller.ResponseError(c, controller.CodeServerBusy)
			c.Abort()
			return
		}
		c.Set(ROLE, role)

		c.Next()
	}
}
