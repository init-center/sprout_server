package post

import (
	"database/sql"
	"sprout_server/common/response/code"
	"sprout_server/dao/mysql"
	"sprout_server/models"

	"go.uber.org/zap"
)

func Create(p *models.ParamsAddPost) int {
	// check the user exist
	exist, err := mysql.CheckUidExist(p.Uid)
	if err != nil {
		zap.L().Error("check uid exist failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodeUserNotExist
	}

	// check the category exist
	exist, err = mysql.CheckCategoryExistById(p.Category)
	if err != nil {
		zap.L().Error("check category exist by id failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodeCategoryNotExist
	}

	// check the tags exist
	for _, tag := range p.Tags {
		exist, err := mysql.CheckTagExistById(tag)
		if err != nil {
			zap.L().Error("check tag exist by id failed", zap.Error(err))
			return code.CodeServerBusy
		}

		if !exist {
			return code.CodeTagNotExist
		}
	}

	//category does not exist, can be created
	if err := mysql.CreatePost(p); err != nil {
		zap.L().Error("create post failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeCreated

}

func GetList(qs *models.QueryStringGetPostList) (models.PostList, int) {
	posts, err := mysql.GetPostList(qs)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get posts failed", zap.Error(err))
		return posts, code.CodeServerBusy
	}

	return posts, code.CodeOK
}

func GetDetail(p *models.UriGetPostDetail) (models.PostDetail, int) {
	post, err := mysql.GetPostDetail(p)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get post detail failed", zap.Error(err))
		return post, code.CodeServerBusy
	}

	if err == sql.ErrNoRows {
		return post, code.CodePostNotExist
	}

	return post, code.CodeOK
}
