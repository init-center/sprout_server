package models

import "time"

type FavoritePost struct {
	*PostListItem
	FavoriteTime time.Time `db:"favorite_time" json:"favoriteTime"`
}

type FavoritePostList struct {
	Page Page           `json:"page"`
	List []FavoritePost `json:"list"`
}
