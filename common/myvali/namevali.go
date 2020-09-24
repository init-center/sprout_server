package myvali

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func checkName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	l := len([]rune(name))
	if l < 2 || l > 12 {
		return false
	}
	pat := `^[A-Za-z\p{Han}][\-_]?([A-Za-z0-9\p{Han}]+[\-_]?)*[A-Za-z0-9\p{Han}]$`
	match, err := regexp.MatchString(pat, name)
	if err != nil {
		zap.L().Error("check name failed", zap.Error(err))
	}
	return match
}
