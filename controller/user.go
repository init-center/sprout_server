package controller

import (
	"sprout_server/common/response"
	"sprout_server/common/response/code"
	"sprout_server/logic/user"
	"sprout_server/models"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (u *UserController) SignUp(c *gin.Context) {
	// 1. verify params
	var p models.ParamsSignUp
	if err := c.ShouldBindJSON(&p); err != nil {
		// params error
		// we use the shouldBindJSON and use the binding tag on model
		// the gin can help us to verify params
		response.Send(c, code.CodeInvalidParams)
		return
	}
	// 2. logic handle
	statusCode := user.Create(&p)
	// 3. response result
	response.Send(c, statusCode)

}
