package mysql

import (
	"sprout_server/models"
	"sprout_server/models/queryfields"
)

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

func UpdateCategory(name string, id uint64) (err error) {
	sqlStr := `UPDATE t_post_category SET name = ? WHERE id = ?`
	_, err = db.Exec(sqlStr, name, id)
	if err != nil {
		return
	}
	return
}

func GetPostCountOfCategoryId(id uint64) (count uint64, err error) {
	sql := `SELECT COUNT(id) FROM t_post WHERE category = ?`
	err = db.Get(&count, sql, id)
	return
}

func DeleteCategory(id uint64) (err error) {
	sqlStr := `DELETE FROM t_post_category WHERE id = ?`
	_, err = db.Exec(sqlStr, id)
	if err != nil {
		return
	}
	return
}

func GetCategories(queryFields *queryfields.CategoryQueryFields) (categories models.CategoryList, err error) {
	sqlStr := `SELECT c.id, c.name, COUNT(p.id) AS post_count FROM t_post_category c LEFT JOIN t_post p ON p.category = c.id `

	sqlStr = dynamicConcatCategorySql(sqlStr, queryFields)
	sqlStr += `GROUP BY c.id ORDER BY c.id, post_count DESC `
	var limit = queryFields.Limit
	if queryFields.Page != 0 && queryFields.Limit != 0 {
		sqlStr += ` LIMIT ? OFFSET ?`
		err = db.Select(&categories.List, sqlStr, queryFields.Id, queryFields.Keyword, limit, (queryFields.Page-1)*limit)
	} else {
		err = db.Select(&categories.List, sqlStr, queryFields.Id, queryFields.Keyword)
	}

	if err != nil {
		return
	}

	if len(categories.List) == 0 {
		categories.List = make([]models.CategoryData, 0, 0)
	}

	// get category count
	countSqlStr := ` SELECT COUNT(id) FROM t_post_category`
	countSqlStr = dynamicConcatCategorySql(countSqlStr, queryFields)
	err = db.Get(&categories.Page.Count, countSqlStr, queryFields.Id, queryFields.Keyword)
	if err != nil {
		return
	}
	categories.Page.CurrentPage = queryFields.Page
	categories.Page.Size = queryFields.Limit

	return

}

func dynamicConcatCategorySql(sqlStr string, queryFields *queryfields.CategoryQueryFields) string {
	if queryFields.Id != 0 {
		sqlStr += ` WHERE id = ? `
	} else {
		sqlStr += ` WHERE ? = 0 `
	}

	if queryFields.Keyword != "" {
		sqlStr += ` AND name LIKE CONCAT("%", ?, "%") `
	} else {
		sqlStr += ` AND LENGTH(?) = 0 `
	}

	return sqlStr
}
