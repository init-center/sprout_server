package controller

import (
	"sprout_server/common/response"
	"sprout_server/common/response/code"
	"sprout_server/logic/friend"
	"sprout_server/models"
	"sprout_server/models/queryfields"

	"github.com/gin-gonic/gin"
)

type FriendController struct{}

func (fc *FriendController) Create(c *gin.Context) {

	var p models.ParamsAddFriend
	if err := c.ShouldBindJSON(&p); err != nil {

		response.Send(c, code.CodeInvalidParams)
		return
	}

	statusCode := friend.Create(&p)

	response.Send(c, statusCode)
	return

}

func (fc *FriendController) Update(c *gin.Context) {
	var p models.ParamsAddFriend
	if err := c.ShouldBindJSON(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	var u models.UriUpdateFriend
	if err := c.ShouldBindUri(&u); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	statusCode := friend.Update(&p, &u)

	response.Send(c, statusCode)
	return

}

func (fc *FriendController) Delete(c *gin.Context) {
	var u models.UriDeleteFriend
	if err := c.ShouldBindUri(&u); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	statusCode := friend.Delete(&u)

	response.Send(c, statusCode)
	return

}

func (fc *FriendController) GetByQuery(c *gin.Context) {
	var p queryfields.FriendQueryFields
	if err := c.ShouldBindQuery(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	friends, statusCode := friend.GetByQuery(&p)

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}

	response.SendWithData(c, statusCode, friends)
}
