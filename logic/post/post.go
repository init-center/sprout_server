package post

import (
	"database/sql"
	"sprout_server/common/response/code"
	"sprout_server/dao/mysql"
	"sprout_server/models"
	"sprout_server/models/queryfields"

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

	if err := mysql.CreatePost(p); err != nil {
		zap.L().Error("create post failed", zap.Error(err))
		if err.Error() == "cannot top deleted article" {
			return code.CodeCantTopDeletePost
		}
		return code.CodeServerBusy
	}

	return code.CodeCreated

}

func Update(p *models.ParamsUpdatePost, u *models.UriUpdatePost) int {

	// check the post exist
	exist, err := mysql.CheckPostExistById(u.Pid)
	if err != nil {
		zap.L().Error("check post exist by id failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodePostNotExist
	}

	// check the category exist
	if p.Category != nil {
		exist, err = mysql.CheckCategoryExistById(*p.Category)
		if err != nil {
			zap.L().Error("check category exist by id failed", zap.Error(err))
			return code.CodeServerBusy
		}

		if !exist {
			return code.CodeCategoryNotExist
		}
	}

	// check the tags exist
	if p.Tags != nil {
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
	}

	if err := mysql.UpdatePost(p, u); err != nil {
		zap.L().Error("update post failed", zap.Error(err))
		if err.Error() == "cannot top deleted article" {
			return code.CodeCantTopDeletePost
		}
		return code.CodeServerBusy
	}

	return code.CodeOK

}

func GetList(qs *models.QueryStringGetPostList) (models.PostList, int) {
	posts, err := mysql.GetPostList(qs)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get posts failed", zap.Error(err))
		return posts, code.CodeServerBusy
	}

	return posts, code.CodeOK
}

func GetListByAdmin(queryFields *queryfields.PostQueryFields) (models.PostListByAdmin, int) {
	posts, err := mysql.GetPostListByAdmin(queryFields)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get posts by admin failed", zap.Error(err))
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

func GetTopPost() (models.PostListItem, int) {
	post, err := mysql.GetTopPost()
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get top post failed", zap.Error(err))
		return post, code.CodeServerBusy
	}

	if err == sql.ErrNoRows {
		return post, code.CodePostNotExist
	}

	return post, code.CodeOK
}

func GetDetailByAdmin(p *models.UriGetPostDetail) (models.PostItemByAdmin, int) {
	post, err := mysql.GetPostDetailByAdmin(p)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get post detail by admin failed", zap.Error(err))
		return post, code.CodeServerBusy
	}

	if err == sql.ErrNoRows {
		return post, code.CodePostNotExist
	}

	return post, code.CodeOK
}
