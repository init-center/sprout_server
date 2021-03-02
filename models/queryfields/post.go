package queryfields

type PostQueryFields struct {
	Pid             string `form:"pid"`
	IsDisplay       uint8  `form:"isDisplay,default=2"`
	IsDelete        uint8  `form:"isDelete,default=2"`
	IsTop           uint8  `form:"isTop,default=2"`
	Tag             uint64 `form:"tag"`
	Category        uint64 `form:"category"`
	Keyword         string `form:"keyword"`
	CreateTimeStart string `form:"createTimeStart"`
	CreateTimeEnd   string `form:"createTimeEnd"`
	IsCommentOpen   uint8  `form:"isCommentOpen,default=2"`
	Page            uint64 `form:"page" binding:"gte=1"`
	Limit           uint64 `form:"limit" binding:"gte=1"`
}
