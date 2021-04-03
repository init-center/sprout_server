package mysql

import (
	"sprout_server/models"
	"sprout_server/models/queryfields"
)

func CreateFriend(p *models.ParamsAddFriend) (err error) {
	sqlStr := `INSERT INTO t_friends(name, url, avatar, intro) VALUES(?, ?, ?, ?)`
	_, err = db.Exec(sqlStr, p.Name, p.Url, p.Avatar, p.Intro)
	return
}

func CheckFriendExistByName(name string) (bool, error) {
	sqlStr := `SELECT count(name) FROM t_friends WHERE name = ?`
	var count int
	if err := db.Get(&count, sqlStr, name); err != nil {
		return false, err
	}
	return count > 0, nil
}

func CheckFriendExistById(id uint64) (bool, error) {
	sqlStr := `SELECT count(id) FROM t_friends WHERE id = ?`
	var count int
	if err := db.Get(&count, sqlStr, id); err != nil {
		return false, err
	}
	return count > 0, nil
}

func UpdateFriend(p *models.ParamsAddFriend, u *models.UriUpdateFriend) (err error) {
	sqlStr := `UPDATE t_friends SET name = ?, url = ?, avatar = ?, intro = ? WHERE id = ?`
	_, err = db.Exec(sqlStr, p.Name, p.Url, p.Avatar, p.Intro, u.Id)
	if err != nil {
		return
	}
	return
}

func DeleteFriend(id uint64) (err error) {
	sqlStr := `DELETE FROM t_friends WHERE id = ?`
	_, err = db.Exec(sqlStr, id)
	if err != nil {
		return
	}
	return
}

func GetFriendList(queryFields *queryfields.FriendQueryFields) (friends models.FriendList, err error) {
	sqlStr := `SELECT id, name, url, avatar, intro FROM t_friends WHERE `

	sqlStr = dynamicConcatFriendSql(sqlStr, queryFields)
	var limit = queryFields.Limit
	if queryFields.Page != 0 && queryFields.Limit != 0 {
		sqlStr += ` LIMIT ? OFFSET ?`
		err = db.Select(&friends.List, sqlStr, queryFields.Keyword, queryFields.Keyword, limit, (queryFields.Page-1)*limit)
	} else {
		err = db.Select(&friends.List, sqlStr, queryFields.Keyword, queryFields.Keyword)
	}

	if err != nil {
		return
	}

	if len(friends.List) == 0 {
		friends.List = make([]models.Friend, 0, 0)
	}

	// get friend count
	countSqlStr := ` SELECT COUNT(id) FROM t_friends WHERE `
	countSqlStr = dynamicConcatFriendSql(countSqlStr, queryFields)
	err = db.Get(&friends.Page.Count, countSqlStr, queryFields.Keyword, queryFields.Keyword)
	if err != nil {
		return
	}
	friends.Page.CurrentPage = queryFields.Page
	friends.Page.Size = queryFields.Limit

	return

}

func dynamicConcatFriendSql(sqlStr string, queryFields *queryfields.FriendQueryFields) string {

	if queryFields.Keyword != "" {
		sqlStr += ` (name LIKE CONCAT("%", ?, "%") OR intro LIKE CONCAT("%", ?, "%") )`
	} else {
		sqlStr += ` LENGTH(?) = 0 AND LENGTH(?) = 0 `
	}

	return sqlStr
}
