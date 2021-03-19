package models

type TagData struct {
	Id        uint64 `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	PostCount string `db:"post_count" json:"postCount"`
}

type Tags = []TagData

type TagList struct {
	Page Page      `json:"page"`
	List []TagData `json:"list"`
}
