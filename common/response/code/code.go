package code

import "net/http"

const (
	CodeOK = 2000 + iota
	CodeCreated
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
	CodePermissionDenied
	CodeCategoryExist
	CodeCategoryNotExist
	CodeTagExist
	CodeTagNotExist
	CodePostNotExist
	CodeCommentNotExist
	CodeAlreadyFavorited
	CodeNotFavorited
	CodeCantTopDeletePost
	CodeUserIsNotBaned
	CodeCategoryHasPost
	CodeTagHasPost
)

var MessageMap = map[int]string{
	CodeOK:               "成功",
	CodeCreated:          "创建成功",
	CodeInvalidParams:    "参数错误",
	CodeUserIdExist:      "ID已被使用",
	CodeUserNameExist:    "用户名已被使用",
	CodeEmailExist:       "Email已被使用",
	CodeUserNotExist:     "用户不存在",
	CodeInvalidPassword:  "用户名或密码错误",
	CodeServerBusy:       "服务器繁忙",
	CodeFrequentRequests: "请求频繁",

	CodeNeedLogin:            "需要登录",
	CodeInvalidToken:         "未登录或登录状态已过期",
	CodeExceedMaxCountInADay: "超过单日最大限制",
	CodeECodeExpired:         "验证码已过期或不存在验证码",
	CodeIncorrectECode:       "验证码不正确",
	CodePermissionDenied:     "拒绝访问",
	CodeCategoryExist:        "分类已存在",
	CodeCategoryNotExist:     "分类不存在",
	CodeTagExist:             "标签已存在",
	CodeTagNotExist:          "标签不存在",
	CodePostNotExist:         "文章不存在",
	CodeCommentNotExist:      "目标评论不存在",
	CodeAlreadyFavorited:     "已经喜欢",
	CodeNotFavorited:         "还未喜欢",
	CodeCantTopDeletePost:    "不能置顶删除的文章",
	CodeUserIsNotBaned:       "用户未被封禁",
	CodeCategoryHasPost:      "分类下存在文章",
	CodeTagHasPost:           "标签下存在文章",
}

// Error code mapping HTTP status code
var HCodeMap = map[int]int{
	CodeOK:               http.StatusOK,
	CodeCreated:          http.StatusCreated,
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
	CodePermissionDenied:     http.StatusForbidden,
	CodeCategoryExist:        http.StatusConflict,
	CodeCategoryNotExist:     http.StatusNotFound,
	CodeTagExist:             http.StatusConflict,
	CodeTagNotExist:          http.StatusNotFound,
	CodePostNotExist:         http.StatusNotFound,
	CodeCommentNotExist:      http.StatusNotFound,
	CodeAlreadyFavorited:     http.StatusNoContent,
	CodeNotFavorited:         http.StatusNotFound,
	CodeCantTopDeletePost:    http.StatusConflict,
	CodeUserIsNotBaned:       http.StatusConflict,
	CodeCategoryHasPost:      http.StatusConflict,
	CodeTagHasPost:           http.StatusConflict,
}

func Msg(code int) string {
	msg, ok := MessageMap[code]
	if !ok {
		return MessageMap[CodeServerBusy]
	}

	return msg
}

func HCode(code int) int {
	hCode, ok := HCodeMap[code]
	if !ok {
		return HCodeMap[CodeServerBusy]
	}

	return hCode
}
