package mysql

import "sprout_server/models"

func AddUserFavoritePost(p *models.ParamsPostFavorite) (err error) {
	sqlStr := `INSERT INTO t_post_favorite(pid, uid) VALUES(?, ?)`
	_, err = db.Exec(sqlStr, p.Pid, p.Uid)
	return
}

func DeleteUserFavoritePost(p *models.ParamsPostFavorite) (err error) {
	sqlStr := `DELETE FROM t_post_favorite WHERE pid = ? AND uid = ?`
	_, err = db.Exec(sqlStr, p.Pid, p.Uid)
	return
}

func CheckUserFavoritePost(p *models.ParamsPostFavorite) (bool, error) {
	sqlStr := `SELECT count(id) FROM t_post_favorite WHERE uid = ? AND pid = ?`
	var count int
	if err := db.Get(&count, sqlStr, p.Uid, p.Pid); err != nil {
		return false, err
	}
	return count > 0, nil
}
