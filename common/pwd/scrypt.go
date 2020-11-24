package pwd

import (
	"encoding/hex"
	"sprout_server/settings"

	"golang.org/x/crypto/scrypt"
)

func Encrypt(password string) (string, error) {
	salt := settings.Conf.SundriesConfig.Salt
	e, err := scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
	encryptPwd := hex.EncodeToString(e)
	if err != nil {
		return encryptPwd, err
	}
	return encryptPwd, nil
}
