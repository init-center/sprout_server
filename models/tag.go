package models

type TagData struct {
	Id   uint64 `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Tags = []TagData

type TagList struct {
	Page Page      `json:"page"`
	List []TagData `json:"list"`
}
