package controller

import (
	"sprout_server/common/response"
	"sprout_server/common/response/code"
	"sprout_server/logic/tag"
	"sprout_server/models"

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

func (t *TagController) GetAll(c *gin.Context) {
	tags, statusCode := tag.GetAll()

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}

	response.SendWithData(c, statusCode, tags)
}
