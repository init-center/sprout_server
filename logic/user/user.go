package user

import (
	"database/sql"
	"sprout_server/common/response/code"
	"sprout_server/dao/mysql"
	"sprout_server/dao/redis"
	"sprout_server/models"
	"sprout_server/models/queryfields"
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
		return code.CodeEmailExist
	}

	// 4. check whether the ecode has expired
	eCode, err := redis.GetECode(p.Email)
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
		UserBasicInfo: &models.UserBasicInfo{
			Uid:    p.Uid,
			Name:   p.Name,
			Avatar: settings.Conf.SundriesConfig.DefaultAvatar,
		},
		PassWord: p.Password,
		Email:    p.Email,
	}
	if err := mysql.InsertUser(u); err != nil {
		zap.L().Error("insert user to db failed", zap.Error(err))
		return code.CodeServerBusy
	}
	// success
	return code.CodeCreated
}

func GetPublicUserInfo(p models.UriGetUserInfo) (models.UserPublicInfo, int) {
	userInfo, err := mysql.GetUserPublicInfo(p.Uid)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error(" get public user info failed", zap.Error(err))
		return userInfo, code.CodeServerBusy
	}

	if err == sql.ErrNoRows {
		return userInfo, code.CodeUserNotExist
	}

	return userInfo, code.CodeOK
}

func GetPrivateUserInfo(p models.UriGetUserInfo) (models.UserPrivateInfo, int) {
	userInfo, err := mysql.GetUserPrivateInfo(p.Uid)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error(" get private user info failed", zap.Error(err))
		return userInfo, code.CodeServerBusy
	}

	if err == sql.ErrNoRows {
		return userInfo, code.CodeUserNotExist
	}

	return userInfo, code.CodeOK
}

func AdminGetUsers(p *queryfields.UserQueryFields) (models.UserDetailList, int) {
	users, err := mysql.AdminGetAllUsers(p)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error(" admin get all users failed", zap.Error(err))
		return users, code.CodeServerBusy
	}

	return users, code.CodeOK
}

func GetBanTime(uid string) (models.BanTime, int) {
	// check the user exist
	exist, err := mysql.CheckUidExist(uid)
	if err != nil {
		zap.L().Error("check uid exist by id failed", zap.Error(err))
		return models.BanTime{}, code.CodeServerBusy
	}
	if !exist {
		return models.BanTime{}, code.CodeUserNotExist
	}
	banTime, err := mysql.GetBanTime(uid)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error(" get ban time failed", zap.Error(err))
		return banTime, code.CodeServerBusy
	}

	return banTime, code.CodeOK
}

func AdminUpdateUser(p *models.ParamsAdminUpdateUser, u *models.UriUpdateUser) int {

	// check the user exist
	exist, err := mysql.CheckUidExist(u.Uid)
	if err != nil {
		zap.L().Error("check uid exist by id failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodeUserNotExist
	}

	if err := mysql.AdminUpdateUser(p, u); err != nil {
		zap.L().Error("update user failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeOK

}

func UpdateUser(p *models.ParamsUpdateUser, uid string) int {

	// check the user exist
	exist, err := mysql.CheckUidExist(uid)
	if err != nil {
		zap.L().Error("check uid exist by id failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodeUserNotExist
	}

	// check ban status
	isBaned, err := mysql.CheckUserBanStatus(uid)
	if err != nil {
		zap.L().Error("check user ban status failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if isBaned {
		return code.CodeUserIsBaned
	}

	user, err := mysql.GetUserPrivateInfo(uid)
	if err != nil && err == sql.ErrNoRows {
		return code.CodeUserNotExist
	}

	if err != nil {
		return code.CodeServerBusy
	}

	eCode, err := redis.GetECode(user.Email)
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

	if p.Password != nil && *p.Password != *p.RePassword {
		return code.CodeInvalidParams
	}

	if err := mysql.UpdateUser(p, uid); err != nil {
		zap.L().Error("update user failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeOK

}

func UpdatePassword(p *models.ParamsUpdatePassword) int {

	// check the email exist
	exist, err := mysql.CheckEmailExist(p.Email)
	if err != nil {
		zap.L().Error("check email exist failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodeEmailNotExist
	}

	eCode, err := redis.GetECode(p.Email)
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

	if err := mysql.UpdatePassword(p); err != nil {
		zap.L().Error("update user password failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeOK

}

func DeleteUser(uid string) int {

	// check the user exist
	exist, err := mysql.CheckUidExist(uid)
	if err != nil {
		zap.L().Error("check uid exist by id failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodeUserNotExist
	}

	if err := mysql.DeleteUser(uid); err != nil {
		zap.L().Error("delete user failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeOK

}

func BanUser(p *models.ParamsBanUser, u *models.UriUpdateUser) int {

	// check the user exist
	exist, err := mysql.CheckUidExist(u.Uid)
	if err != nil {
		zap.L().Error("check uid exist by id failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodeUserNotExist
	}

	if err := mysql.BanUser(p, u); err != nil {
		zap.L().Error("ban user failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeOK

}

func UnblockUser(u *models.UriUpdateUser) int {

	// check the user exist
	exist, err := mysql.CheckUidExist(u.Uid)
	if err != nil {
		zap.L().Error("check uid exist by id failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodeUserNotExist
	}

	err = mysql.UnblockUser(u)
	if err == sql.ErrNoRows {
		return code.CodeUserIsNotBaned
	} else if err != nil {
		zap.L().Error("unblock user failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeOK

}
