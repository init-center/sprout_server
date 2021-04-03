package mysql

import "sprout_server/models"

func CreatePageViews(uid string, url string, ip string, ua string, os string, engine string, browser string) (err error) {
	sqlStr := `INSERT INTO t_page_views(uid, url, ip, user_agent, os, engine, browser) VALUES(?, ?,?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlStr, uid, url, ip, ua, os, engine, browser)
	return
}

func GetPageViews() (pv models.PageViews, err error) {
	sqlStr := `SELECT COUNT(id) AS visits, COUNT(DISTINCT ip) AS distinctions FROM t_page_views`
	err = db.Get(&pv, sqlStr)
	return
}
