package mysql

import (
	"sprout_server/common/snowflake"
	"sprout_server/models"
)

func getUidByCid(cid int64) (uid string, err error) {
	sqlStr := `SELECT uid FROM t_post_comment WHERE cid = ?`
	err = db.Get(&uid, sqlStr, cid)
	return
}

func getParentCidByTargetCid(targetCid int64) (parentCid int64, err error) {
	sqlStr := `SELECT parent_cid FROM t_post_comment WHERE cid = ?`
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

func CheckPostCommentExist(cid int64) (bool, error) {
	sqlStr := `SELECT count(id) FROM t_post_comment WHERE cid = ?`
	var count int
	if err := db.Get(&count, sqlStr, cid); err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetPostCommentCount(pid int64) (int64, error) {
	sqlStr := `SELECT count(id) FROM t_post_comment WHERE pid = ?`
	var count int64
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

	childTargetNameSql := `SELECT name AS target_name FROM t_user u WHERE u.uid = ?`
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
		if err != nil {
			return
		}
	}
	return
}
