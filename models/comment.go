package models

import "time"

type BaseCommentItem struct {
	Cid        int64      `db:"cid" json:"cid,string"`
	Pid        int64      `db:"pid" json:"pid"`
	Uid        string     `db:"uid" json:"uid"`
	UserName   string     `db:"user_name" json:"userName"`
	Avatar     string     `db:"avatar" json:"avatar"`
	Content    string     `db:"content" json:"content"`
	CreateTime time.Time  `db:"create_time" json:"createTime"`
	UpdateTime time.Time  `db:"update_time" json:"updateTime"`
	DeleteTime *time.Time `db:"delete_time" json:"deleteTime"`
}

type CommentItem struct {
	*BaseCommentItem
	TargetCid  int64  `db:"target_cid" json:"targetCid,string"`
	TargetUid  string `db:"target_uid" json:"targetUid"`
	TargetName string `db:"target_name" json:"targetName"`
	ParentCid  int64  `db:"parent_cid" json:"parentCid,string"`
	ParentUid  string `db:"parent_uid" json:"parentUid"`
}

type CommentPage struct {
	Count       int64 `json:"count"`
	CurrentPage int64 `json:"currentPage"`
	Size        int64 `json:"size"`
}

type Comment struct {
	*BaseCommentItem
	Page     CommentPage   `json:"page"`
	Children []CommentItem `json:"children"`
}

type CommentList struct {
	Page CommentPage `json:"page"`
	List []Comment   `json:"list"`
}

type ParentCommentChildren struct {
	Page CommentPage   `json:"page"`
	List []CommentItem `json:"list"`
}
