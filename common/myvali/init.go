package myvali

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func Init() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("checkPwd", checkPwd); err != nil {
			zap.L().Error("register pwd validator failed", zap.Error(err))
			return err
		}

		if err := v.RegisterValidation("checkUid", checkUid); err != nil {
			zap.L().Error("register uid validator failed", zap.Error(err))
			return err
		}

		if err := v.RegisterValidation("checkName", checkName); err != nil {
			zap.L().Error("register name validator failed", zap.Error(err))
			return err
		}
	}
	return nil
}
