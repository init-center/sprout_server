package mysql

import (
	"database/sql"
	"sprout_server/common/constants"
	"sprout_server/common/pwd"
	"sprout_server/models"
	"sprout_server/models/queryfields"
	"time"
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

func GetUserGroup(uid string) (int, error) {
	sqlStr := `SELECT` + " `group`" + ` FROM t_user WHERE uid = ?`
	var group int
	if err := db.Get(&group, sqlStr, uid); err != nil {
		return 0, err
	}
	return group, nil
}

func InsertUser(user *models.User) (err error) {
	// encrypt the password
	password, err := pwd.Encrypt(user.PassWord)
	if err != nil {
		return
	}
	sqlStr := `INSERT INTO t_user(uid, name, password, email, avatar) VALUES(?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlStr, user.Uid, user.Name, password, user.Email, user.Avatar)
	return
}

func Login(p *models.ParamsSignIn) (models.User, error) {
	var u models.User
	sqlStr := `SELECT uid, name, password FROM t_user WHERE (uid=? OR email=?) AND password=?`
	password, _ := pwd.Encrypt(p.Password)
	if err := db.Get(&u, sqlStr, p.Uid, p.Uid, password); err != nil {
		return u, err
	}

	return u, nil

}

func AdminUpdateUser(p *models.ParamsAdminUpdateUser, u *models.UriUpdateUser) (err error) {
	sqlStr := `
	UPDATE t_user u SET 
	u.delete_time = CASE ? WHEN NULL THEN u.delete_time 
	WHEN 0 THEN NULL 
	WHEN 1 THEN NOW() 
	ELSE u.delete_time END,
	u.name = IFNULL(?, u.name), 
	u.avatar = IFNULL(?, u.avatar),
	u.password = IFNULL(?, u.password),
	u.email = IFNULL(?, u.email),
	u.tel = IFNULL(?, u.tel),
	u.birthday = IFNULL(?, u.birthday),
	u.intro = IFNULL(?, u.intro),
	u.group = CASE ? WHEN NULL THEN u.group 
	WHEN 1 THEN 1 
	WHEN 2 THEN 2 
	ELSE u.group END 
	WHERE u.uid = ?`

	_, err = db.Exec(sqlStr, p.IsDelete, p.Name, p.Avatar, p.Password, p.Email, p.Tel, p.Birthday, p.Intro, p.Group, u.Uid)
	if err != nil {
		return
	}
	return
}

func UpdateUser(p *models.ParamsUpdateUser, uid string) (err error) {
	sqlStr := `
	UPDATE t_user u SET 
	u.name = IFNULL(?, u.name), 
	u.avatar = IFNULL(?, u.avatar),
	u.password = IFNULL(?, u.password),
	u.email = IFNULL(?, u.email),
	u.tel = IFNULL(?, u.tel),
	u.birthday = IFNULL(?, u.birthday),
	u.intro = IFNULL(?, u.intro) 
	WHERE u.uid = ?`

	_, err = db.Exec(sqlStr, p.Name, p.Avatar, p.Password, p.Email, p.Tel, p.Birthday, p.Intro, uid)
	if err != nil {
		return
	}
	return
}

func DeleteUser(uid string) (err error) {
	sqlStr := `
	UPDATE t_user u SET 
	u.delete_time = NOW() 
	WHERE u.uid = ?`

	_, err = db.Exec(sqlStr, uid)
	if err != nil {
		return
	}
	return
}

func BanUser(p *models.ParamsBanUser, u *models.UriUpdateUser) (err error) {
	userBanExistSql := `SELECT COUNT(id) FROM t_user_ban WHERE uid = ?`
	var count uint8
	err = db.Get(&count, userBanExistSql, u.Uid)
	if err != nil {
		return
	}
	if count == 0 {
		_, err = db.Exec(`INSERT INTO t_user_ban(uid, start_time, end_time) VALUES(?, ?, ?)`, u.Uid, p.BanStartTime, p.BanEndTime)
		return err
	}

	_, err = db.Exec(`UPDATE t_user_ban SET start_time = ?, end_time = ? WHERE uid = ?`, p.BanStartTime, p.BanEndTime, u.Uid)
	return err
}

func UnblockUser(u *models.UriUpdateUser) (err error) {
	endTimeSql := `SELECT end_time, IFNULL(end_time > NOW(), 0) AS is_baned FROM t_user_ban WHERE uid = ?`
	type BanInfo struct {
		EndTime *time.Time `db:"end_time"`
		IsBaned uint8      `db:"is_baned"`
	}
	var banInfo BanInfo
	err = db.Get(&banInfo, endTimeSql, u.Uid)
	if err != nil {
		return
	}

	if banInfo.EndTime == nil || banInfo.IsBaned == 0 {
		return sql.ErrNoRows
	}

	_, err = db.Exec(`UPDATE t_user_ban SET start_time = NULL, end_time = NULL WHERE uid = ?`, u.Uid)
	return err
}

func GetUserPublicInfo(uid string) (userInfo models.UserPublicInfo, err error) {
	sqlStr := `SELECT 
	u.uid,
	u.name,
	u.avatar,
	u.group,
	u.intro,
	u.create_time,
	IFNULL(ub.end_time > NOW(), 0) AS is_baned 
	FROM t_user u 
	LEFT JOIN t_user_ban ub 
	ON u.uid = ub.uid 
	WHERE u.uid = ? 
	AND u.delete_time IS NULL 
	`

	err = db.Get(&userInfo, sqlStr, uid)

	return
}

func GetUserPrivateInfo(uid string) (userInfo models.UserPrivateInfo, err error) {
	sqlStr := `
	SELECT 
	u.uid,
	u.name,
	u.avatar,
	u.email,
	u.group,
	u.tel,
	u.gender,
	u.birthday,
	u.intro,
	u.create_time,
	u.update_time,
	u.delete_time,
	ub.start_time AS ban_start_time,
	ub.end_time AS ban_end_time, 
	IFNULL(ub.end_time > NOW(), 0) AS is_baned 
	FROM t_user u 
	LEFT JOIN t_user_ban ub 
	ON u.uid = ub.uid 
	WHERE u.uid = ? 
	AND u.delete_time IS NULL 
	`

	err = db.Get(&userInfo, sqlStr, uid)

	return
}

func AdminGetAllUsers(queryFields *queryfields.UserQueryFields) (users models.UserDetailList, err error) {
	sqlStr := `
	SELECT 
	u.uid,
	u.name,
	u.avatar,
	u.email,
	u.group,
	u.tel,
	u.gender,
	u.birthday,
	u.intro,
	u.create_time,
	u.update_time,
	u.delete_time,
	ub.start_time AS ban_start_time,
	ub.end_time AS ban_end_time, 
	IFNULL(ub.end_time > NOW(), 0) AS is_baned 
	FROM t_user u 
	LEFT JOIN t_user_ban ub 
	ON u.uid = ub.uid 
	WHERE `

	sqlStr = dynamicConcatUserSql(sqlStr, queryFields)
	sqlStr += `ORDER BY u.create_time DESC `

	var limit = queryFields.Limit

	if queryFields.Limit != 0 && queryFields.Page != 0 {
		sqlStr += `LIMIT ? OFFSET ?`
		err = db.Select(&users.List, sqlStr, queryFields.Uid, queryFields.Name,
			queryFields.Email, queryFields.Tel,
			queryFields.CreateTimeStart, queryFields.CreateTimeEnd,
			queryFields.BanTimeStart, queryFields.BanTimeStart,
			limit, (queryFields.Page-1)*limit)
	} else {
		err = db.Select(&users.List, sqlStr, queryFields.Uid, queryFields.Name,
			queryFields.Email, queryFields.Tel,
			queryFields.CreateTimeStart, queryFields.CreateTimeEnd,
			queryFields.BanTimeStart, queryFields.BanTimeStart)
	}
	if err != nil {
		return
	}

	countSqlStr := `
	SELECT COUNT(DISTINCT u.id) 
	FROM t_user u 
	LEFT JOIN t_user_ban ub 
	ON u.uid = ub.uid 
	WHERE `
	countSqlStr = dynamicConcatUserSql(countSqlStr, queryFields)
	err = db.Get(&users.Page.Count, countSqlStr, queryFields.Uid, queryFields.Name,
		queryFields.Email, queryFields.Tel,
		queryFields.CreateTimeStart, queryFields.CreateTimeEnd,
		queryFields.BanTimeStart, queryFields.BanTimeStart)
	if err != nil {
		return
	}

	users.Page.CurrentPage = queryFields.Page
	users.Page.Size = queryFields.Limit

	if len(users.List) == 0 {
		users.List = make([]models.UserDetailByAdmin, 0, 0)
	}

	return

}

func dynamicConcatUserSql(sqlStr string, queryFields *queryfields.UserQueryFields) string {
	if queryFields.Uid == "" {
		sqlStr += ` LENGTH(?) = 0 `
	} else {
		sqlStr += ` u.uid = ? `
	}

	if queryFields.Name == "" {
		sqlStr += ` AND LENGTH(?) = 0 `
	} else {
		sqlStr += ` AND u.name = ? `
	}

	if queryFields.Email == "" {
		sqlStr += ` AND LENGTH(?) = 0 `
	} else {
		sqlStr += ` AND u.email = ? `
	}

	if queryFields.Tel == "" {
		sqlStr += ` AND LENGTH(?) = 0 `
	} else {
		sqlStr += ` AND u.tel = ? `
	}

	if queryFields.Gender == 0 {
		sqlStr += ` AND u.gender = 0 `
	} else if queryFields.Gender == 1 {
		sqlStr += ` AND u.delete_time = 1 `
	}

	if queryFields.IsDelete == 0 {
		sqlStr += ` AND u.delete_time IS NULL `
	} else if queryFields.IsDelete == 1 {
		sqlStr += ` AND u.delete_time IS NOT NULL `
	}

	if queryFields.IsBaned == 0 {
		sqlStr += ` AND (ub.end_time IS NULL OR ub.end_time < NOW() ) `
	} else if queryFields.IsBaned == 1 {
		sqlStr += ` AND ub.end_time > Now() `
	}

	if queryFields.Group == constants.UserGroupAdmin {
		sqlStr += ` AND u.group = 1 `
	} else if queryFields.Group == constants.UserGroupDefault {
		sqlStr += ` AND u.group = 2 `
	}

	if queryFields.CreateTimeStart != "" && queryFields.CreateTimeEnd != "" {
		sqlStr += ` AND (u.create_time >= ? AND u.create_time <= ?) `
	} else {
		sqlStr += ` AND LENGTH(?)= 0 AND LENGTH(?) = 0 `
	}

	if queryFields.BanTimeStart != "" && queryFields.BanTimeEnd != "" {
		sqlStr += ` AND (ub.start_time >= ? AND ub.end_time <= ?) `
	} else {
		sqlStr += ` AND LENGTH(?)= 0 AND LENGTH(?) = 0 `
	}

	return sqlStr
}
