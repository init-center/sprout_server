package models

type Page struct {
	Count       uint64 `json:"count"`
	CurrentPage uint64 `json:"currentPage"`
	Size        uint64 `json:"size"`
}
