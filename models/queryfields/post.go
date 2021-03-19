package queryfields

type PostQueryFields struct {
	Pid             string `form:"pid" json:"pid"`
	IsDisplay       uint8  `form:"isDisplay,default=2" json:"isDisplay"`
	IsDelete        uint8  `form:"isDelete,default=2" json:"isDelete"`
	IsTop           uint8  `form:"isTop,default=2" json:"isTop"`
	Tag             uint64 `form:"tag" json:"tag"`
	TagName         string `form:"tagName" json:"tagName"`
	CategoryName    string `form:"categoryName" json:"categoryName"`
	Category        uint64 `form:"category" json:"category"`
	Keyword         string `form:"keyword" json:"keyword"`
	CreateTimeStart string `form:"createTimeStart" json:"createTimeStart"`
	CreateTimeEnd   string `form:"createTimeEnd" json:"createTimeEnd"`
	IsCommentOpen   uint8  `form:"isCommentOpen,default=2" json:"isCommentOpen"`
	Page            uint64 `form:"page" binding:"omitempty,gte=1" json:"page"`
	Limit           uint64 `form:"limit" binding:"omitempty,gte=1" json:"limit"`
}
