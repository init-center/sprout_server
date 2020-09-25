package models

type User struct {
	Uid      string `db:"uid"`
	PassWord string `db:"password"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Avatar   string `db:"avatar"`
}
