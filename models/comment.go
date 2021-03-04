package models

import "time"

type BaseCommentItem struct {
	Cid        uint64     `db:"cid" json:"cid,string"`
	Pid        uint64     `db:"pid" json:"pid,string"`
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
	TargetCid  *uint64 `db:"target_cid" json:"targetCid,string"`
	TargetUid  *string `db:"target_uid" json:"targetUid"`
	TargetName *string `db:"target_name" json:"targetName"`
	ParentCid  *uint64 `db:"parent_cid" json:"parentCid,string"`
	ParentUid  *string `db:"parent_uid" json:"parentUid"`
}

type CommentItemByAdmin struct {
	*CommentItem
	PostTitle    string `db:"post_title" json:"postTitle"`
	ReviewStatus uint8  `db:"review_status" json:"reviewStatus"`
}

type CommentPage struct {
	Count       uint64 `json:"count"`
	CurrentPage uint64 `json:"currentPage"`
	Size        uint64 `json:"size"`
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

type CommentItemListByAdmin struct {
	List []CommentItemByAdmin `json:"list"`
	Page CommentPage          `json:"page"`
}

type ParentCommentChildren struct {
	Page CommentPage   `json:"page"`
	List []CommentItem `json:"list"`
}
