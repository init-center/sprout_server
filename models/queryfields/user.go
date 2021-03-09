package queryfields

type UserQueryFields struct {
	Uid             string `form:"uid"`
	Name            string `form:"name"`
	Email           string `form:"email"`
	IsBaned         uint8  `form:"isBaned,default=2"`
	IsDelete        uint8  `form:"isDelete,default=2"`
	Group           uint8  `form:"group,default=3"`
	CreateTimeStart string `form:"createTimeStart"`
	CreateTimeEnd   string `form:"createTimeEnd"`
	BanTimeStart    string `form:"banTimeStart"`
	BanTimeEnd      string `form:"banTimeEnd"`
	Gender          uint8  `form:"gender,default=2"`
	Tel             string `form:"tel"`
	Page            uint64 `form:"page" binding:"omitempty,gte=1"`
	Limit           uint64 `form:"limit" binding:"omitempty,gte=1"`
}
