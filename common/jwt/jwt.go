package jwt

import (
	"errors"
	"sprout_server/common/constant"
	"sprout_server/settings"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	Uid string `json:"uid"`
	jwt.StandardClaims
}

func GenToken(uid string) (string, error) {
	mc := MyClaims{
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "init.center",
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Add(constant.TokenExpireDuration).Unix(),
		},
	}
	secret := genSecret()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mc)
	return token.SignedString(secret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return genSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid { // verify token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}

func genSecret() []byte {
	//secret must be a []byte
	return []byte(strings.ToLower(settings.Conf.SundriesConfig.JwtSecret))
}
