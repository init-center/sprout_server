package models

import "time"

type BaseCommentItem struct {
	Cid        uint64     `db:"cid" json:"cid,string"`
	Pid        uint64     `db:"pid" json:"pid,string"`
	Uid        string     `db:"uid" json:"uid"`
	UserName   string     `db:"user_name" json:"userName"`
	Avatar     string     `db:"avatar" json:"avatar"`
	Content    string     `db:"content" json:"content"`
	OS         string     `db:"os" json:"os"`
	Browser    string     `db:"browser" json:"browser"`
	Engine     string     `db:"engine" json:"engine"`
	CreateTime time.Time  `db:"create_time" json:"createTime"`
	UpdateTime time.Time  `db:"update_time" json:"updateTime"`
	DeleteTime *time.Time `db:"delete_time" json:"deleteTime"`
}

type CommentItem struct {
	*BaseCommentItem
	TargetCid  *uint64 `db:"target_cid" json:"targetCid,string"`
	TargetUid  *string `db:"target_uid" json:"targetUid"`
	TargetName *string `db:"target_name" json:"targetName"`
	ParentCid  *uint64 `db:"parent_cid" json:"parentCid,string"`
	ParentUid  *string `db:"parent_uid" json:"parentUid"`
}

type CommentItemByAdmin struct {
	*CommentItem
	Ip           string `db:"ip" json:"ip"`
	PostTitle    string `db:"post_title" json:"postTitle"`
	ReviewStatus uint8  `db:"review_status" json:"reviewStatus"`
}

type Comment struct {
	*BaseCommentItem
	Page     Page          `json:"page"`
	Children []CommentItem `json:"children"`
}

type CommentList struct {
	Page Page      `json:"page"`
	List []Comment `json:"list"`
}

type CommentItemListByAdmin struct {
	List []CommentItemByAdmin `json:"list"`
	Page Page                 `json:"page"`
}

type ParentCommentChildren struct {
	Page Page          `json:"page"`
	List []CommentItem `json:"list"`
}
