package models

import "time"

type Tag struct {
	Id   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type PostListItem struct {
	Uid          string     `json:"uid" db:"uid"`
	Pid          int64      `json:"pid,string" db:"pid"`
	CategoryId   int64      `json:"categoryId" db:"category"`
	CategoryName string     `json:"categoryName" db:"name"`
	Tags         []Tag      `json:"tags" db:"tags"`
	Title        string     `json:"title" db:"title"`
	Cover        string     `json:"cover" db:"cover"`
	Summary      string     `json:"summary" db:"summary"`
	Views        int64      `json:"views" db:"views"`
	CreateTime   time.Time  `json:"createTime" db:"create_time"`
	TopTime      *time.Time `json:"topTime" db:"top_time"`
	CommentCount int64      `json:"commentCount" db:"comment_count"`
}

type PostDetail struct {
	*PostListItem
	Bgm         string    `json:"bgm" db:"bgm"`
	Content     string    `json:"content" db:"content"`
	CommentOpen int8      `json:"commentOpen" db:"comment_open"`
	UpdateTime  time.Time `json:"updateTime" db:"update_time"`
}

type PostDetailByAdmin struct {
	*PostDetail
	Display    int8       `json:"display" db:"display"`
	DeleteTime *time.Time `json:"deleteTime" db:"delete_time"`
}

type PostPage struct {
	Count       int64 `json:"count"`
	CurrentPage int64 `json:"currentPage"`
	Size        int64 `json:"size"`
}

type PostList struct {
	Page PostPage       `json:"page"`
	List []PostListItem `json:"list"`
}
