package mysql

import (
	"database/sql"
	"fmt"
	"sprout_server/common/snowflake"
	"sprout_server/common/transaction"
	"sprout_server/models"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

func CreatePost(p *models.ParamsAddPost) (err error) {
	var pid = snowflake.GenID()

	//gen txFunc
	txFunc := func(tx *sqlx.Tx) (err error) {
		postSqlStr := `INSERT INTO t_post(uid, pid, category, title, cover, bgm, summary, content) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = tx.Exec(postSqlStr, p.Uid, pid, p.Category, p.Title, p.Cover, p.Bgm, p.Summary, p.Content)
		if err != nil {
			return
		}

		postConfigSqlStr := `INSERT INTO t_post_config(pid, display, comment_open, top_time) VALUES(?, ?, ?, ?)`
		if p.Top == 1 {
			_, err = tx.Exec(postConfigSqlStr, pid, p.Display, p.CommentOpen, time.Now())
		} else {
			_, err = tx.Exec(postConfigSqlStr, pid, p.Display, p.CommentOpen, nil)
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

func CheckPostExistById(pid int64) (bool, error) {
	sqlStr := `SELECT count(id) FROM t_post WHERE pid = ?`
	var count int
	if err := db.Get(&count, sqlStr, pid); err != nil {
		return false, err
	}
	return count > 0, nil
}

func getTopPost() (topPost models.PostListItem, err error) {
	sqlStr := `
	SELECT 
    p.pid, 
    p.uid, 
    p.cover, 
    p.title, 
    p.summary, 
    p.category, 
	c.name,
    p.create_time, 
    pc.top_time,
	pv.views 
	FROM t_post p 
	LEFT JOIN t_post_config pc 
	ON p.pid = pc.pid 
	LEFT JOIN t_post_views pv 
	ON pc.pid = pv.pid 
	LEFT JOIN t_post_category c 
	ON p.category = c.id 
	WHERE pc.display = 1 
	AND p.delete_time IS NULL 
	AND pc.top_time IS NOT NULL 
	LIMIT 1`

	err = db.Get(&topPost, sqlStr)
	return
}

func GetPostList(qs *models.QueryStringGetPostList) (postList models.PostList, err error) {
	sqlStr := `
	SELECT 
    p.pid, 
    p.uid, 
    p.cover, 
    p.title, 
    p.summary, 
    p.category, 
	c.name,
    p.create_time, 
    pc.top_time,
	pv.views 
	FROM t_post p 
	LEFT JOIN t_post_config pc 
	ON p.pid = pc.pid 
	LEFT JOIN t_post_views pv 
	ON pc.pid = pv.pid 
	LEFT JOIN t_post_category c 
	ON p.category = c.id 
	WHERE pc.display = 1 
	AND p.delete_time IS NULL 
	AND pc.top_time IS NULL 
	ORDER BY p.create_time DESC 
	LIMIT ? 
	OFFSET ?`

	var limit = qs.Limit

	if qs.Page == 1 {
		// if the page is 1, get the top post
		var topPost models.PostListItem
		topPost, err = getTopPost()
		if err != nil && err != sql.ErrNoRows {
			return
		}
		if err == nil {
			postList.List = append(postList.List, topPost)
			limit -= 1
		}
	}

	err = db.Select(&postList.List, sqlStr, qs.Limit, (qs.Page-1)*limit)
	if len(postList.List) == 0 {
		postList.List = make([]models.PostListItem, 0, 0)
		return
	}
	if err != nil {
		return
	}

	// get post count
	postCountSql := `
	SELECT COUNT(p.id) 
	FROM t_post p 
	LEFT JOIN t_post_config pc 
	ON p.pid = pc.pid 
	WHERE pc.display = 1 
	AND p.delete_time IS NULL`

	err = db.Get(&postList.Page.Count, postCountSql)
	if err != nil {
		return
	}

	postList.Page.CurrentPage = qs.Page
	postList.Page.Size = qs.Limit

	// get tags
	tagsSqlStr := `SELECT pt.id, pt.name FROM t_post_tag_relation ptr LEFT JOIN t_post_tag pt ON pt.id = ptr.tid WHERE ptr.pid = ?`
	for i := range postList.List {
		err = db.Select(&postList.List[i].Tags, tagsSqlStr, postList.List[i].Pid)
		if err != nil {
			return
		}

		// get commentCount
		postList.List[i].CommentCount, err = GetPostCommentCount(postList.List[i].Pid)
		if err != nil {
			return
		}
	}

	return
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
	c.name,
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

	// get commentCount
	post.CommentCount, err = GetPostCommentCount(post.Pid)
	if err != nil {
		return
	}

	return
}
