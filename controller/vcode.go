package controller

import (
	"sprout_server/common/response"
	"sprout_server/common/response/code"
	"sprout_server/logic/vcode"
	"sprout_server/models"

	"github.com/gin-gonic/gin"
)

type VCodeController struct{}

func (v *VCodeController) SendECode(c *gin.Context) {
	var p models.ParamsGetECode
	if err := c.ShouldBindJSON(&p); err != nil {
		// params error
		// we use the shouldBindJSON and use the binding tag on model
		// the gin can help us to verify params
		response.Send(c, code.CodeInvalidParams)
		return
	}

	statusCode := vcode.SendCodeToEmail(&p)
	response.Send(c, statusCode)

}
