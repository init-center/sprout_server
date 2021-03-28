package mysql

import (
	"sprout_server/models"
	"sprout_server/models/queryfields"
)

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

func GetPostFavoriteCount(pid uint64) (count uint64, err error) {
	sqlStr := `SELECT count(id) FROM t_post_favorite WHERE pid = ?`
	if err := db.Get(&count, sqlStr, pid); err != nil {
		return 0, err
	}

	return count, nil
}

func GetFavorites(q *queryfields.FavoriteQueryFields) (favorites models.FavoritePostList, err error) {
	sqlStr := `SELECT 
	p.pid, 
    p.uid, 
    p.cover, 
    p.title, 
    p.summary, 
    p.category, 
	c.name AS category_name,
    p.create_time, 
    pc.top_time,
	pv.views,
	pf.create_time AS favorite_time 
	FROM t_post_favorite pf 
	LEFT JOIN t_post p ON pf.pid = p.pid 
	LEFT JOIN t_post_config pc 
	ON p.pid = pc.pid 
	LEFT JOIN t_post_category c ON c.id = p.category 
	LEFT JOIN t_post_views pv ON pc.pid = pv.pid 
	WHERE `
	concatSql := func(sqlStr string, qf *queryfields.FavoriteQueryFields) string {
		if qf.Uid != "" {
			sqlStr += ` pf.uid = ? `
		} else {
			sqlStr += ` LENGTH(?) = 0 `
		}

		if qf.Pid != 0 {
			sqlStr += ` AND pf.pid = ? `
		} else {
			sqlStr += ` AND ? = 0 `
		}

		return sqlStr
	}

	sqlStr = concatSql(sqlStr, q)

	sqlStr += ` ORDER BY p.create_time DESC `

	var limit = q.Limit

	if q.Limit != 0 && q.Page != 0 {
		sqlStr += ` LIMIT ? OFFSET ?`
		err = db.Select(&favorites.List, sqlStr, q.Uid, q.Pid, limit, (q.Page-1)*limit)
	} else {
		err = db.Select(&favorites.List, sqlStr, q.Uid, q.Pid)
	}

	// get post count
	postCountSql := `
	SELECT COUNT(DISTINCT pf.id) 
	FROM t_post_favorite pf WHERE `

	postCountSql = concatSql(postCountSql, q)

	err = db.Get(&favorites.Page.Count, postCountSql, q.Uid, q.Pid)
	if err != nil {
		return
	}

	favorites.Page.CurrentPage = q.Page
	favorites.Page.Size = q.Limit

	if len(favorites.List) == 0 {
		favorites.List = make([]models.FavoritePost, 0, 0)
		return
	}

	// get tags
	for i := range favorites.List {
		favorites.List[i].Tags, err = GetTagsByPid(favorites.List[i].Pid)
		if err != nil {
			return
		}

		// get favorites
		favorites.List[i].Favorites, err = GetPostFavoriteCount(favorites.List[i].Pid)

		// get commentCount
		favorites.List[i].CommentCount, err = GetPostCommentCount(favorites.List[i].Pid)
		if err != nil {
			return
		}
	}

	return

}
