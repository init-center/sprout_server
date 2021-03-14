package analysis

import (
	"database/sql"
	"sprout_server/common/response/code"
	"sprout_server/dao/mysql"
	"sprout_server/models"

	"go.uber.org/zap"
)

func GetUserAnalysis(days uint8) (models.BaseAnalysisData, int) {

	userAnalysis, err := mysql.GetUserAnalysis(days)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get userAnalysis failed", zap.Error(err))
		return userAnalysis, code.CodeServerBusy
	}

	return userAnalysis, code.CodeOK
}

func GetCommentAnalysis(days uint8) (models.BaseAnalysisData, int) {

	commentAnalysis, err := mysql.GetCommentAnalysis(days)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get commentAnalysis failed", zap.Error(err))
		return commentAnalysis, code.CodeServerBusy
	}

	return commentAnalysis, code.CodeOK
}

func GetComplexAnalysis(months uint8) (models.MonthComplexIncreaseList, int) {

	complexAnalysis, err := mysql.GetComplexAnalysis(months)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get complex Analysis failed", zap.Error(err))
		return complexAnalysis, code.CodeServerBusy
	}

	return complexAnalysis, code.CodeOK
}

func GetViewsAnalysis(days uint8) (models.BaseAnalysisData, int) {

	viewsAnalysis, err := mysql.GetViewsAnalysis(days)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get viewsAnalysis failed", zap.Error(err))
		return viewsAnalysis, code.CodeServerBusy
	}

	return viewsAnalysis, code.CodeOK
}

func GetPostAnalysis() (models.PostAnalysisData, int) {

	postAnalysis, err := mysql.GetPostAnalysis()
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get postAnalysis failed", zap.Error(err))
		return postAnalysis, code.CodeServerBusy
	}

	return postAnalysis, code.CodeOK
}

func GetPostViewsRank(limit uint8) (models.PostViewsRank, int) {

	postViewRank, err := mysql.GetPostViewsRank(limit)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get post views rank failed", zap.Error(err))
		return postViewRank, code.CodeServerBusy
	}

	if len(postViewRank) == 0 {
		postViewRank = make(models.PostViewsRank, 0, 0)
	}

	return postViewRank, code.CodeOK
}

func GetCategoriesPostsCount() ([]models.CategoriesPostsCount, int) {

	list, err := mysql.GetCategoriesPostsCount()
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get categories posts count failed", zap.Error(err))
		return list, code.CodeServerBusy
	}

	if len(list) == 0 {
		list = make([]models.CategoriesPostsCount, 0, 0)
	}

	return list, code.CodeOK
}

func GetTagsPostsCount() ([]models.TagsPostsCount, int) {

	list, err := mysql.GetTagsPostsCount()
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get tags posts count failed", zap.Error(err))
		return list, code.CodeServerBusy
	}

	if len(list) == 0 {
		list = make([]models.TagsPostsCount, 0, 0)
	}

	return list, code.CodeOK
}
