package queryfields

type FavoriteQueryFields struct {
	Pid   uint64 `form:"pid"`
	Uid   string `form:"uid"`
	Page  uint64 `form:"page" binding:"omitempty,gte=1"`
	Limit uint64 `form:"limit" binding:"omitempty,gte=1"`
}
