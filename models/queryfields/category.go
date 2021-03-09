package queryfields

type CategoryQueryFields struct {
	Id      uint64 `form:"id"`
	Keyword string `form:"keyword"`
	Page    uint64 `form:"page" binding:"omitempty,gte=1"`
	Limit   uint64 `form:"limit" binding:"omitempty,gte=1"`
}
