package middleware

import (
	"TalkSphere/controller"
	"TalkSphere/pkg/jwt"
	"TalkSphere/pkg/rbac"
	"strconv"
	"strings"

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
		// parts[1] 是获取的 tokenString, 我们使用之前定义好的解析 JWT 的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		c.Set(controller.CtxtUserID, mc.UserID)
		c.Set(controller.CtxUserName, mc.Username)

		// 获取并设置用户角色
		role, err := rbac.GetUserRole(strconv.Itoa(int(mc.UserID)))
		if err != nil {
			controller.ResponseError(c, controller.CodeServerBusy)
			c.Abort()
			return
		}
		c.Set(ROLE, role)

		c.Next()
	}
}
