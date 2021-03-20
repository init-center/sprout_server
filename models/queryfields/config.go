package queryfields

type ConfigQueryFields struct {
	Key     string `form:"key"`
	Value   string `form:"value"`
	Explain string `form:"explain"`
	Page    uint64 `form:"page" binding:"omitempty,gte=1"`
	Limit   uint64 `form:"limit" binding:"omitempty,gte=1"`
}
