package category

import (
	"database/sql"
	"sprout_server/common/response/code"
	"sprout_server/dao/mysql"
	"sprout_server/models"
	"sprout_server/models/queryfields"

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

func Update(p *models.ParamsAddCategory, u *models.UriUpdateCategory) int {
	// check the category exist

	exist, err := mysql.CheckCategoryExistById(u.Id)
	if err != nil {
		zap.L().Error("check category exist by id failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodeCategoryNotExist
	}

	exist, err = mysql.CheckCategoryExistByName(p.Name)
	if err != nil {
		zap.L().Error("check category exist by name failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if exist {
		return code.CodeCategoryExist
	}

	//category does not exist, can be update
	if err := mysql.UpdateCategory(p.Name, u.Id); err != nil {
		zap.L().Error("update category failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeOK

}

func Delete(u *models.UriDeleteCategory) int {

	exist, err := mysql.CheckCategoryExistById(u.Id)
	if err != nil {
		zap.L().Error("check category exist by id failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodeCategoryNotExist
	}

	count, err := mysql.GetPostCountOfCategoryId(u.Id)
	if err != nil {
		zap.L().Error("get post count of category id failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if count > 0 {
		return code.CodeCategoryHasPost
	}

	if err := mysql.DeleteCategory(u.Id); err != nil {
		zap.L().Error("delete category failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeOK

}

func GetByQuery(p *queryfields.CategoryQueryFields) (models.CategoryList, int) {
	categories, err := mysql.GetCategories(p)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get categories failed", zap.Error(err))
		return categories, code.CodeServerBusy
	}

	return categories, code.CodeOK
}
