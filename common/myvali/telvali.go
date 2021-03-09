package myvali

import (
	"regexp"

	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"
)

func checkTel(fl validator.FieldLevel) bool {
	uid := fl.Field().String()
	pat := `^1[3-9]\d{9}$`
	match, err := regexp.MatchString(pat, uid)
	if err != nil {
		zap.L().Error("check tel failed", zap.Error(err))
	}
	return match
}
