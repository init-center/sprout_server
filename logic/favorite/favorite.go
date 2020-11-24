package favorite

import (
	"sprout_server/common/response/code"
	"sprout_server/dao/mysql"
	"sprout_server/models"

	"go.uber.org/zap"
)

func CheckUserFavoritePost(p *models.ParamsPostFavorite) int {
	exist, err := mysql.CheckUserFavoritePost(p)
	if err != nil {
		zap.L().Error("check user favorite post failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if exist {
		return code.CodeAlreadyFavorited
	}

	return code.CodeNotFavorited

}

func AddUserFavoritePost(p *models.ParamsPostFavorite) int {
	exist, err := mysql.CheckUserFavoritePost(p)
	if err != nil {
		zap.L().Error("check user favorite post failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if exist {
		return code.CodeAlreadyFavorited
	}

	exist, err = mysql.CheckPostExistById(p.Pid)
	if err != nil {
		zap.L().Error("check post exist by id failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodePostNotExist
	}

	if err := mysql.AddUserFavoritePost(p); err != nil {
		zap.L().Error("add user favorite post failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeCreated
}

func DeleteUserFavoritePost(p *models.ParamsPostFavorite) int {
	exist, err := mysql.CheckUserFavoritePost(p)
	if err != nil {
		zap.L().Error("check user favorite post failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodeNotFavorited
	}

	if err := mysql.DeleteUserFavoritePost(p); err != nil {
		zap.L().Error("delete user favorite post failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeOK
}
