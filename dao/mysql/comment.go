package mysql

import (
	"database/sql"
	"sprout_server/common/snowflake"
	"sprout_server/models"
	"sprout_server/models/queryfields"
)

func getUidByCid(cid uint64) (uid string, err error) {
	sqlStr := `SELECT uid FROM t_post_comment WHERE cid = ?`
	err = db.Get(&uid, sqlStr, cid)
	return
}

func getParentCidByTargetCid(targetCid uint64) (parentCid uint64, err error) {
	sqlStr := `SELECT IFNULL(parent_cid, 0) FROM t_post_comment WHERE cid = ?`
	err = db.Get(&parentCid, sqlStr, targetCid)
	return
}

func CreatePostComment(p *models.ParamsAddComment) (err error) {
	var cid = snowflake.GenID()

	hasTarget := p.TargetCid != 0
	var sqlStr string
	if hasTarget {
		targetUid, err := getUidByCid(p.TargetCid)
		if err != nil {
			return err
		}

		// get parentCid by targetCid
		parentCid, err := getParentCidByTargetCid(p.TargetCid)
		if err != nil {
			return err
		}

		// If parentCid is zero, then targetCid becomes parentCid
		if parentCid == 0 {
			parentCid = p.TargetCid
		}

		parentUid, err := getUidByCid(parentCid)
		if err != nil {
			return err
		}

		sqlStr = `INSERT INTO t_post_comment(cid, pid, uid, target_cid, target_uid, parent_cid, parent_uid, content) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = db.Exec(sqlStr, cid, p.Pid, p.Uid, p.TargetCid, targetUid, parentCid, parentUid, p.Content)
	} else {
		sqlStr = `INSERT INTO t_post_comment(cid, pid, uid, content) VALUES(?, ?, ?, ?)`
		_, err = db.Exec(sqlStr, cid, p.Pid, p.Uid, p.Content)
	}

	return
}

func CheckPostCommentExist(cid uint64) (bool, error) {
	sqlStr := `SELECT count(id) FROM t_post_comment WHERE cid = ?`
	var count int
	if err := db.Get(&count, sqlStr, cid); err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetPostCommentCount(pid uint64) (uint64, error) {
	sqlStr := `SELECT count(id) FROM t_post_comment WHERE pid = ? AND review_status = 1 AND delete_time IS NULL`
	var count uint64
	if err := db.Get(&count, sqlStr, pid); err != nil {
		return 0, err
	}

	return count, nil
}

func GetPostCommentList(p *models.ParamsGetCommentList) (commentList models.CommentList, err error) {
	sqlStr := `
	SELECT 
	pc.cid,
    pc.pid, 
    pc.uid,
	pc.content,
	u.name AS user_name,
	u.avatar,
    pc.create_time, 
	pc.update_time 
	FROM t_post_comment pc 
	LEFT JOIN t_user u 
	ON pc.uid = u.uid 
	WHERE pc.pid = ? 
	AND pc.parent_cid IS NULL 
	AND pc.review_status = 1 
	AND pc.delete_time IS NULL 
	ORDER BY pc.create_time DESC 
	LIMIT ? 
	OFFSET ?`

	err = db.Select(&commentList.List, sqlStr, p.Pid, p.Limit, (p.Page-1)*p.Limit)
	if err != nil {
		return
	}

	// get post parentComment count
	parentCountSqlStr := `
	SELECT 
	COUNT(id) 
	FROM t_post_comment pc 
	WHERE pc.pid = ? 
	AND pc.parent_cid IS NULL 
	AND pc.review_status = 1 
	AND pc.delete_time IS NULL`

	err = db.Get(&commentList.Page.Count, parentCountSqlStr, p.Pid)
	if err != nil {
		return
	}

	commentList.Page.CurrentPage = p.Page
	commentList.Page.Size = p.Limit

	if len(commentList.List) == 0 {
		commentList.List = make([]models.Comment, 0, 0)
		return
	}

	// get children
	childSqlStr := `
	SELECT 
	pc.cid,
    pc.pid, 
    pc.uid,
	pc.content,
	pc.target_cid,
	pc.target_uid,
	pc.parent_cid,
	pc.parent_uid,
	u.name AS user_name,
	u.avatar,
    pc.create_time, 
	pc.update_time 
	FROM t_post_comment pc 
	LEFT JOIN t_user u 
	ON pc.uid = u.uid 
	WHERE pc.pid = ? 
	AND pc.parent_cid = ? 
	AND pc.review_status = 1 
	AND pc.delete_time IS NULL 
	ORDER BY pc.create_time DESC 
	LIMIT ? 
	OFFSET ?`

	// get child targetName sql
	childTargetNameSql := `SELECT name AS target_name FROM t_user u WHERE u.uid = ?`

	// get comment childComment count
	childCommentCountSqlStr := `
	SELECT 
	COUNT(id) 
	FROM t_post_comment pc 
	WHERE pc.pid = ? 
	AND pc.parent_cid = ?  
	AND pc.review_status = 1 
	AND pc.delete_time IS NULL`

	for i := range commentList.List {
		err = db.Select(&commentList.List[i].Children, childSqlStr, p.Pid, commentList.List[i].Cid, p.ChildLimit, (p.Page-1)*p.Limit)
		if err != nil {
			return
		}

		if len(commentList.List[i].Children) == 0 {
			commentList.List[i].Children = make([]models.CommentItem, 0, 0)
		}

		err = db.Get(&commentList.List[i].Page.Count, childCommentCountSqlStr, p.Pid, commentList.List[i].Cid)
		if err != nil {
			return
		}

		commentList.List[i].Page.CurrentPage = p.Page
		commentList.List[i].Page.Size = p.ChildLimit

		for j := range commentList.List[i].Children {
			err = db.Get(&commentList.List[i].Children[j].TargetName, childTargetNameSql, commentList.List[i].Children[j].TargetUid)
			if err != nil {
				return
			}
		}
	}

	return
}

func GetPostParentCommentChildren(p *models.ParamsGetParentCommentChildren) (parentCommentChildren models.ParentCommentChildren, err error) {
	childSqlStr := `
	SELECT 
	pc.cid,
    pc.pid, 
    pc.uid,
	pc.content,
	pc.target_cid,
	pc.target_uid,
	pc.parent_cid,
	pc.parent_uid,
	u.name AS user_name,
	u.avatar,
    pc.create_time, 
	pc.update_time 
	FROM t_post_comment pc 
	LEFT JOIN t_user u 
	ON pc.uid = u.uid 
	WHERE pc.pid = ? 
	AND pc.parent_cid = ? 
	AND pc.review_status = 1 
	AND pc.delete_time IS NULL 
	ORDER BY pc.create_time DESC 
	LIMIT ? 
	OFFSET ?`

	childTargetNameSql := `SELECT name AS target_name FROM t_user WHERE uid = ?`
	childCommentCountSqlStr := `
	SELECT 
	COUNT(id) 
	FROM t_post_comment pc 
	WHERE pc.pid = ? 
	AND pc.parent_cid = ?  
	AND pc.review_status = 1 
	AND pc.delete_time IS NULL`

	err = db.Select(&parentCommentChildren.List, childSqlStr, p.Pid, p.Cid, p.Limit, (p.Page-1)*p.Limit)
	if err != nil {
		return
	}

	if len(parentCommentChildren.List) == 0 {
		parentCommentChildren.List = make([]models.CommentItem, 0, 0)
	}

	err = db.Get(&parentCommentChildren.Page.Count, childCommentCountSqlStr, p.Pid, p.Cid)
	if err != nil {
		return
	}

	parentCommentChildren.Page.CurrentPage = p.Page
	parentCommentChildren.Page.Size = p.Limit

	for i := range parentCommentChildren.List {
		err = db.Get(&parentCommentChildren.List[i].TargetName, childTargetNameSql, parentCommentChildren.List[i].TargetUid)
		if err != nil && err != sql.ErrNoRows {
			return
		}
	}
	return
}

func GetPostComments(queryFields *queryfields.CommentQueryFields) (comments models.CommentItemListByAdmin, err error) {
	sqlStr := `
	SELECT 
	pc.cid,
    pc.pid, 
	p.title AS post_title, 
    pc.uid,
	pc.content,
	pc.target_cid,
	pc.target_uid,
	pc.parent_cid,
	pc.parent_uid,
	u.name AS user_name,
	u.avatar,
	pc.review_status,
    pc.create_time, 
	pc.update_time,
	pc.delete_time 
	FROM t_post_comment pc 
	LEFT JOIN t_user u 
	ON pc.uid = u.uid 
	LEFT JOIN t_post p 
	ON pc.pid = p.pid 
	WHERE `

	sqlStr = dynamicConcatCommentSql(sqlStr, queryFields)
	sqlStr += `ORDER BY pc.create_time DESC `

	var limit = queryFields.Limit

	if queryFields.Limit != 0 && queryFields.Page != 0 {
		sqlStr += `LIMIT ? OFFSET ?`
		err = db.Select(&comments.List, sqlStr, queryFields.Uid, queryFields.Pid,
			queryFields.CreateTimeStart, queryFields.CreateTimeEnd,
			limit, (queryFields.Page-1)*limit)
	} else {
		err = db.Select(&comments.List, sqlStr, queryFields.Uid, queryFields.Pid,
			queryFields.CreateTimeStart, queryFields.CreateTimeEnd)
	}
	if err != nil {
		return
	}

	if len(comments.List) == 0 {
		comments.List = make([]models.CommentItemByAdmin, 0, 0)
		return
	}

	targetNameSql := `SELECT name AS target_name FROM t_user WHERE uid = ?`
	for i := range comments.List {
		err = db.Get(&comments.List[i].TargetName, targetNameSql, comments.List[i].TargetUid)
		if err != nil && err != sql.ErrNoRows {
			return
		}
	}

	countSqlStr := `
	SELECT COUNT(DISTINCT pc.id) 
	FROM t_post_comment pc 
	LEFT JOIN t_user u 
	ON pc.uid = u.uid 
	LEFT JOIN t_post p 
	ON pc.pid = p.pid 
	WHERE `
	countSqlStr = dynamicConcatCommentSql(countSqlStr, queryFields)
	err = db.Get(&comments.Page.Count, countSqlStr, queryFields.Uid, queryFields.Pid,
		queryFields.CreateTimeStart, queryFields.CreateTimeEnd)
	if err != nil {
		return
	}

	comments.Page.CurrentPage = queryFields.Page
	comments.Page.Size = queryFields.Limit

	return

}

func AdminUpdatePostComment(p *models.ParamsAdminUpdateComment, u *models.UriUpdateComment) (err error) {
	sqlStr := `
	UPDATE t_post_comment SET 
	delete_time = CASE ? WHEN NULL THEN delete_time 
	WHEN 0 THEN NULL 
	WHEN 1 THEN NOW() 
	ELSE delete_time END,
	review_status = CASE ? WHEN NULL THEN review_status 
	WHEN 0 THEN 0 
	WHEN 1 THEN 1 
	WHEN 2 THEN 2 
	ELSE review_status END,
	content = IFNULL(?, content) 
	WHERE cid = ?`

	_, err = db.Exec(sqlStr, p.IsDelete, p.ReviewState, p.Content, u.Cid)
	if err != nil {
		return
	}
	return
}

func dynamicConcatCommentSql(sqlStr string, queryFields *queryfields.CommentQueryFields) string {
	if queryFields.Uid == "" {
		sqlStr += ` LENGTH(?) = 0 `
	} else {
		sqlStr += ` u.uid = ? `
	}

	if queryFields.Pid == "" {
		sqlStr += ` AND LENGTH(?) = 0 `
	} else {
		sqlStr += ` AND pc.pid = ? `
	}

	if queryFields.IsDelete == 0 {
		sqlStr += ` AND pc.delete_time IS NULL `
	} else if queryFields.IsDelete == 1 {
		sqlStr += ` AND pc.delete_time IS NOT NULL `
	}

	switch queryFields.ReviewStatus {
	case 0:
		sqlStr += ` AND pc.review_status = 0 `
	case 1:
		sqlStr += ` AND pc.review_status = 1 `
	case 2:
		sqlStr += ` AND pc.review_status = 2 `
	default:

	}

	if queryFields.CreateTimeStart != "" && queryFields.CreateTimeEnd != "" {
		sqlStr += ` AND (p.create_time >= ? AND p.create_time <= ?) `
	} else {
		sqlStr += ` AND LENGTH(?)= 0 AND LENGTH(?) = 0 `
	}

	return sqlStr
}
