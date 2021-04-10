package models

import (
	"sprout_server/models/queryfields"
	"time"
)

type PostListItem struct {
	Uid          string     `json:"uid" db:"uid"`
	UserName     string     `json:"userName" db:"user_name"`
	Pid          uint64     `json:"pid,string" db:"pid"`
	CategoryId   uint64     `json:"categoryId" db:"category"`
	CategoryName string     `json:"categoryName" db:"category_name"`
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
	Page   Page                        `json:"page"`
	List   []PostItemByAdmin           `json:"list"`
	Search queryfields.PostQueryFields `json:"search"`
}

type PostList struct {
	Page   Page                   `json:"page"`
	List   []PostListItem         `json:"list"`
	Search QueryStringGetPostList `json:"search"`
}

type PostDetailList struct {
	Page   Page                   `json:"page"`
	List   []PostDetail           `json:"list"`
	Search QueryStringGetPostList `json:"search"`
}
