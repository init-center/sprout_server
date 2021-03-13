package controller

import (
	"sprout_server/common/constants"
	"sprout_server/common/response"
	"sprout_server/common/response/code"
	"sprout_server/logic/pageViews"

	"github.com/gin-gonic/gin"
)

type PageViewsController struct{}

func (pvc *PageViewsController) CreatePageViews(c *gin.Context) {

	uid, _ := c.Get(constants.CtxUidKey)
	ip, _ := c.Get(constants.CtxOriginIpKey)
	ua, _ := c.Get(constants.CtxOriginUAKey)
	os, _ := c.Get(constants.CtxOriginOsKey)
	engine, _ := c.Get(constants.CtxOriginEngineKey)
	browser, _ := c.Get(constants.CtxOriginBrowserKey)
	if uid == nil {
		uid = ""
	}
	if ip == nil {
		uid = ""
	}
	if ua == nil {
		ua = ""
	}
	if os == nil {
		os = ""
	}
	if engine == nil {
		engine = ""
	}
	if browser == nil {
		browser = ""
	}
	statusCode := pageViews.CreatePageViews(uid.(string), ip.(string), ua.(string), os.(string), engine.(string), browser.(string))

	response.Send(c, statusCode)
	return

}

func (pvc *PageViewsController) GetPageViews(c *gin.Context) {

	pv, statusCode := pageViews.GetPageViews()

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}

	response.SendWithData(c, statusCode, pv)
	return

}
