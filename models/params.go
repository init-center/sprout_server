package models

type ParamsSignUp struct {
	Uid        string `json:"uid" binding:"required,min=2,max=12,checkUid"`
	Name       string `json:"name" binding:"required,checkName"`
	Password   string `json:"password" binding:"required,min=6,max=16,checkPwd"`
	RePassword string `json:"rePassword" binding:"required,eqfield=Password,min=6,max=16,checkPwd"`
	Email      string `json:"email" binding:"required,email"`
	ECode      string `json:"eCode" binding:"required"` // email verify code
}

type ParamsGetECode struct {
	Uid   string `json:"uid" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Type  string `json:"type" binding:"required"`
}

type ParamsSignIn struct {
	Uid      string `json:"uid" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=16,checkPwd"`
}
