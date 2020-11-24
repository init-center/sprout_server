package models

type TagData struct {
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Tags []TagData
