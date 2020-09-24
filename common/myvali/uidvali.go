package myvali

import (
	"regexp"

	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"
)

func checkUid(fl validator.FieldLevel) bool {
	uid := fl.Field().String()
	pat := `^[A-Za-z][\-_]?([A-Za-z0-9]+[\-_]?)*[A-Za-z0-9]$`
	match, err := regexp.MatchString(pat, uid)
	if err != nil {
		zap.L().Error("check uid failed", zap.Error(err))
	}
	return match
}
