package controller

import (
	"sprout_server/common/constants"
	"sprout_server/common/response"
	"sprout_server/common/response/code"
	"sprout_server/logic/session"
	"sprout_server/models"

	"github.com/gin-gonic/gin"
)

type SessionController struct{}

func (s *SessionController) SignIn(c *gin.Context) {
	// 1. verify params
	var p models.ParamsSignIn
	if err := c.ShouldBindJSON(&p); err != nil {
		// params error
		// we use the shouldBindJSON and use the binding tag on model
		// the gin can help us to verify params
		response.Send(c, code.CodeInvalidParams)
		return
	}
	// 2. logic handle
	token, statusCode := session.Create(&p)

	// 3. response result
	if statusCode != code.CodeCreated {
		response.Send(c, statusCode)
		return
	}
	response.SendWithData(c, statusCode, gin.H{
		"token": token,
	})

}

func (s *SessionController) CheckSignIn(c *gin.Context) {
	uid, exists := c.Get(constants.CtxUidKey)
	if !exists {
		response.Send(c, code.CodeNeedLogin)
		return
	}
	response.SendWithData(c, code.CodeOK, gin.H{
		"uid": uid,
	})

}
