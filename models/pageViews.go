package models

type PageViews struct {
	Visits       uint64 `json:"visits" db:"visits"`
	Distinctions uint64 `json:"distinctions" db:"distinctions"`
}
