package config

import (
	"database/sql"
	"sprout_server/common/response/code"
	"sprout_server/dao/mysql"
	"sprout_server/models"
	"sprout_server/models/queryfields"

	"go.uber.org/zap"
)

func Create(p *models.ParamsAddConfig) int {
	// check the Config exist
	exist, err := mysql.CheckConfigExistByKey(p.Key)
	if err != nil {
		zap.L().Error("check Config exist by key failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if exist {
		return code.CodeConfigKeyExist
	}

	//Config does not exist, can be created
	if err := mysql.CreateConfig(p); err != nil {
		zap.L().Error("create Config failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeCreated

}

func Update(p *models.ParamsUpdateConfig, u *models.UriUpdateConfig) int {
	// check the Config exist

	// u.key is the oldKey, if update it must be exist
	exist, err := mysql.CheckConfigExistByKey(u.Key)
	if err != nil {
		zap.L().Error("check Config exist by key failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodeConfigKeyNotExist
	}

	exist, err = mysql.CheckConfigExistByKey(p.Key)
	if err != nil {
		zap.L().Error("check Config exist by key failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if exist {
		return code.CodeConfigKeyExist
	}

	//new config key does not exist, can be update
	if err := mysql.UpdateConfig(p, u); err != nil {
		zap.L().Error("update Config failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeOK

}

func Delete(u *models.UriDeleteConfig) int {

	if err := mysql.DeleteConfig(u.Key); err != nil {
		zap.L().Error("delete Config failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeOK

}

func GetItem(u *models.UriGetConfigItem) (models.ConfigItem, int) {
	config, err := mysql.GetConfigByKey(u.Key)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get config by key failed", zap.Error(err))
		return config, code.CodeServerBusy
	}

	return config, code.CodeOK
}

func GetByQuery(p *queryfields.ConfigQueryFields) (models.ConfigList, int) {
	configs, err := mysql.GetConfigs(p)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get configs failed", zap.Error(err))
		return configs, code.CodeServerBusy
	}

	return configs, code.CodeOK
}
