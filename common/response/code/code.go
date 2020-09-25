package code

import "net/http"

const (
	CodeOK = 2000 + iota
	CodeInvalidParams
	CodeUserIdExist
	CodeUserNameExist
	CodeEmailExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	CodeFrequentRequests
	CodeNeedLogin
	CodeInvalidToken
	CodeExceedMaxCountInADay
	CodeECodeExpired
	CodeIncorrectECode
)

var CodeMessageMap = map[int]string{
	CodeOK:               "OK",
	CodeInvalidParams:    "Invalid params",
	CodeUserIdExist:      "User id already exist",
	CodeUserNameExist:    "UserName already exist",
	CodeEmailExist:       "Email already exist",
	CodeUserNotExist:     "User not exist",
	CodeInvalidPassword:  "Incorrect user id or password",
	CodeServerBusy:       "Server busy",
	CodeFrequentRequests: "Frequent requests",

	CodeNeedLogin:            "Need login",
	CodeInvalidToken:         "Invalid token",
	CodeExceedMaxCountInADay: "Exceed the maximum limit of a single day",
	CodeECodeExpired:         "Email verification code expired or no verification code",
	CodeIncorrectECode:       "Incorrect Email verification code",
}

// Error code mapping HTTP status code
var CodeHCodeMap = map[int]int{
	CodeOK:               http.StatusOK,
	CodeInvalidParams:    http.StatusUnprocessableEntity,
	CodeUserIdExist:      http.StatusConflict,
	CodeUserNameExist:    http.StatusConflict,
	CodeEmailExist:       http.StatusConflict,
	CodeUserNotExist:     http.StatusNotFound,
	CodeInvalidPassword:  http.StatusUnauthorized,
	CodeServerBusy:       http.StatusInternalServerError,
	CodeFrequentRequests: http.StatusTooManyRequests,

	CodeNeedLogin:            http.StatusUnauthorized,
	CodeInvalidToken:         http.StatusUnauthorized,
	CodeExceedMaxCountInADay: http.StatusTooManyRequests,
	CodeECodeExpired:         http.StatusUnprocessableEntity,
	CodeIncorrectECode:       http.StatusUnprocessableEntity,
}

func Msg(code int) string {
	msg, ok := CodeMessageMap[code]
	if !ok {
		return CodeMessageMap[CodeServerBusy]
	}

	return msg
}

func HCode(code int) int {
	hCode, ok := CodeHCodeMap[code]
	if !ok {
		return CodeHCodeMap[CodeServerBusy]
	}

	return hCode
}
