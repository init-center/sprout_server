package mysql

import (
	"sprout_server/models"
	"sprout_server/models/queryfields"
)

func CreateTag(name string) (err error) {
	sqlStr := `INSERT INTO t_post_tag(name) VALUES(?)`
	_, err = db.Exec(sqlStr, name)
	return
}

func GetPostCountOfTagId(id uint64) (count uint64, err error) {
	sql := `SELECT COUNT(id) FROM t_post_tag_relation WHERE tid = ?`
	err = db.Get(&count, sql, id)
	return
}

func DeleteTag(id uint64) (err error) {
	sqlStr := `DELETE FROM t_post_tag WHERE id = ?`
	_, err = db.Exec(sqlStr, id)
	if err != nil {
		return
	}
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

func CheckTagExistById(id uint64) (bool, error) {
	sqlStr := `SELECT count(id) FROM t_post_tag WHERE id = ?`
	var count int
	if err := db.Get(&count, sqlStr, id); err != nil {
		return false, err
	}
	return count > 0, nil
}

func UpdateTag(name string, id uint64) (err error) {
	sqlStr := `UPDATE t_post_tag SET name = ? WHERE id = ?`
	_, err = db.Exec(sqlStr, name, id)
	if err != nil {
		return
	}
	return
}

func GetTagsByPid(pid uint64) (tags models.Tags, err error) {
	sqlStr := `SELECT pt.id, pt.name FROM t_post_tag_relation ptr LEFT JOIN t_post_tag pt ON pt.id = ptr.tid WHERE ptr.pid = ?`
	err = db.Select(&tags, sqlStr, pid)
	return
}

func GetTags(queryFields *queryfields.TagQueryFields) (tags models.TagList, err error) {
	sqlStr := `SELECT t.id, t.name, COUNT(ptr.id) AS post_count FROM t_post_tag t LEFT JOIN t_post_tag_relation ptr ON t.id = ptr.tid `

	sqlStr = dynamicConcatTagSql(sqlStr, queryFields)
	sqlStr += ` GROUP BY t.id ORDER BY t.id, post_count DESC `

	var limit = queryFields.Limit
	if queryFields.Page != 0 && queryFields.Limit != 0 {
		sqlStr += `LIMIT ? OFFSET ? `
		err = db.Select(&tags.List, sqlStr, queryFields.Id, queryFields.Keyword, limit, (queryFields.Page-1)*limit)
	} else {
		err = db.Select(&tags.List, sqlStr, queryFields.Id, queryFields.Keyword)
	}

	if err != nil {
		return
	}

	if len(tags.List) == 0 {
		tags.List = make([]models.TagData, 0, 0)
	}

	// get tag count
	countSqlStr := ` SELECT COUNT(id) FROM t_post_tag`
	countSqlStr = dynamicConcatTagSql(countSqlStr, queryFields)
	err = db.Get(&tags.Page.Count, countSqlStr, queryFields.Id, queryFields.Keyword)
	if err != nil {
		return
	}
	tags.Page.CurrentPage = queryFields.Page
	tags.Page.Size = queryFields.Limit

	return

}

func dynamicConcatTagSql(sqlStr string, queryFields *queryfields.TagQueryFields) string {
	if queryFields.Id != 0 {
		sqlStr += ` WHERE t.id = ? `
	} else {
		sqlStr += ` WHERE ? = 0 `
	}

	if queryFields.Keyword != "" {
		sqlStr += ` AND t.name LIKE CONCAT("%", ?, "%") `
	} else {
		sqlStr += ` AND LENGTH(?) = 0 `
	}

	return sqlStr
}
