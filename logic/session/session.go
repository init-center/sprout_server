package session

import (
	"database/sql"
	"sprout_server/common/jwt"
	"sprout_server/common/response/code"
	"sprout_server/dao/mysql"
	"sprout_server/models"

	"go.uber.org/zap"
)

func Create(p *models.ParamsSignIn) (string, int) {
	isUser, err := mysql.CheckUidExist(p.Uid)
	if err != nil {
		zap.L().Error("check uid failed", zap.Error(err))
		return "", code.CodeServerBusy
	}

	if !isUser {
		isUser, err = mysql.CheckEmailExist(p.Uid)
		if err != nil {
			zap.L().Error("check email failed", zap.Error(err))
			return "", code.CodeServerBusy
		}
	}

	if !isUser {
		return "", code.CodeUserNotExist
	}

	isBaned, err := mysql.CheckUserBanStatus(p.Uid)
	if err != nil {
		zap.L().Error("check user ban status failed", zap.Error(err))
		return "", code.CodeServerBusy
	}

	if isBaned {
		return "", code.CodeUserIsBaned
	}
	u, err := mysql.Login(p)
	if err == sql.ErrNoRows {
		return "", code.CodeInvalidPassword
	}
	if err != nil {
		zap.L().Error("login failed", zap.Error(err))
		return "", code.CodeServerBusy
	}

	// verify success, gen token
	token, err := jwt.GenToken(u.Uid)
	if err != nil {
		zap.L().Error("gen token failed", zap.Error(err))
		return "", code.CodeServerBusy
	}

	return token, code.CodeCreated

}
