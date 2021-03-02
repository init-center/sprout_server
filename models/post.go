package models

import "time"

type PostListItem struct {
	Uid          string     `json:"uid" db:"uid"`
	Pid          uint64     `json:"pid,string" db:"pid"`
	CategoryId   uint64     `json:"categoryId" db:"category"`
	CategoryName string     `json:"categoryName" db:"name"`
	Tags         Tags       `json:"tags" db:"tags"`
	Title        string     `json:"title" db:"title"`
	Cover        string     `json:"cover" db:"cover"`
	Summary      string     `json:"summary" db:"summary"`
	Views        uint64     `json:"views" db:"views"`
	Favorites    uint64     `json:"favorites" db:"favorites"`
	CreateTime   time.Time  `json:"createTime" db:"create_time"`
	TopTime      *time.Time `json:"topTime" db:"top_time"`
	CommentCount uint64     `json:"commentCount" db:"comment_count"`
}

type PostDetail struct {
	*PostListItem
	Bgm         string    `json:"bgm" db:"bgm"`
	Content     string    `json:"content" db:"content"`
	CommentOpen uint8     `json:"isCommentOpen" db:"comment_open"`
	UpdateTime  time.Time `json:"updateTime" db:"update_time"`
}

type PostItemByAdmin struct {
	*PostDetail
	IsDisplay  uint8      `json:"isDisplay" db:"display"`
	DeleteTime *time.Time `json:"deleteTime" db:"delete_time"`
}

type PostListByAdmin struct {
	Page PostPage          `json:"page"`
	List []PostItemByAdmin `json:"list"`
}

type PostPage struct {
	Count       uint64 `json:"count"`
	CurrentPage uint64 `json:"currentPage"`
	Size        uint64 `json:"size"`
}

type PostList struct {
	Page PostPage       `json:"page"`
	List []PostListItem `json:"list"`
}
