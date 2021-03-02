package constants

import (
	"time"
)

const (
	MaxEcodeCount       = 5
	ECodeExpireTime     = 2
	TokenExpireDuration = time.Hour * 24
)

const CtxUidKey = "uid"

// user group constant
const (
	UserGroupAdmin = iota + 1
	UserGroupDefault
)

const (
	EcodeSignUpType = iota + 1
)

var EcodeTypeNameMap = map[int]string{
	EcodeSignUpType: "注册",
}
