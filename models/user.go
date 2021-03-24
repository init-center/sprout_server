package models

import "time"

type UserPublicInfo struct {
	Uid    string `db:"uid" json:"uid"`
	Name   string `db:"name" json:"name"`
	Avatar string `db:"avatar" json:"avatar"`
	Intro  string `db:"intro" json:"intro"`
}

type User struct {
	*UserPublicInfo
	PassWord string `db:"password" json:"password"`
	Email    string `db:"email" json:"email"`
}

type UserDetailByAdmin struct {
	*UserPublicInfo
	Email        string     `db:"email" json:"email"`
	Gender       *uint8     `db:"gender" json:"gender"`
	Tel          *string    `db:"tel" json:"tel"`
	Birthday     *time.Time `db:"birthday" json:"birthday"`
	Group        uint8      `db:"group" json:"group"`
	CreateTime   *time.Time `db:"create_time" json:"createTime"`
	UpdateTime   *time.Time `db:"update_time" json:"updateTime"`
	DeleteTime   *time.Time `db:"delete_time" json:"deleteTime"`
	BanStartTime *time.Time `db:"ban_start_time" json:"banStartTime"`
	BanEndTime   *time.Time `db:"ban_end_time" json:"banEndTime"`
	IsBaned      uint8      `db:"is_baned" json:"isBaned"`
}

type UserDetailList struct {
	Page Page                `json:"page"`
	List []UserDetailByAdmin `json:"list"`
}
