package tag

import (
	"database/sql"
	"sprout_server/common/response/code"
	"sprout_server/dao/mysql"
	"sprout_server/models"

	"go.uber.org/zap"
)

func Create(p *models.ParamsAddTag) int {
	// check the tag exist
	exist, err := mysql.CheckTagExistByName(p.Name)
	if err != nil {
		zap.L().Error("check tag exist failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if exist {
		return code.CodeTagExist
	}

	//tag does not exist, can be created
	if err := mysql.CreateTag(p.Name); err != nil {
		zap.L().Error("create tag failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeCreated

}

func GetAll() (models.Tags, int) {
	tags, err := mysql.GetAllTags()
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get tags failed", zap.Error(err))
		return tags, code.CodeServerBusy
	}

	// if get the empty slice, wanna an empty array but not a null in frontend
	if len(tags) == 0 {
		tags = make(models.Tags, 0, 0)
	}

	return tags, code.CodeOK
}
