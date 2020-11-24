package category

import (
	"database/sql"
	"sprout_server/common/response/code"
	"sprout_server/dao/mysql"
	"sprout_server/models"

	"go.uber.org/zap"
)

func Create(p *models.ParamsAddCategory) int {
	// check the category exist
	exist, err := mysql.CheckCategoryExistByName(p.Name)
	if err != nil {
		zap.L().Error("check category exist by name failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if exist {
		return code.CodeCategoryExist
	}

	//category does not exist, can be created
	if err := mysql.CreateCategory(p.Name); err != nil {
		zap.L().Error("create category failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeCreated

}

func GetAll() (models.Categories, int) {
	categories, err := mysql.GetAllCategory()
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get categories failed", zap.Error(err))
		return categories, code.CodeServerBusy
	}

	// if get the empty slice, wanna an empty array but not a null in frontend
	if len(categories) == 0 {
		categories = make(models.Categories, 0, 0)
	}

	return categories, code.CodeOK
}
