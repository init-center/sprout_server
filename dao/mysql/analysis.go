package mysql

import (
	"fmt"
	"sprout_server/models"
)

func GetUserAnalysis(days uint8) (userAnalysis models.BaseAnalysisData, err error) {
	totalSql := `SELECT COUNT(id) FROM t_user`
	err = db.Get(&userAnalysis.Total, totalSql)
	if err != nil {
		return
	}

	recentIncreaseListSql := createRecentDaysIncreaseSql(days, "t_user")

	err = db.Select(&userAnalysis.RecentIncreaseList, recentIncreaseListSql)
	if err != nil {
		return
	}

	err = db.Get(&userAnalysis.TodayIncrease, `SELECT COUNT(id) AS today_increase FROM t_user WHERE DATE(create_time) = CURDATE()`)
	return
}

func GetCommentAnalysis(days uint8) (commentAnalysis models.BaseAnalysisData, err error) {
	totalSql := `SELECT COUNT(id) FROM t_post_comment`
	err = db.Get(&commentAnalysis.Total, totalSql)
	if err != nil {
		return
	}

	recentIncreaseListSql := createRecentDaysIncreaseSql(days, "t_post_comment")

	err = db.Select(&commentAnalysis.RecentIncreaseList, recentIncreaseListSql)
	if err != nil {
		return
	}

	err = db.Get(&commentAnalysis.TodayIncrease, `SELECT COUNT(id) AS today_increase FROM t_post_comment WHERE DATE(create_time) = CURDATE()`)
	return
}

func GetComplexAnalysis(months uint8) (complexAnalysis models.MonthComplexIncreaseList, err error) {

	sql := `SELECT a.month,IFNULL(b.count,0) AS comments, IFNULL(c.count,0) AS users, IFNULL(d.count,0) AS views from (`
	for i := 0; i < int(months); i++ {
		if i != 0 {
			sql += ` UNION ALL `
		}
		sql += fmt.Sprint(`SELECT DATE_FORMAT((CURDATE() - INTERVAL `, i, ` MONTH), '%Y-%m')`, ` AS month `)
	}
	sql += `) a LEFT JOIN (SELECT DATE_FORMAT(create_time, '%Y-%m') AS month, COUNT(*) AS count FROM t_post_comment GROUP BY month) b 
			ON a.month = b.month LEFT JOIN (SELECT DATE_FORMAT(create_time, '%Y-%m') AS month, COUNT(*) AS count FROM t_user GROUP BY month) c ON a.month = c.month 
			LEFT JOIN (SELECT DATE_FORMAT(create_time, '%Y-%m') AS month, COUNT(*) AS count FROM t_page_views GROUP BY month) d ON a.month = d.month ORDER BY a.month;`

	err = db.Select(&complexAnalysis, sql)
	return
}

func GetViewsAnalysis(days uint8) (viewsAnalysis models.BaseAnalysisData, err error) {
	totalSql := `SELECT COUNT(id) FROM t_page_views`
	err = db.Get(&viewsAnalysis.Total, totalSql)
	if err != nil {
		return
	}

	recentIncreaseListSql := createRecentDaysIncreaseSql(days, "t_page_views")

	err = db.Select(&viewsAnalysis.RecentIncreaseList, recentIncreaseListSql)
	if err != nil {
		return
	}

	err = db.Get(&viewsAnalysis.TodayIncrease, `SELECT COUNT(id) AS today_increase FROM t_page_views WHERE DATE(create_time) = CURDATE()`)
	return
}

func GetPostAnalysis() (postAnalysis models.PostAnalysisData, err error) {
	totalSql := `SELECT COUNT(id) FROM t_post`
	err = db.Get(&postAnalysis.Total, totalSql)
	if err != nil {
		return
	}

	MonthAverageSql := `SELECT IFNULL(ROUND(AVG(a.count)), 0) AS average FROM (SELECT COUNT(id) AS count FROM t_post GROUP BY YEAR(create_time), MONTH(create_time)) a`

	err = db.Get(&postAnalysis.Average, MonthAverageSql)
	if err != nil {
		return
	}

	err = db.Get(&postAnalysis.MonthIncrease, `SELECT COUNT(id) FROM t_post WHERE YEAR(create_time) = YEAR(NOW()) AND MONTH(create_time) = MONTH(NOW())`)
	return
}

func GetPostViewsRank(limit uint8) (postViewsRank models.PostViewsRank, err error) {

	sql := `SELECT p.pid, p.title, pv.views FROM t_post p LEFT JOIN t_post_views pv on p.pid = pv.pid ORDER BY pv.views DESC LIMIT ?`

	err = db.Select(&postViewsRank, sql, limit)
	return
}

func GetCategoriesPostsCount() (countList []models.CategoriesPostsCount, err error) {

	sql := `SELECT pc.name,COUNT(p.id) AS value FROM t_post p LEFT JOIN t_post_category pc ON p.category = pc.id GROUP BY p.category`

	err = db.Select(&countList, sql)
	return
}

func GetTagsPostsCount() (countList []models.TagsPostsCount, err error) {

	sql := `SELECT pt.name,COUNT(p.id) AS value FROM t_post p LEFT JOIN t_post_tag_relation ptr ON p.pid = ptr.pid LEFT JOIN t_post_tag pt ON ptr.tid = pt.id GROUP BY pt.name`

	err = db.Select(&countList, sql)
	return
}

func createRecentDaysIncreaseSql(days uint8, tableName string) string {
	sql := `SELECT DATE_FORMAT(a.date, '%Y-%m-%d') AS date,IFNULL(b.count,0) AS increase from (`
	for i := 0; i < int(days); i++ {
		if i != 0 {
			sql += ` UNION ALL `
		}
		sql += fmt.Sprint(`SELECT DATE_SUB(CURDATE(), INTERVAL `, i, ` DAY)`, ` AS date `)
	}
	sql += `) a LEFT JOIN (SELECT DATE(create_time) AS date, COUNT(*) AS count FROM ` + tableName + ` GROUP BY DATE(create_time)) b 
			ON a.date = b.date ORDER BY a.date;`
	return sql
}
