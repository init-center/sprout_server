package controller

import (
	"sprout_server/common/response"
	"sprout_server/common/response/code"
	"sprout_server/logic/tag"
	"sprout_server/models"
	"sprout_server/models/queryfields"

	"github.com/gin-gonic/gin"
)

type TagController struct{}

func (t *TagController) Create(c *gin.Context) {
	// 1. verify params
	var p models.ParamsAddTag
	if err := c.ShouldBindJSON(&p); err != nil {
		// params error
		// we use the shouldBindJSON and use the binding tag on model
		// the gin can help us to verify params
		response.Send(c, code.CodeInvalidParams)
		return
	}
	// 2. logic handle
	statusCode := tag.Create(&p)

	// 3. response result
	response.Send(c, statusCode)
	return

}

func (t *TagController) Update(c *gin.Context) {
	var p models.ParamsAddTag
	if err := c.ShouldBindJSON(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	var u models.UriUpdateTag
	if err := c.ShouldBindUri(&u); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	statusCode := tag.Update(&p, &u)

	response.Send(c, statusCode)
	return

}

func (t *TagController) Delete(c *gin.Context) {
	var u models.UriDeleteTag
	if err := c.ShouldBindUri(&u); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	statusCode := tag.Delete(&u)

	response.Send(c, statusCode)
	return

}

func (t *TagController) GetByQuery(c *gin.Context) {
	var p queryfields.TagQueryFields
	if err := c.ShouldBindQuery(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	categories, statusCode := tag.GetByQuery(&p)

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}

	response.SendWithData(c, statusCode, categories)
}
