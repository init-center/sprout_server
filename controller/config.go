package controller

import (
	"sprout_server/common/response"
	"sprout_server/common/response/code"
	"sprout_server/logic/config"
	"sprout_server/models"
	"sprout_server/models/queryfields"

	"github.com/gin-gonic/gin"
)

type ConfigController struct{}

func (cc *ConfigController) Create(c *gin.Context) {
	var p models.ParamsAddConfig
	if err := c.ShouldBindJSON(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	statusCode := config.Create(&p)

	response.Send(c, statusCode)
	return

}

func (cc *ConfigController) Update(c *gin.Context) {
	var p models.ParamsUpdateConfig
	if err := c.ShouldBindJSON(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	var u models.UriUpdateConfig
	if err := c.ShouldBindUri(&u); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	statusCode := config.Update(&p, &u)

	response.Send(c, statusCode)
	return

}

func (cc *ConfigController) Delete(c *gin.Context) {
	var u models.UriDeleteConfig
	if err := c.ShouldBindUri(&u); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	statusCode := config.Delete(&u)

	response.Send(c, statusCode)
	return

}

func (cc *ConfigController) GetConfigByKey(c *gin.Context) {
	var u models.UriGetConfigItem
	if err := c.ShouldBindQuery(&u); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	configItem, statusCode := config.GetItem(&u)

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}

	response.SendWithData(c, statusCode, configItem)
}

func (cc *ConfigController) GetByQuery(c *gin.Context) {
	var p queryfields.ConfigQueryFields
	if err := c.ShouldBindQuery(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	categories, statusCode := config.GetByQuery(&p)

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}

	response.SendWithData(c, statusCode, categories)
}
