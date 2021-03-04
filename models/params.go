package models

type ParamsSignUp struct {
	Uid        string `json:"uid" binding:"required,min=2,max=12,ne=admin,checkUid"`
	Name       string `json:"name" binding:"required,ne=admin,checkName"`
	Password   string `json:"password" binding:"required,min=6,max=16,checkPwd"`
	RePassword string `json:"rePassword" binding:"required,eqfield=Password,min=6,max=16,checkPwd"`
	Email      string `json:"email" binding:"required,email"`
	ECode      string `json:"eCode" binding:"required"` // email verify code
}

type ParamsGetECode struct {
	Email string `json:"email" binding:"required,email"`
	Type  int    `json:"type" binding:"required"`
}

type ParamsSignIn struct {
	Uid      string `json:"uid" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=16,checkPwd"`
}

type ParamsAddCategory struct {
	Name string `json:"name" binding:"required,min=1,max=64"`
}

type ParamsAddTag struct {
	Name string `json:"name" binding:"required,min=1,max=64"`
}

type ParamsAddPost struct {
	Category      uint64   `json:"category" binding:"required"`
	Tags          []uint64 `json:"tags" binding:"required,min=1,max=20"`
	Title         string   `json:"title" binding:"required,min=1,max=128"`
	Cover         string   `json:"cover" binding:"required,url"`
	Bgm           string   `json:"bgm" binding:"required,url"`
	Summary       string   `json:"summary" binding:"required,gte=2,max=128"`
	Content       string   `json:"content" binding:"required,gte=2"`
	IsDisplay     uint8    `json:"isDisplay" binding:"oneof=0 1"`
	IsCommentOpen uint8    `json:"isCommentOpen" binding:"oneof=0 1"`
	IsTop         uint8    `json:"isTop" binding:"oneof=0 1"`
	IsDelete      *uint8   `json:"isDelete" binding:"oneof=0 1"`
}

type QueryStringGetPostList struct {
	Page  uint64 `form:"page" binding:"gte=1"`
	Limit uint64 `form:"limit" binding:"gte=1"`
}

type UriGetPostDetail struct {
	Pid uint64 `uri:"pid" binding:"required"`
}

type ParamsAddComment struct {
	Pid       uint64 `uri:"pid"`
	Uid       string `json:"uid"`
	TargetCid uint64 `json:"targetCid,string"`
	Content   string `json:"content" binding:"required,min=1,max=1024"`
}

type ParamsGetCommentList struct {
	Pid        uint64 `uri:"pid"`
	Page       uint64 `form:"page" binding:"gte=1"`
	Limit      uint64 `form:"limit" binding:"gte=1"`
	ChildLimit uint64 `form:"child_limit,default=2" binding:"gte=1"`
}

type ParamsGetParentCommentChildren struct {
	Pid   uint64 `uri:"pid"`
	Cid   uint64 `uri:"cid"`
	Page  uint64 `form:"page" binding:"gte=1"`
	Limit uint64 `form:"limit" binding:"gte=1"`
}

type ParamsPostFavorite struct {
	Pid uint64 `uri:"pid"`
	Uid string
}

type UriUpdatePost struct {
	Pid uint64 `uri:"pid" binding:"required"`
}

type ParamsUpdatePost struct {
	Category      *uint64  `json:"category" binding:"omitempty"`
	Tags          []uint64 `json:"tags" binding:"omitempty,min=1,max=20"`
	Title         *string  `json:"title" binding:"omitempty,min=1,max=128"`
	Cover         *string  `json:"cover" binding:"omitempty,url"`
	Bgm           *string  `json:"bgm" binding:"omitempty,url"`
	Summary       *string  `json:"summary" binding:"omitempty,gte=2,max=128"`
	Content       *string  `json:"content" binding:"omitempty,gte=2"`
	IsDisplay     *uint8   `json:"isDisplay" binding:"omitempty,oneof=0 1"`
	IsCommentOpen *uint8   `json:"isCommentOpen" binding:"omitempty,oneof=0 1"`
	IsTop         *uint8   `json:"isTop" binding:"omitempty,oneof=0 1"`
	IsDelete      *uint8   `json:"isDelete" binding:"omitempty,oneof=0 1"`
}

type UriUpdateComment struct {
	Cid uint64 `uri:"cid" binding:"required"`
}

type ParamsAdminUpdateComment struct {
	IsDelete    *uint8  `json:"isDelete" binding:"omitempty,oneof=0 1"`
	ReviewState *uint8  `json:"reviewStatus" binding:"omitempty,oneof=0 1 2"`
	Content     *string `json:"content" binding:"omitempty,gte=2"`
}
