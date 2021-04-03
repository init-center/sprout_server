package pageViews

import (
	"sprout_server/common/response/code"
	"sprout_server/dao/mysql"
	"sprout_server/models"

	"go.uber.org/zap"
)

func CreatePageViews(uid string, url string, ip string, ua string, os string, engine string, browser string) int {
	err := mysql.CreatePageViews(uid, url, ip, ua, os, engine, browser)
	if err != nil {
		zap.L().Error("create page views failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeOK

}

func GetPageViews() (pv models.PageViews, statusCode int) {
	pv, err := mysql.GetPageViews()
	if err != nil {
		zap.L().Error("get page views failed", zap.Error(err))
		return pv, code.CodeServerBusy
	}

	return pv, code.CodeOK

}
