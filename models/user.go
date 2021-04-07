package models

import "time"

type UserBasicInfo struct {
	Uid    string `db:"uid" json:"uid"`
	Name   string `db:"name" json:"name"`
	Avatar string `db:"avatar" json:"avatar"`
}

type UserPublicInfo struct {
	*UserBasicInfo
	Intro      string     `db:"intro" json:"intro"`
	CreateTime *time.Time `db:"create_time" json:"createTime"`
	IsBaned    uint8      `db:"is_baned" json:"isBaned"`
	Group      uint8      `db:"group" json:"group"`
}

type User struct {
	*UserBasicInfo
	PassWord string `db:"password" json:"password"`
	Email    string `db:"email" json:"email"`
}

type UserDetailByAdmin struct {
	*UserPublicInfo
	Email        string     `db:"email" json:"email"`
	Gender       *uint8     `db:"gender" json:"gender"`
	Tel          *string    `db:"tel" json:"tel"`
	Birthday     *time.Time `db:"birthday" json:"birthday"`
	UpdateTime   *time.Time `db:"update_time" json:"updateTime"`
	DeleteTime   *time.Time `db:"delete_time" json:"deleteTime"`
	BanStartTime *time.Time `db:"ban_start_time" json:"banStartTime"`
	BanEndTime   *time.Time `db:"ban_end_time" json:"banEndTime"`
}

type UserPrivateInfo = UserDetailByAdmin

type UserDetailList struct {
	Page Page                `json:"page"`
	List []UserDetailByAdmin `json:"list"`
}

type BanTime struct {
	BanStartTime *time.Time `db:"ban_start_time" json:"banStartTime"`
	BanEndTime   *time.Time `db:"ban_end_time" json:"banEndTime"`
}
