package vcode

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"net/smtp"
	"sprout_server/common/constant"
	"sprout_server/common/response/code"
	"sprout_server/dao/redis"
	"sprout_server/models"
	"sprout_server/settings"
	"strings"
	"time"

	"github.com/jordan-wright/email"
	"go.uber.org/zap"
)

func SendCodeToEmail(p *models.ParamsGetECode) int {
	// 1.Check whether the eCode has expired
	_, err := redis.GetECode(p.Uid)
	if err == redis.Nil {
		//Nil value, can send mail
		count, err := redis.GetECodeCount(p.Uid)
		if err != nil && err != redis.Nil {
			zap.L().Error("get ecode count failed", zap.Error(err))
			return code.CodeServerBusy
		} else if count < constant.MaxEcodeCount {
			// gen ecode
			eCode := genECode()
			// send code to email
			if err := send(p, eCode); err != nil {
				return code.CodeServerBusy
			}
			// send success
			// store the ecode to redis
			_, err := redis.SetECode(p.Uid, eCode, constant.ECodeExpireTime*time.Minute)
			if err != nil {
				zap.L().Error("store the ecode to rdb failed", zap.Error(err))
				return code.CodeServerBusy
			}
			// let the send count add 1 (max 5 times every day)
			_, err = redis.IncrECodeCount(p.Uid)
			if err != nil {
				zap.L().Error("incr the ecode failed", zap.Error(err))
			}
			return code.CodeOK
		} else {
			return code.CodeExceedMaxCountInADay
		}
	} else if err != nil {
		// db error
		zap.L().Error("get ecode failed", zap.Error(err))
		return code.CodeServerBusy

	} else {
		return code.CodeFrequentRequests
	}
}

func send(p *models.ParamsGetECode, eCode string) error {
	tmpl, err := template.ParseFiles("./template/email_code.tmpl")
	if err != nil {
		zap.L().Error("get tmpl failed", zap.Error(err))
		return err
	}
	data := &models.ECodeData{
		Email: p.Email,
		Type:  p.Type,
		ECode: eCode,
		Time:  time.Now().Format("2006-01-02 15:04:05"),
	}
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		zap.L().Error("execute tmpl failed", zap.Error(err))
		return err
	}
	e := &email.Email{
		To:      []string{p.Email},
		From:    settings.Conf.SmtpConfig.UserName,
		Subject: "init.center verification code",
		HTML:    []byte(tpl.String()),
	}
	smtpConf := settings.Conf.SmtpConfig
	if err := e.Send(smtpConf.Host+":"+smtpConf.Port, smtp.PlainAuth("", smtpConf.User, smtpConf.Password, smtpConf.Host)); err != nil {
		zap.L().Error("send eCode to email failed", zap.Error(err))
		return err
	}
	return nil
}

func genECode() string {
	rand.Seed(time.Now().UnixNano())
	randFloat := rand.Float64()
	c := fmt.Sprintf("%x", randFloat)
	return strings.ToUpper(c[6:10])
}
