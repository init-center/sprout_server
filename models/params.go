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
	//Uid   string `json:"uid" binding:"required"`
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
	Uid         string  `json:"uid" binding:"required,min=2,max=12,checkUid"`
	Category    int64   `json:"category" binding:"required"`
	Tags        []int64 `json:"tags" binding:"required,min=1,max=20"`
	Title       string  `json:"title" binding:"required,min=1,max=128"`
	Cover       string  `json:"cover" binding:"required,url"`
	Bgm         string  `json:"bgm" binding:"required,url"`
	Summary     string  `json:"summary" binding:"required,max=128"`
	Content     string  `json:"content" binding:"required"`
	Display     int8    `json:"display" binding:"oneof=0 1"`
	CommentOpen int8    `json:"commentOpen" binding:"oneof=0 1"`
	Top         int8    `json:"top" binding:"oneof=0 1"`
}

type QueryStringGetPostList struct {
	Page  int64 `form:"page,default=1",binding:"gte=1"`
	Limit int64 `form:"limit,default=10",binding:"gte=1"`
}

type UriGetPostDetail struct {
	Pid int64 `uri:"pid" binding:"required"`
}

type ParamsAddComment struct {
	Pid       int64  `uri:"pid"`
	Uid       string `json:"uid"`
	TargetCid int64  `json:"targetCid,string"`
	Content   string `json:"content" binding:"required,min=1,max=1024"`
}

type ParamsGetCommentList struct {
	Pid        int64 `uri:"pid"`
	Page       int64 `form:"page,default=1",binding:"gte=1"`
	Limit      int64 `form:"limit,default=5",binding:"gte=1"`
	ChildLimit int64 `form:"child_limit,default=2",binding:"gte=1"`
}

type ParamsGetParentCommentChildren struct {
	Pid   int64 `uri:"pid"`
	Cid   int64 `uri:"cid"`
	Page  int64 `form:"page,default=1",binding:"gte=1"`
	Limit int64 `form:"limit,default=5",binding:"gte=1"`
}

type ParamsPostFavorite struct {
	Pid int64 `uri:"pid"`
	Uid string
}
