package mysql

import "sprout_server/models"

func CreateTag(name string) (err error) {
	sqlStr := `INSERT INTO t_post_tag(name) VALUES(?)`
	_, err = db.Exec(sqlStr, name)
	return
}

func CheckTagExistByName(name string) (bool, error) {
	sqlStr := `SELECT count(name) FROM t_post_tag WHERE name = ?`
	var count int
	if err := db.Get(&count, sqlStr, name); err != nil {
		return false, err
	}
	return count > 0, nil
}

func CheckTagExistById(id int64) (bool, error) {
	sqlStr := `SELECT count(id) FROM t_post_tag WHERE id = ?`
	var count int
	if err := db.Get(&count, sqlStr, id); err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetAllTag() (tags models.Tags, err error) {
	sqlStr := `SELECT id,name FROM t_post_tag`
	err = db.Select(&tags, sqlStr)
	return

}
