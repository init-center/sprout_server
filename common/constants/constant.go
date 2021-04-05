package constants

import (
	"time"
)

const (
	MaxEcodeCount       = 5
	ECodeExpireTime     = 2
	TokenExpireDuration = time.Hour * 24
)

const (
	CtxUidKey           = "uid"
	CtxOriginEngineKey  = "engine"
	CtxOriginOsKey      = "origin-os"
	CtxOriginBrowserKey = "origin-browser"
	CtxOriginUAKey      = "origin-user-agent"
	CtxOriginIpKey      = "origin-ip"
)

// user group constant
const (
	UserGroupAdmin = iota + 1
	UserGroupDefault
)

const (
	EcodeSignUpType = iota + 1
	EcodeUpdateUserType
)

var EcodeTypeNameMap = map[int]string{
	EcodeSignUpType:     "注册",
	EcodeUpdateUserType: "修改用户信息",
}

const ShouldReplyCommentChildLimit = 5
