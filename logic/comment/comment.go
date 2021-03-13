package comment

import (
	"database/sql"
	"sprout_server/common/response/code"
	"sprout_server/dao/mysql"
	"sprout_server/models"
	"sprout_server/models/queryfields"

	"go.uber.org/zap"
)

func CreatePostComment(p *models.ParamsAddComment) int {
	// check the post exist
	exist, err := mysql.CheckPostExistById(p.Pid)
	if err != nil {
		zap.L().Error("check post exist by id failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodePostNotExist
	}

	// check the reply target exist
	if p.TargetCid != 0 {
		exist, err = mysql.CheckPostCommentExist(p.TargetCid)
		if err != nil {
			zap.L().Error("check post comment exist failed", zap.Error(err))
		}

		if !exist {
			return code.CodeCommentNotExist
		}
	}

	//to create
	if err := mysql.CreatePostComment(p); err != nil {
		zap.L().Error("create post comment failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeCreated
}

func AdminUpdatePostComment(p *models.ParamsAdminUpdateComment, u *models.UriUpdateComment) int {
	// check the comment exist
	exist, err := mysql.CheckPostCommentExist(u.Cid)
	if err != nil {
		zap.L().Error("check comment exist failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodePostNotExist
	}

	//to update
	if err := mysql.AdminUpdatePostComment(p, u); err != nil {
		zap.L().Error("update post comment failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeOK
}

func GetPostCommentList(p *models.ParamsGetCommentList) (models.CommentList, int) {
	commentList, err := mysql.GetPostCommentList(p)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get post commentList failed", zap.Error(err))
		return commentList, code.CodeServerBusy
	}

	return commentList, code.CodeOK
}

func GetPostComments(p *queryfields.CommentQueryFields) (models.CommentItemListByAdmin, int) {
	comments, err := mysql.GetPostComments(p)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get post comments failed", zap.Error(err))
		return comments, code.CodeServerBusy
	}

	return comments, code.CodeOK
}

func GetPostParentCommentChildren(p *models.ParamsGetParentCommentChildren) (models.ParentCommentChildren, int) {
	parentCommentChildren, err := mysql.GetPostParentCommentChildren(p)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get post parent comment children failed", zap.Error(err))
		return parentCommentChildren, code.CodeServerBusy
	}

	return parentCommentChildren, code.CodeOK
}