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

type UriUpdateCategory struct {
	Id uint64 `uri:"id" binding:"required"`
}

type UriDeleteCategory = UriUpdateCategory

type ParamsAddTag struct {
	Name string `json:"name" binding:"required,min=1,max=64"`
}

type UriUpdateTag struct {
	Id uint64 `uri:"id" binding:"required"`
}

type UriDeleteTag = UriUpdateTag

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
	Tag             uint64 `form:"tag" json:"tag"`
	Category        uint64 `form:"category" json:"category"`
	TagName         string `form:"tagName" json:"tagName"`
	CategoryName    string `form:"categoryName" json:"categoryName"`
	Keyword         string `form:"keyword" json:"keyword"`
	FirstPageGetTop uint8  `form:"firstPageGetTop" binding:"oneof=0 1" json:"firstPageGetTop"`
	Page            uint64 `form:"page" binding:"omitempty,gte=1" json:"page"`
	Limit           uint64 `form:"limit" binding:"omitempty,gte=1" json:"limit"`
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
	Cid        uint64 `form:"cid"`
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

type UriUpdateUser struct {
	Uid string `uri:"uid" binding:"required"`
}

type UriGetUserInfo = UriUpdateUser
type UriGetBanTime = UriUpdateUser

type BaseParamsUpdateUser struct {
	Name     *string `json:"name" binding:"omitempty,ne=admin,checkName"`
	Avatar   *string `json:"avatar" binding:"omitempty,url"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Gender   *uint8  `json:"gender" binding:"omitempty,oneof=0 1"`
	Tel      *string `json:"tel" binding:"omitempty,checkTel"`
	Birthday *string `json:"birthday" binding:"omitempty,datetime=2006-01-02"`
	Intro    *string `json:"intro" binding:"omitempty,max=128"`
}

type ParamsUpdateUser struct {
	*BaseParamsUpdateUser
	ECode      string  `json:"eCode" binding:"required"`
	Password   *string `json:"password" binding:"omitempty,min=6,max=16,checkPwd"`
	RePassword *string `json:"rePassword" binding:"omitempty,eqfield=Password,min=6,max=16,checkPwd"`
}

type ParamsUpdatePassword struct {
	ECode      string `json:"eCode" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6,max=16,checkPwd"`
	RePassword string `json:"rePassword" binding:"required,eqfield=Password,min=6,max=16,checkPwd"`
}

type ParamsAdminUpdateUser struct {
	*BaseParamsUpdateUser
	Password *string `json:"password" binding:"omitempty,min=6,max=16,checkPwd"`
	Group    *uint8  `json:"group" binding:"omitempty,oneof=1 2"`
	IsDelete *uint8  `json:"isDelete" binding:"omitempty,oneof=0 1"`
}

type ParamsBanUser struct {
	BanStartTime string `json:"banStartTime" binding:"datetime=2006-01-02 15:04:05"`
	BanEndTime   string `json:"banEndTime" binding:"datetime=2006-01-02 15:04:05"`
}

type ParamsRecentDaysAnalysis struct {
	Days uint8 `json:"days" binding:"max=30"`
}

type ParamsRecentMonthsAnalysis struct {
	Months uint8 `json:"months" binding:"max=12"`
}

type ParamsGetViewsRank struct {
	Limit uint8 `form:"limit,default=7"`
}

type ParamsAddConfig struct {
	Key     string `json:"key" binding:"required,min=1,max=100"`
	Value   string `json:"value" binding:"required,min=1"`
	Explain string `json:"explain" binding:"required,min=1"`
}

type ParamsUpdateConfig = ParamsAddConfig

type UriUpdateConfig struct {
	Key string `uri:"key" binding:"required,min=1,max=100"`
}

type UriGetConfigItem = UriUpdateConfig

type UriDeleteConfig = UriUpdateConfig

type ParamsCreatePageView struct {
	Url string `json:"url" binding:"required,url"`
}

type ParamsAddFriend struct {
	Name   string `json:"name" binding:"required,min=1,max=50"`
	Url    string `json:"url" binding:"required,url"`
	Avatar string `json:"avatar" binding:"required,url"`
	Intro  string `json:"intro" binding:"required,min=1,max=128"`
}

type UriUpdateFriend struct {
	Id uint64 `uri:"id" binding:"required"`
}

type UriDeleteFriend = UriUpdateTag
