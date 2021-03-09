package queryfields

type CommentQueryFields struct {
	Pid             string `form:"pid"`
	Uid             string `form:"uid"`
	ReviewStatus    uint8  `form:"reviewStatus,default=3"`
	CreateTimeStart string `form:"createTimeStart"`
	CreateTimeEnd   string `form:"createTimeEnd"`
	IsDelete        uint8  `form:"isDelete,default=2"`
	Page            uint64 `form:"page" binding:"omitempty,gte=1"`
	Limit           uint64 `form:"limit" binding:"omitempty,gte=1"`
}
