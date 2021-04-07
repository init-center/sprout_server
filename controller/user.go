package controller

import (
	"sprout_server/common/constants"
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

func (u *UserController) GetPublicUserInfo(c *gin.Context) {
	var p models.UriGetUserInfo
	if err := c.ShouldBindUri(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	userInfo, statusCode := user.GetPublicUserInfo(p)

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}
	response.SendWithData(c, statusCode, userInfo)

}

func (u *UserController) GetPrivateUserInfo(c *gin.Context) {
	var p models.UriGetUserInfo
	uid, exists := c.Get(constants.CtxUidKey)
	if !exists {
		response.Send(c, code.CodeNeedLogin)
		return
	}
	p.Uid = uid.(string)
	userInfo, statusCode := user.GetPrivateUserInfo(p)

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}
	response.SendWithData(c, statusCode, userInfo)

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

func (u *UserController) GetBanTime(c *gin.Context) {
	var p models.UriGetBanTime
	if err := c.ShouldBindQuery(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	banTime, statusCode := user.GetBanTime(p.Uid)

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}
	response.SendWithData(c, statusCode, banTime)

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

func (u *UserController) UpdateUser(c *gin.Context) {
	var p models.ParamsUpdateUser
	if err := c.ShouldBindJSON(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	var uid, exists = c.Get(constants.CtxUidKey)
	if !exists {
		response.Send(c, code.CodeNeedLogin)
		return
	}

	statusCode := user.UpdateUser(&p, uid.(string))
	response.Send(c, statusCode)

}

func (u *UserController) UpdatePassword(c *gin.Context) {
	var p models.ParamsUpdatePassword
	if err := c.ShouldBindJSON(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}

	statusCode := user.UpdatePassword(&p)
	response.Send(c, statusCode)

}

func (u *UserController) DeleteUser(c *gin.Context) {
	var uid, exists = c.Get(constants.CtxUidKey)
	if !exists {
		response.Send(c, code.CodeNeedLogin)
		return
	}

	statusCode := user.DeleteUser(uid.(string))
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
