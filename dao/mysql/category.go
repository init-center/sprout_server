package mysql

import "sprout_server/models"

func CreateCategory(name string) (err error) {
	sqlStr := `INSERT INTO t_post_category(name) VALUES(?)`
	_, err = db.Exec(sqlStr, name)
	return
}

func CheckCategoryExistByName(name string) (bool, error) {
	sqlStr := `SELECT count(name) FROM t_post_category WHERE name = ?`
	var count int
	if err := db.Get(&count, sqlStr, name); err != nil {
		return false, err
	}
	return count > 0, nil
}

func CheckCategoryExistById(id uint64) (bool, error) {
	sqlStr := `SELECT count(id) FROM t_post_category WHERE id = ?`
	var count int
	if err := db.Get(&count, sqlStr, id); err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetAllCategory() (categories models.Categories, err error) {
	sqlStr := `SELECT id,name FROM t_post_category`
	err = db.Select(&categories, sqlStr)
	return

}
