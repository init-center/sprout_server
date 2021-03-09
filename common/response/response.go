package response

import (
	"sprout_server/common/response/code"
	"time"

	"github.com/gin-gonic/gin"
)

type WithNoData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Time    int64  `json:"time"`
}

type WithData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Time    int64       `json:"time"`
}

func Send(c *gin.Context, statusCode int) {
	c.JSON(code.HCode(statusCode), WithNoData{Code: statusCode, Message: code.Msg(statusCode), Time: time.Now().Unix()})
}

func SendWithData(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(code.HCode(statusCode), WithData{Code: statusCode, Message: code.Msg(statusCode), Data: data, Time: time.Now().Unix()})
}
