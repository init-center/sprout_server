package middlewares

import (
	"sprout_server/common/constant"
	"sprout_server/common/response"
	"sprout_server/common/response/code"
	"sprout_server/dao/mysql"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := c.Get(constant.CtxUidKey)
		if !ok {
			response.Send(c, code.CodeNeedLogin)
			c.Abort()
			return
		}

		group, err := mysql.GetUserGroup(uid.(string))
		if err != nil || group != constant.UserGroupAdmin {
			response.Send(c, code.CodePermissionDenied)
			c.Abort()
			return
		}

		// verify success
		c.Next()
	}
}
