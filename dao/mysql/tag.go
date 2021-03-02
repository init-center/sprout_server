package mysql

import (
	"sprout_server/common/transaction"
	"sprout_server/models"

	"github.com/jmoiron/sqlx"
)

func CreateTag(name string) (err error) {
	sqlStr := `INSERT INTO t_post_tag(name) VALUES(?)`
	_, err = db.Exec(sqlStr, name)
	return
}

func DeleteTagById(id uint64) (err error) {
	txFunc := func(tx *sqlx.Tx) (err error) {
		relationSqlStr := `DELETE FROM t_post_tag_relation WHERE tid = ?`
		sqlStr := `DELETE FROM t_post_tag WHERE id = ?`
		_, err = tx.Exec(relationSqlStr, id)
		_, err = tx.Exec(sqlStr, id)
		return
	}

	err = transaction.Start(db, txFunc)

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

func GetAllTags() (tags models.Tags, err error) {
	sqlStr := `SELECT id,name FROM t_post_tag`
	err = db.Select(&tags, sqlStr)
	return

}

func GetTagsByPid(pid uint64) (tags models.Tags, err error) {
	sqlStr := `SELECT pt.id, pt.name FROM t_post_tag_relation ptr LEFT JOIN t_post_tag pt ON pt.id = ptr.tid WHERE ptr.pid = ?`
	err = db.Select(&tags, sqlStr, pid)
	return
}

func DisconnectPostTagRelation(pid uint64, tid uint64) (err error) {
	deleteTagSqlStr := `DELETE FROM t_post_tag_relation WHERE pid = ? AND tid = ?`
	_, err = db.Exec(deleteTagSqlStr, pid, tid)
	return
}
