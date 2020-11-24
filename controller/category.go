package controller

import (
	"sprout_server/common/response"
	"sprout_server/common/response/code"
	"sprout_server/logic/category"
	"sprout_server/models"

	"github.com/gin-gonic/gin"
)

type CategoryController struct{}

func (cc *CategoryController) Create(c *gin.Context) {
	// 1. verify params
	var p models.ParamsAddCategory
	if err := c.ShouldBindJSON(&p); err != nil {
		// params error
		// we use the shouldBindJSON and use the binding tag on model
		// the gin can help us to verify params
		response.Send(c, code.CodeInvalidParams)
		return
	}
	// 2. logic handle
	statusCode := category.Create(&p)

	// 3. response result
	response.Send(c, statusCode)
	return

}

func (cc *CategoryController) GetAll(c *gin.Context) {
	categories, statusCode := category.GetAll()

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}

	response.SendWithData(c, statusCode, categories)
}
