package models

type UserPublicInfo struct {
	Uid    string `db:"uid"`
	Name   string `db:"name"`
	Avatar string `db:"avatar"`
}

type User struct {
	*UserPublicInfo
	PassWord string `db:"password"`
	Email    string `db:"email"`
}
