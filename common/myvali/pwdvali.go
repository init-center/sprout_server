package myvali

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

func checkPwd(fl validator.FieldLevel) bool {
	pwd := fl.Field().String()
	var hasNumber, hasUpperCase, hasLowercase bool
	for _, c := range pwd {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpperCase = true
		case unicode.IsLower(c):
			hasLowercase = true
		}
	}
	return hasNumber && hasUpperCase && hasLowercase
}
