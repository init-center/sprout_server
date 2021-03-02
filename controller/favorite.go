package controller

import (
	"sprout_server/common/constants"
	"sprout_server/common/response"
	"sprout_server/common/response/code"
	"sprout_server/logic/favorite"
	"sprout_server/models"

	"github.com/gin-gonic/gin"
)

type FavoriteController struct{}

func (fc *FavoriteController) CheckUserFavoritePost(c *gin.Context) {
	// 1. verify params
	var p models.ParamsPostFavorite
	if err := c.ShouldBindUri(&p); err != nil {
		// params error
		response.Send(c, code.CodeInvalidParams)
		return
	}
	// 2. logic handle
	statusCode := favorite.CheckUserFavoritePost(&p)

	// 3. response result
	response.Send(c, statusCode)
	return

}

func (fc *FavoriteController) AddUserFavoritePost(c *gin.Context) {
	// 1. verify params
	var p models.ParamsPostFavorite
	if err := c.ShouldBindUri(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	// 2. logic handle
	// get uid from context
	uid, exist := c.Get(constants.CtxUidKey)
	if !exist {
		response.Send(c, code.CodeInvalidToken)
		return
	}
	p.Uid = uid.(string)
	statusCode := favorite.AddUserFavoritePost(&p)

	// 3. response result
	response.Send(c, statusCode)
	return
}

func (fc *FavoriteController) DeleteUserFavoritePost(c *gin.Context) {
	// 1. verify params
	var p models.ParamsPostFavorite
	if err := c.ShouldBindUri(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	// 2. logic handle
	// get uid from context
	uid, exist := c.Get(constants.CtxUidKey)
	if !exist {
		response.Send(c, code.CodeInvalidToken)
		return
	}
	p.Uid = uid.(string)
	statusCode := favorite.DeleteUserFavoritePost(&p)

	// 3. response result
	response.Send(c, statusCode)
	return
}
