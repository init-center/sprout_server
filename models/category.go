package models

type CategoryData struct {
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Categories []CategoryData
