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
