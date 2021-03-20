package tag

import (
	"database/sql"
	"sprout_server/common/response/code"
	"sprout_server/dao/mysql"
	"sprout_server/models"
	"sprout_server/models/queryfields"

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

func Update(p *models.ParamsAddTag, u *models.UriUpdateTag) int {
	// check the tag exist
	exist, err := mysql.CheckCategoryExistById(u.Id)
	if err != nil {
		zap.L().Error("check tag exist by id failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodeTagNotExist
	}

	exist, err = mysql.CheckTagExistByName(p.Name)
	if err != nil {
		zap.L().Error("check tag exist by name failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if exist {
		return code.CodeTagExist
	}

	//tag does not exist, can be update
	if err := mysql.UpdateTag(p.Name, u.Id); err != nil {
		zap.L().Error("update tag failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeOK

}

func Delete(u *models.UriDeleteTag) int {

	exist, err := mysql.CheckTagExistById(u.Id)
	if err != nil {
		zap.L().Error("check tag exist by id failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodeTagNotExist
	}

	count, err := mysql.GetPostCountOfTagId(u.Id)
	if err != nil {
		zap.L().Error("get post count of tag id failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if count > 0 {
		return code.CodeTagHasPost
	}

	if err := mysql.DeleteTag(u.Id); err != nil {
		zap.L().Error("delete tag failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeOK

}

func GetByQuery(p *queryfields.TagQueryFields) (models.TagList, int) {
	tags, err := mysql.GetTags(p)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get tags failed", zap.Error(err))
		return tags, code.CodeServerBusy
	}

	return tags, code.CodeOK
}
