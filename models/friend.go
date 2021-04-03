package models

type Friend struct {
	Id     uint64 `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Url    string `db:"url" json:"url"`
	Avatar string `db:"avatar" json:"avatar"`
	Intro  string `db:"intro" json:"intro"`
}

type FriendList struct {
	Page Page     `json:"page"`
	List []Friend `json:"list"`
}
