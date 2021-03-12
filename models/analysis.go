package models

type RecentIncrease struct {
	Date     string `json:"date" db:"date"`
	Increase uint64 `json:"increase" db:"increase"`
}

type BaseAnalysisData struct {
	Total              uint64           `json:"total" db:"total"`
	RecentIncreaseList []RecentIncrease `json:"recentIncreaseList"`
	TodayIncrease      uint64           `json:"todayIncrease" db:"today_increase"`
}

type PostAnalysisData struct {
	Total         uint64 `json:"total" db:"total"`
	Average       uint64 `json:"average" db:"average"`
	MonthIncrease uint64 `json:"monthIncrease" db:"month_increase"`
}

type PostViewsItem struct {
	Id    uint64 `json:"pid" db:"pid"`
	Title string `json:"title" db:"title"`
	Views uint64 `json:"views" db:"views"`
}

type PostViewsRank = []PostViewsItem

type CategoriesPostsCount struct {
	Name  string `json:"name" db:"name"`
	Value uint64 `json:"value" db:"value"`
}

type TagsPostsCount = CategoriesPostsCount
