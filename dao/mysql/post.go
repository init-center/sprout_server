package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"sprout_server/common/snowflake"
	"sprout_server/common/tools"
	"sprout_server/common/transaction"
	"sprout_server/models"
	"sprout_server/models/queryfields"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type TopPost struct {
	TopTime *time.Time `json:"topTime" db:"top_time"`
}

func CreatePost(p *models.ParamsAddPost, uid string) (err error) {
	var pid = snowflake.GenID()

	//gen txFunc
	txFunc := func(tx *sqlx.Tx) (err error) {
		postSqlStr := `INSERT INTO t_post(uid, pid, category, title, cover, bgm, summary, content) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = tx.Exec(postSqlStr, uid, pid, p.Category, p.Title, p.Cover, p.Bgm, p.Summary, p.Content)
		if err != nil {
			return
		}

		postConfigSqlStr := `INSERT INTO t_post_config(pid, display, comment_open, top_time) VALUES(?, ?, ?, ?)`
		if p.IsTop == 1 {
			if *p.IsDelete == 1 {
				return errors.New("cannot top deleted article")
			}
			cancelTopSql := `UPDATE t_post_config SET top_time = NULL WHERE top_time IS NOT NULL;`
			_, err = tx.Exec(cancelTopSql)
			_, err = tx.Exec(postConfigSqlStr, pid, p.IsDisplay, p.IsCommentOpen, time.Now())
		} else if p.IsTop == 0 {
			_, err = tx.Exec(postConfigSqlStr, pid, p.IsDisplay, p.IsCommentOpen, nil)

			// If there is no pinned post, set the first post as a pinned article
			var tp TopPost

			err = tx.Get(&tp, `SELECT top_time FROM t_post_config WHERE top_time IS NOT NULL LIMIT 1`)
			if err != nil && err != sql.ErrNoRows {
				return
			}

			if err == sql.ErrNoRows {
				_, err = tx.Exec(`
					UPDATE t_post_config pc 
					LEFT JOIN t_post p 
					ON pc.pid = p.pid 
					SET pc.top_time = NOW() 
					WHERE p.delete_time IS NULL 
					AND pc.display = 1
					AND p.id = (SELECT MIN(id) FROM t_post)`)
				if err != nil {
					return
				}
			}

		}

		if err != nil {
			return
		}

		// store (?, ?) slice
		valueStrings := make([]string, 0, len(p.Tags))
		// store values slice
		valueArgs := make([]interface{}, 0, len(p.Tags)*2)
		// range tags to prepare data
		for _, u := range p.Tags {
			valueStrings = append(valueStrings, "(?, ?)")
			valueArgs = append(valueArgs, pid)
			valueArgs = append(valueArgs, u)
		}
		// join stmt
		postTagRelationSqlStr := fmt.Sprintf("INSERT INTO t_post_tag_relation(pid, tid) VALUES%s",
			strings.Join(valueStrings, ","))
		_, err = tx.Exec(postTagRelationSqlStr, valueArgs...)
		if err != nil {
			return
		}

		postViewsSqlStr := `INSERT INTO t_post_views(pid) VALUES(?)`
		_, err = tx.Exec(postViewsSqlStr, pid)
		return
	}

	// start a transaction
	err = transaction.Start(db, txFunc)
	return
}

func UpdatePost(p *models.ParamsUpdatePost, u *models.UriUpdatePost) (err error) {
	pid := u.Pid
	//gen txFunc
	txFunc := func(tx *sqlx.Tx) (err error) {
		postSqlStr := `UPDATE t_post SET 
		category = IFNULL(?, category), 
		title = IFNULL(?, title), 
		cover = IFNULL(?, cover),
		delete_time = CASE ? WHEN NULL THEN delete_time 
		WHEN 0 THEN NULL 
		WHEN 1 THEN NOW() 
		ELSE delete_time END,
		bgm = IFNULL(?, bgm),
		summary = IFNULL(?, summary),
		content = IFNULL(?, content) 
		WHERE pid = ?
		`
		_, err = tx.Exec(postSqlStr, p.Category, p.Title, p.Cover, p.IsDelete, p.Bgm, p.Summary, p.Content, pid)
		if err != nil {
			return
		}

		uriGetPostDetail := &models.UriGetPostDetail{
			Pid: pid,
		}

		if p.IsTop != nil && *p.IsTop == 1 {
			post, err := GetPostDetailByAdmin(uriGetPostDetail)
			if err != nil {
				return err
			}

			if (p.IsDelete != nil && *p.IsDelete == 1) || post.DeleteTime != nil {
				return errors.New("cannot top deleted article")
			}
			cancelTopSql := `UPDATE t_post_config SET top_time = NULL WHERE top_time IS NOT NULL;`
			_, err = tx.Exec(cancelTopSql)
			if err != nil {
				return err
			}
		}

		postConfigSqlStr := `UPDATE t_post_config SET 
		display = CASE ? WHEN NULL THEN display 
		WHEN 0 THEN 0 
		WHEN 1 THEN 1 
		ELSE display END, 
		comment_open = CASE ? WHEN NULL THEN comment_open 
		WHEN 0 THEN 0 
		WHEN 1 THEN 1 
		ELSE comment_open END, 
		top_time = CASE ? WHEN NULL THEN top_time 
		WHEN 0 THEN NULL 
		WHEN 1 THEN NOW() 
		ELSE top_time END 
		WHERE pid = ?`

		_, err = tx.Exec(postConfigSqlStr, p.IsDisplay, p.IsCommentOpen, p.IsTop, pid)

		if err != nil {
			return
		}

		if p.IsTop != nil && *p.IsTop == 0 {

			var tp TopPost

			err = tx.Get(&tp, `SELECT top_time FROM t_post_config WHERE top_time IS NOT NULL LIMIT 1`)
			if err != nil && err != sql.ErrNoRows {
				return
			}

			if err == sql.ErrNoRows {
				_, err = tx.Exec(`
					UPDATE t_post_config pc
					LEFT JOIN t_post p
					ON pc.pid = p.pid
					SET pc.top_time = NOW()
					WHERE p.delete_time IS NULL
					AND pc.display = 1
					AND p.id = (SELECT MIN(id) FROM t_post)`)

				if err != nil {
					return
				}
			}
		}

		if p.Tags != nil {
			existedTags, err := GetTagsByPid(pid)
			if err != nil {
				return err
			}

			// convert tags to []interface{}
			s := make([]interface{}, len(p.Tags))
			for i, v := range p.Tags {
				s[i] = v
			}

			var WillDisConnectTags []uint64
			for _, tag := range existedTags {

				includes := tools.Contains(s, tag.Id)
				if !includes {
					WillDisConnectTags = append(WillDisConnectTags, tag.Id)
				}
			}

			var WillAddTags []uint64
			for _, tid := range p.Tags {
				found := false
				for _, tag := range existedTags {
					existedTid := tag.Id
					if existedTid == tid {
						found = true
						break
					}
				}
				if !found {
					WillAddTags = append(WillAddTags, tid)
				}
			}

			// delete tags
			deleteTagSqlStr := `DELETE FROM t_post_tag_relation WHERE pid = ? AND tid = ?`
			for _, tid := range WillDisConnectTags {
				_, err = tx.Exec(deleteTagSqlStr, pid, tid)
				if err != nil {
					return err
				}
			}

			// add tags
			if len(WillAddTags) > 0 {
				// store (?, ?) slice
				valueStrings := make([]string, 0, len(WillAddTags))
				// store values slice
				valueArgs := make([]interface{}, 0, len(WillAddTags)*2)
				// range tags to prepare data
				for _, u := range WillAddTags {
					valueStrings = append(valueStrings, "(?, ?)")
					valueArgs = append(valueArgs, pid)
					valueArgs = append(valueArgs, u)
				}
				// join stmt
				postTagRelationSqlStr := fmt.Sprintf("INSERT INTO t_post_tag_relation(pid, tid) VALUES%s",
					strings.Join(valueStrings, ","))
				_, err = tx.Exec(postTagRelationSqlStr, valueArgs...)
				if err != nil {
					return err
				}
			}
		}
		return
	}

	// start a transaction
	err = transaction.Start(db, txFunc)
	return
}

func CheckPostExistById(pid uint64) (bool, error) {
	sqlStr := `SELECT count(id) FROM t_post WHERE pid = ?`
	var count int
	if err := db.Get(&count, sqlStr, pid); err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetPostListByAdmin(queryFields *queryfields.PostQueryFields) (postList models.PostListByAdmin, err error) {
	sqlStr := `
	SELECT 
    DISTINCT p.pid, 
    p.uid, 
    p.cover, 
    p.title, 
    p.summary, 
	p.content, 
    p.category, 
	p.bgm, 
	c.name AS category_name,
    p.create_time, 
	p.update_time,
	p.delete_time,
    pc.top_time,
	pc.comment_open,
	pc.display, 
	pv.views 
	FROM t_post p 
	LEFT JOIN t_post_config pc 
	ON p.pid = pc.pid 
	LEFT JOIN t_post_views pv 
	ON pc.pid = pv.pid 
	LEFT JOIN t_post_category c 
	ON p.category = c.id `

	if queryFields.Tag != 0 || queryFields.TagName != "" {
		sqlStr += `LEFT JOIN t_post_tag_relation ptr ON p.pid = ptr.pid LEFT JOIN t_post_tag pt ON ptr.tid = pt.id `
	}
	sqlStr += ` WHERE `

	sqlStr = dynamicConcatPostSql(sqlStr, queryFields)

	sqlStr += ` ORDER BY p.create_time DESC `
	var limit = queryFields.Limit
	if queryFields.Limit != 0 && queryFields.Page != 0 {
		sqlStr += ` LIMIT ? OFFSET ?`
		err = db.Select(&postList.List, sqlStr, queryFields.Pid, queryFields.Category,
			queryFields.Tag, queryFields.CategoryName, queryFields.TagName, queryFields.CreateTimeStart, queryFields.CreateTimeEnd,
			queryFields.Keyword, queryFields.Keyword, limit, (queryFields.Page-1)*limit)
	} else {
		err = db.Select(&postList.List, sqlStr, queryFields.Pid, queryFields.Category,
			queryFields.Tag, queryFields.CategoryName, queryFields.TagName, queryFields.CreateTimeStart, queryFields.CreateTimeEnd,
			queryFields.Keyword, queryFields.Keyword)
	}

	if err != nil {
		return
	}

	// get post count
	postCountSql := `
	SELECT COUNT(DISTINCT p.id) 
	FROM t_post p 
	LEFT JOIN t_post_config pc 
	ON p.pid = pc.pid 
	LEFT JOIN t_post_category c 
	ON p.category = c.id `

	if queryFields.Tag != 0 || queryFields.TagName != "" {
		postCountSql += `LEFT JOIN t_post_tag_relation ptr ON p.pid = ptr.pid LEFT JOIN t_post_tag pt ON ptr.tid = pt.id `
	}
	postCountSql += ` WHERE `

	postCountSql = dynamicConcatPostSql(postCountSql, queryFields)

	err = db.Get(&postList.Page.Count, postCountSql, queryFields.Pid, queryFields.Category,
		queryFields.Tag, queryFields.CategoryName, queryFields.TagName, queryFields.CreateTimeStart, queryFields.CreateTimeEnd,
		queryFields.Keyword, queryFields.Keyword)

	postList.Page.CurrentPage = queryFields.Page
	postList.Page.Size = queryFields.Limit

	if len(postList.List) == 0 {
		postList.List = make([]models.PostItemByAdmin, 0, 0)
		return
	}

	// get tags
	tagsSqlStr := `SELECT pt.id, pt.name FROM t_post_tag_relation ptr LEFT JOIN t_post_tag pt ON pt.id = ptr.tid WHERE ptr.pid = ?`
	for i := range postList.List {
		err = db.Select(&postList.List[i].Tags, tagsSqlStr, postList.List[i].Pid)
		if err != nil {
			return
		}

		// get favorites
		postList.List[i].Favorites, err = GetPostFavoriteCount(postList.List[i].Pid)

		// get commentCount
		postList.List[i].CommentCount, err = GetPostCommentCount(postList.List[i].Pid)
		if err != nil {
			return
		}
	}

	postList.Search = *queryFields

	return
}

func dynamicConcatPostSql(sqlStr string, queryFields *queryfields.PostQueryFields) string {
	if queryFields.IsDisplay == 0 {
		sqlStr += ` pc.display = 0 `
	} else if queryFields.IsDisplay == 1 {
		sqlStr += ` pc.display = 1 `
	} else {
		sqlStr += `1 = 1`
	}

	if queryFields.IsDelete == 0 {
		sqlStr += ` AND p.delete_time IS NULL `
	} else if queryFields.IsDelete == 1 {
		sqlStr += ` AND p.delete_time IS NOT NULL `
	}

	if queryFields.IsTop == 0 {
		sqlStr += ` AND pc.top_time IS NULL `
	} else if queryFields.IsTop == 1 {
		sqlStr += ` AND pc.top_time IS NOT NULL `
	}

	if queryFields.IsCommentOpen == 0 {
		sqlStr += ` AND pc.comment_open = 0 `
	} else if queryFields.IsCommentOpen == 1 {
		sqlStr += ` AND pc.comment_open = 1 `
	}

	if queryFields.Pid != "" {
		sqlStr += ` AND p.pid = ? `
	} else {
		// always true
		sqlStr += ` AND LENGTH(?) = 0 `
	}

	if queryFields.Category != 0 {
		sqlStr += ` AND p.category = ? `
	} else {
		sqlStr += ` AND 0 = ? `
	}

	if queryFields.Tag != 0 {
		sqlStr += ` AND ptr.tid = ? `
	} else {
		sqlStr += ` AND 0 = ? `
	}

	if queryFields.CategoryName != "" {
		sqlStr += ` AND c.name = ? `
	} else {
		sqlStr += ` AND LENGTH(?) = 0 `
	}

	if queryFields.TagName != "" {
		sqlStr += ` AND pt.name = ? `
	} else {
		sqlStr += ` AND LENGTH(?) = 0 `
	}

	if queryFields.CreateTimeStart != "" && queryFields.CreateTimeEnd != "" {
		sqlStr += ` AND (p.create_time >= ? AND p.create_time <= ?) `
	} else {
		sqlStr += ` AND LENGTH(?)= 0 AND LENGTH(?) = 0 `
	}

	if queryFields.Keyword != "" {
		sqlStr += ` AND (p.content LIKE CONCAT("%", ?, "%") OR p.title LIKE CONCAT("%", ?, "%")) `
	} else {
		sqlStr += ` AND LENGTH(?) = 0 AND LENGTH(?) = 0 `
	}

	return sqlStr
}

func GetTopPost(qs *models.QueryStringGetPostList) (topPost models.PostListItem, err error) {
	sqlStr := `
	SELECT 
    p.pid, 
    p.uid, 
    p.cover, 
    p.title, 
    p.summary, 
    p.category, 
	c.name AS category_name,
    p.create_time, 
    pc.top_time,
	pv.views 
	FROM t_post p 
	LEFT JOIN t_post_config pc 
	ON p.pid = pc.pid 
	LEFT JOIN t_post_views pv 
	ON pc.pid = pv.pid 
	LEFT JOIN t_post_category c 
	ON p.category = c.id `
	if qs.Tag != 0 || qs.TagName != "" {
		sqlStr += ` LEFT JOIN t_post_tag_relation ptr ON p.pid = ptr.pid LEFT JOIN t_post_tag pt ON ptr.tid = pt.id `
	}
	sqlStr += `WHERE pc.display = 1 
	AND p.delete_time IS NULL 
	AND pc.top_time IS NOT NULL `

	sqlStr = concatPostListSql(sqlStr, qs)

	sqlStr += `LIMIT 1`

	err = db.Get(&topPost, sqlStr, qs.Category, qs.Tag, qs.CategoryName, qs.TagName, qs.Keyword, qs.Keyword)
	if err != nil {
		return
	}
	topPost.Tags, err = GetTagsByPid(topPost.Pid)
	if err != nil {
		return
	}

	// get favorites
	topPost.Favorites, err = GetPostFavoriteCount(topPost.Pid)

	topPost.CommentCount, err = GetPostCommentCount(topPost.Pid)

	return
}

func GetTopPostDetail(qs *models.QueryStringGetPostList) (topPost models.PostDetail, err error) {
	sqlStr := `
	SELECT 
    p.uid, 
    p.cover, 
    p.title, 
    p.summary, 
	p.bgm,
	p.content,
	pc.comment_open,
	p.update_time,
    p.category, 
	c.name AS category_name,
	pv.views,
    p.create_time, 
    pc.top_time, 
    p.create_time 
	FROM t_post p 
	LEFT JOIN t_post_config pc 
	ON p.pid = pc.pid 
	LEFT JOIN t_post_views pv 
	ON pc.pid = pv.pid 
	LEFT JOIN t_post_category c 
	ON p.category = c.id `
	if qs.Tag != 0 || qs.TagName != "" {
		sqlStr += ` LEFT JOIN t_post_tag_relation ptr ON p.pid = ptr.pid LEFT JOIN t_post_tag pt ON ptr.tid = pt.id `
	}
	sqlStr += `WHERE pc.display = 1 
	AND p.delete_time IS NULL 
	AND pc.top_time IS NOT NULL `

	sqlStr = concatPostListSql(sqlStr, qs)

	sqlStr += `LIMIT 1`

	err = db.Get(&topPost, sqlStr, qs.Category, qs.Tag, qs.CategoryName, qs.TagName, qs.Keyword, qs.Keyword)
	if err != nil {
		return
	}
	topPost.Tags, err = GetTagsByPid(topPost.Pid)
	if err != nil {
		return
	}

	// get favorites
	topPost.Favorites, err = GetPostFavoriteCount(topPost.Pid)

	topPost.CommentCount, err = GetPostCommentCount(topPost.Pid)

	return
}

func GetPostList(qs *models.QueryStringGetPostList) (postList models.PostList, err error) {
	sqlStr := `
	SELECT 
    DISTINCT p.pid, 
    p.uid, 
    p.cover, 
    p.title, 
    p.summary, 
    p.category, 
	c.name AS category_name,
    p.create_time, 
    pc.top_time,
	pv.views 
	FROM t_post p 
	LEFT JOIN t_post_config pc 
	ON p.pid = pc.pid 
	LEFT JOIN t_post_views pv 
	ON pc.pid = pv.pid 
	LEFT JOIN t_post_category c 
	ON p.category = c.id `

	if qs.Tag != 0 || qs.TagName != "" {
		sqlStr += `LEFT JOIN t_post_tag_relation ptr ON p.pid = ptr.pid LEFT JOIN t_post_tag pt ON ptr.tid = pt.id `
	}

	// if the first page does not get top, subsequent requests still get it
	sqlStr += ` WHERE pc.display = 1 AND p.delete_time IS NULL `

	// if the first page get top, subsequent requests dont get it
	if qs.FirstPageGetTop == 1 {
		sqlStr += ` AND pc.top_time IS NULL `
	}

	sqlStr = concatPostListSql(sqlStr, qs)

	sqlStr += ` ORDER BY p.create_time DESC `

	var limit = qs.Limit

	var offset = (qs.Page - 1) * limit

	if qs.Page == 1 && qs.FirstPageGetTop == 1 {
		// if the page is 1 , get the top post
		var topPost models.PostListItem
		topPost, err = GetTopPost(qs)
		if err != nil && err != sql.ErrNoRows {
			return
		}
		if err == nil {
			postList.List = append(postList.List, topPost)
		}
	}

	if qs.Limit != 0 && qs.Page != 0 {
		sqlStr += ` LIMIT ? OFFSET ?`
		err = db.Select(&postList.List, sqlStr, qs.Category, qs.Tag, qs.CategoryName, qs.TagName, qs.Keyword, qs.Keyword, limit, offset)
	} else {
		err = db.Select(&postList.List, sqlStr, qs.Category, qs.Tag, qs.CategoryName, qs.TagName, qs.Keyword, qs.Keyword)
	}

	// get post count
	postCountSql := `
	SELECT COUNT(DISTINCT p.id) 
	FROM t_post p 
	LEFT JOIN t_post_config pc 
	ON p.pid = pc.pid 
	LEFT JOIN t_post_views pv 
	ON pc.pid = pv.pid 
	LEFT JOIN t_post_category c 
	ON p.category = c.id `

	if qs.Tag != 0 || qs.TagName != "" {
		postCountSql += `LEFT JOIN t_post_tag_relation ptr ON p.pid = ptr.pid LEFT JOIN t_post_tag pt ON ptr.tid = pt.id `
	}

	postCountSql += ` WHERE pc.display = 1 AND p.delete_time IS NULL `

	postCountSql = concatPostListSql(postCountSql, qs)

	err = db.Get(&postList.Page.Count, postCountSql, qs.Category, qs.Tag, qs.CategoryName, qs.TagName, qs.Keyword, qs.Keyword)

	if len(postList.List) == 0 {
		postList.List = make([]models.PostListItem, 0, 0)
		return
	}

	postList.Page.CurrentPage = qs.Page
	postList.Page.Size = qs.Limit

	// get tags
	for i := range postList.List {
		postList.List[i].Tags, err = GetTagsByPid(postList.List[i].Pid)
		if err != nil {
			return
		}

		// get favorites
		postList.List[i].Favorites, err = GetPostFavoriteCount(postList.List[i].Pid)

		// get commentCount
		postList.List[i].CommentCount, err = GetPostCommentCount(postList.List[i].Pid)
		if err != nil {
			return
		}
	}

	postList.Search = *qs

	return
}

func GetPostDetailList(qs *models.QueryStringGetPostList) (postList models.PostDetailList, err error) {
	sqlStr := `
	SELECT 
    DISTINCT p.pid,
    p.uid, 
    p.cover, 
    p.title, 
    p.summary, 
	p.bgm,
	p.content,
	pc.comment_open,
	p.update_time,
    p.category, 
	c.name AS category_name,
	pv.views,
    p.create_time, 
    pc.top_time, 
    p.create_time 
	FROM t_post p 
	LEFT JOIN t_post_config pc 
	ON p.pid = pc.pid 
	LEFT JOIN t_post_views pv 
	ON pc.pid = pv.pid 
	LEFT JOIN t_post_category c 
	ON p.category = c.id `

	if qs.Tag != 0 || qs.TagName != "" {
		sqlStr += `LEFT JOIN t_post_tag_relation ptr ON p.pid = ptr.pid LEFT JOIN t_post_tag pt ON ptr.tid = pt.id `
	}

	// if the first page does not get top, subsequent requests still get it
	sqlStr += ` WHERE pc.display = 1 AND p.delete_time IS NULL `

	// if the first page get top, subsequent requests dont get it
	if qs.FirstPageGetTop == 1 {
		sqlStr += ` AND pc.top_time IS NULL `
	}

	sqlStr = concatPostListSql(sqlStr, qs)

	sqlStr += ` ORDER BY p.create_time DESC `

	var limit = qs.Limit

	var offset = (qs.Page - 1) * limit

	if qs.Page == 1 && qs.FirstPageGetTop == 1 {
		// if the page is 1 , get the top post
		var topPost models.PostDetail
		topPost, err = GetTopPostDetail(qs)
		if err != nil && err != sql.ErrNoRows {
			return
		}
		if err == nil {
			postList.List = append(postList.List, topPost)
		}
	}

	if qs.Limit != 0 && qs.Page != 0 {
		sqlStr += ` LIMIT ? OFFSET ?`
		err = db.Select(&postList.List, sqlStr, qs.Category, qs.Tag, qs.CategoryName, qs.TagName, qs.Keyword, qs.Keyword, limit, offset)
	} else {
		err = db.Select(&postList.List, sqlStr, qs.Category, qs.Tag, qs.CategoryName, qs.TagName, qs.Keyword, qs.Keyword)
	}

	// get post count
	postCountSql := `
	SELECT COUNT(DISTINCT p.id) 
	FROM t_post p 
	LEFT JOIN t_post_config pc 
	ON p.pid = pc.pid 
	LEFT JOIN t_post_views pv 
	ON pc.pid = pv.pid 
	LEFT JOIN t_post_category c 
	ON p.category = c.id `

	if qs.Tag != 0 || qs.TagName != "" {
		postCountSql += `LEFT JOIN t_post_tag_relation ptr ON p.pid = ptr.pid LEFT JOIN t_post_tag pt ON ptr.tid = pt.id `
	}

	postCountSql += ` WHERE pc.display = 1 AND p.delete_time IS NULL `

	postCountSql = concatPostListSql(postCountSql, qs)

	err = db.Get(&postList.Page.Count, postCountSql, qs.Category, qs.Tag, qs.CategoryName, qs.TagName, qs.Keyword, qs.Keyword)

	if len(postList.List) == 0 {
		postList.List = make([]models.PostDetail, 0, 0)
		return
	}

	postList.Page.CurrentPage = qs.Page
	postList.Page.Size = qs.Limit

	// get tags
	for i := range postList.List {
		postList.List[i].Tags, err = GetTagsByPid(postList.List[i].Pid)
		if err != nil {
			return
		}

		// get favorites
		postList.List[i].Favorites, err = GetPostFavoriteCount(postList.List[i].Pid)

		// get commentCount
		postList.List[i].CommentCount, err = GetPostCommentCount(postList.List[i].Pid)
		if err != nil {
			return
		}

	}

	postList.Search = *qs

	return
}

func concatPostListSql(sqlStr string, qs *models.QueryStringGetPostList) string {

	if qs.Category != 0 {
		sqlStr += ` AND p.category = ? `
	} else {
		sqlStr += ` AND 0 = ? `
	}

	if qs.Tag != 0 {
		sqlStr += ` AND ptr.tid = ? `
	} else {
		sqlStr += ` AND 0 = ? `
	}

	if qs.CategoryName != "" {
		sqlStr += ` AND c.name = ? `
	} else {
		sqlStr += ` AND LENGTH(?) = 0 `
	}

	if qs.TagName != "" {
		sqlStr += ` AND pt.name = ? `
	} else {
		sqlStr += ` AND LENGTH(?) = 0 `
	}

	if qs.Keyword != "" {
		sqlStr += ` AND (p.content LIKE CONCAT("%", ?, "%") OR p.title LIKE CONCAT("%", ?, "%")) `
	} else {
		sqlStr += ` AND LENGTH(?) = 0 AND LENGTH(?) = 0 `
	}
	return sqlStr
}

func GetPostDetail(p *models.UriGetPostDetail) (post models.PostDetail, err error) {

	sqlStr := `
	SELECT 
    p.pid, 
    p.uid, 
    p.cover, 
    p.title, 
    p.summary, 
	p.bgm,
	p.content,
	pc.comment_open,
	p.update_time,
    p.category, 
	c.name AS category_name,
	pv.views,
    p.create_time, 
    pc.top_time, 
    p.create_time
	FROM t_post p 
	LEFT JOIN t_post_config pc 
	ON p.pid = pc.pid 
	LEFT JOIN t_post_views pv 
	ON pc.pid = pv.pid 
	LEFT JOIN t_post_category c 
	ON p.category = c.id 
	WHERE p.pid = ?
	AND pc.display = 1
	AND p.delete_time IS NULL`

	err = db.Get(&post, sqlStr, p.Pid)
	if err != nil {
		return
	}

	// get tags
	tagsSqlStr := `SELECT pt.id, pt.name FROM t_post_tag_relation ptr LEFT JOIN t_post_tag pt ON pt.id = ptr.tid WHERE ptr.pid = ?`
	err = db.Select(&post.Tags, tagsSqlStr, post.Pid)
	if err != nil {
		return
	}

	// get favorites
	post.Favorites, err = GetPostFavoriteCount(post.Pid)

	// get commentCount
	post.CommentCount, err = GetPostCommentCount(post.Pid)
	if err != nil {
		return
	}

	// add post view
	_, err = db.Exec(`UPDATE t_post_views SET views = views + 1 WHERE pid = ?`, p.Pid)
	if err != nil {
		return
	}

	return
}

func GetPostDetailByAdmin(p *models.UriGetPostDetail) (post models.PostItemByAdmin, err error) {

	sqlStr := `
	SELECT 
    p.pid, 
    p.uid, 
    p.cover, 
    p.title, 
    p.summary, 
	p.bgm,
	p.content,
	pc.comment_open,
	p.update_time,
    p.category, 
	c.name AS category_name,
	pv.views,
    p.create_time, 
    pc.top_time, 
    p.create_time,
    p.delete_time,
    pc.display 
	FROM t_post p 
	LEFT JOIN t_post_config pc 
	ON p.pid = pc.pid 
	LEFT JOIN t_post_views pv 
	ON pc.pid = pv.pid 
	LEFT JOIN t_post_category c 
	ON p.category = c.id 
	WHERE p.pid = ?`

	err = db.Get(&post, sqlStr, p.Pid)
	if err != nil {
		return
	}

	// get tags
	tagsSqlStr := `SELECT pt.id, pt.name FROM t_post_tag_relation ptr LEFT JOIN t_post_tag pt ON pt.id = ptr.tid WHERE ptr.pid = ?`
	err = db.Select(&post.Tags, tagsSqlStr, post.Pid)
	if err != nil {
		return
	}

	// get favorites
	post.Favorites, err = GetPostFavoriteCount(post.Pid)

	// get commentCount
	post.CommentCount, err = GetPostCommentCount(post.Pid)
	if err != nil {
		return
	}

	return
}
