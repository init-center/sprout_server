package response

import (
	"sprout_server/common/response/code"

	"github.com/gin-gonic/gin"
)

type WithNoData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type WithData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Send(c *gin.Context, statusCode int) {
	c.JSON(code.HCode(statusCode), WithNoData{Code: statusCode, Message: code.Msg(statusCode)})
}

func SendWithData(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(code.HCode(statusCode), WithData{Code: statusCode, Message: code.Msg(statusCode), Data: data})
}
