package models

type CategoryData struct {
	Id   uint64 `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Categories []CategoryData
