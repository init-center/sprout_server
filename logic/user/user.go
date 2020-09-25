package user

import (
	"sprout_server/common/response/code"
	"sprout_server/dao/mysql"
	"sprout_server/dao/redis"
	"sprout_server/models"
	"sprout_server/settings"
	"strings"

	"go.uber.org/zap"
)

func Create(p *models.ParamsSignUp) int {
	// 1. check uid exist or not
	uidExist, err := mysql.CheckUidExist(p.Uid)
	if err != nil {
		zap.L().Error("check uid exist failed", zap.Error(err))
		return code.CodeServerBusy
	}
	if uidExist {
		return code.CodeUserIdExist
	}
	// 2. check name exist or not
	userNameExist, err := mysql.CheckUserNameExist(p.Name)
	if err != nil {
		zap.L().Error("check userName exist failed", zap.Error(err))
		return code.CodeServerBusy
	}
	if userNameExist {
		return code.CodeUserNameExist
	}

	// 3. check email exist or not
	emailExist, err := mysql.CheckEmailExist(p.Email)
	if err != nil {
		zap.L().Error("check email exist failed", zap.Error(err))
		return code.CodeServerBusy
	}
	if emailExist {
		return code.CodeUserNameExist
	}

	// 4. check whether the ecode has expired
	eCode, err := redis.GetECode(p.Uid)
	if err == redis.Nil {
		// no ecode or ecode expired
		return code.CodeECodeExpired
	} else if err != nil {
		// db error
		zap.L().Error("get ecode failed", zap.Error(err))
		return code.CodeServerBusy
	}

	// 5. check the ecode is equal
	isECodeEqual := strings.EqualFold(eCode, p.ECode)
	if !isECodeEqual {
		return code.CodeIncorrectECode
	}

	// 6. insert the new user to db
	u := &models.User{
		Uid:      p.Uid,
		PassWord: p.Password,
		Name:     p.Name,
		Email:    p.Email,
		Avatar:   settings.Conf.SundriesConfig.DefaultAvatar,
	}
	if err := mysql.InsertUser(u); err != nil {
		zap.L().Error("insert user to db failed", zap.Error(err))
		return code.CodeServerBusy
	}
	// success
	return code.CodeOK
}
