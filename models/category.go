package models

type CategoryData struct {
	Id        uint64 `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	PostCount string `db:"post_count" json:"postCount"`
}

type CategoryList struct {
	Page Page           `json:"page"`
	List []CategoryData `json:"list"`
}
