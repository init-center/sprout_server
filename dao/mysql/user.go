package mysql

import (
	"sprout_server/common/pwd"
	"sprout_server/models"
	"sprout_server/settings"
)

func CheckUidExist(uid string) (bool, error) {
	sqlStr := `SELECT count(uid) FROM t_user WHERE uid = ?`
	var count int
	if err := db.Get(&count, sqlStr, uid); err != nil {
		return false, err
	}
	return count > 0, nil
}

func CheckUserNameExist(userName string) (bool, error) {
	sqlStr := `SELECT count(name) FROM t_user WHERE name = ?`
	var count int
	if err := db.Get(&count, sqlStr, userName); err != nil {
		return false, err
	}
	return count > 0, nil
}

func CheckEmailExist(email string) (bool, error) {
	sqlStr := `SELECT count(email) FROM t_user WHERE email = ?`
	var count int
	if err := db.Get(&count, sqlStr, email); err != nil {
		return false, err
	}
	return count > 0, nil
}

func InsertUser(user *models.User) (err error) {
	// encrypt the password
	password, err := pwd.Encrypt(user.PassWord, settings.Conf.SaltPrefix+user.Uid)
	if err != nil {
		return
	}
	sqlStr := `INSERT INTO t_user(uid, name, password, email, avatar) VALUES(?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlStr, user.Uid, user.Name, password, user.Email, user.Avatar)
	return
}
