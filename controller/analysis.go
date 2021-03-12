package controller

import (
	"sprout_server/common/response"
	"sprout_server/common/response/code"
	"sprout_server/logic/analysis"
	"sprout_server/models"

	"github.com/gin-gonic/gin"
)

type AnalysisController struct{}

func (ac *AnalysisController) GetUserAnalysis(c *gin.Context) {
	var p models.ParamsRecentDaysAnalysis
	if err := c.ShouldBindQuery(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}

	if p.Days == 0 {
		p.Days = 7
	}
	userAnalysis, statusCode := analysis.GetUserAnalysis(p.Days)

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}

	response.SendWithData(c, statusCode, userAnalysis)
}

func (ac *AnalysisController) GetCommentAnalysis(c *gin.Context) {
	var p models.ParamsRecentDaysAnalysis
	if err := c.ShouldBindQuery(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}

	if p.Days == 0 {
		p.Days = 7
	}
	commentAnalysis, statusCode := analysis.GetCommentAnalysis(p.Days)

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}

	response.SendWithData(c, statusCode, commentAnalysis)
}

func (ac *AnalysisController) GetPostAnalysis(c *gin.Context) {
	postAnalysis, statusCode := analysis.GetPostAnalysis()

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}

	response.SendWithData(c, statusCode, postAnalysis)
}

func (ac *AnalysisController) GetPostViewsRank(c *gin.Context) {
	var p models.ParamsGetViewsRank
	if err := c.ShouldBindQuery(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}

	if p.Limit == 0 {
		p.Limit = 7
	}
	postViewRank, statusCode := analysis.GetPostViewsRank(p.Limit)

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}

	response.SendWithData(c, statusCode, postViewRank)
}

func (ac *AnalysisController) GetCategoriesPostsCount(c *gin.Context) {
	list, statusCode := analysis.GetCategoriesPostsCount()

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}

	response.SendWithData(c, statusCode, list)
}

func (ac *AnalysisController) GetTagsPostsCount(c *gin.Context) {
	list, statusCode := analysis.GetTagsPostsCount()

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}

	response.SendWithData(c, statusCode, list)
}
