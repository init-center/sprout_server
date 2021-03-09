package controller

import (
	"sprout_server/common/response"
	"sprout_server/common/response/code"
	"sprout_server/logic/user"
	"sprout_server/models"
	"sprout_server/models/queryfields"

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

func (u *UserController) AdminGetUsers(c *gin.Context) {
	var p queryfields.UserQueryFields
	if err := c.ShouldBindQuery(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	users, statusCode := user.AdminGetUsers(&p)

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}
	response.SendWithData(c, statusCode, users)

}

func (u *UserController) AdminUpdateUser(c *gin.Context) {
	var p models.ParamsAdminUpdateUser
	if err := c.ShouldBindJSON(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	var uri models.UriUpdateUser
	if err := c.ShouldBindUri(&uri); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}

	statusCode := user.AdminUpdateUser(&p, &uri)
	response.Send(c, statusCode)

}

func (u *UserController) BanUser(c *gin.Context) {
	var p models.ParamsBanUser
	if err := c.ShouldBindJSON(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	var uri models.UriUpdateUser
	if err := c.ShouldBindUri(&uri); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}

	statusCode := user.BanUser(&p, &uri)
	response.Send(c, statusCode)
}

func (u *UserController) UnblockUser(c *gin.Context) {
	var uri models.UriUpdateUser
	if err := c.ShouldBindUri(&uri); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}

	statusCode := user.UnblockUser(&uri)
	response.Send(c, statusCode)
}
