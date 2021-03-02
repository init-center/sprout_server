package middlewares

import (
	"sprout_server/common/constants"
	"sprout_server/common/jwt"
	"sprout_server/common/response"
	"sprout_server/common/response/code"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			response.Send(c, code.CodeNeedLogin)
			c.Abort()
			return
		}

		// split token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Send(c, code.CodeInvalidToken)
			c.Abort()
			return
		}
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			response.Send(c, code.CodeInvalidToken)
			c.Abort()
			return
		}

		// verify token success
		// set the uid to gin context
		c.Set(constants.CtxUidKey, mc.Uid)
		c.Next()
	}
}
