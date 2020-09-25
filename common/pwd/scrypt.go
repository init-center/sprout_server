package pwd

import (
	"encoding/hex"
	"sprout_server/settings"
	"strings"

	"golang.org/x/crypto/scrypt"
)

func Encrypt(password string, uid string) (string, error) {
	// lowercase the uid
	uid = strings.ToLower(uid)
	salt := settings.Conf.SundriesConfig.SaltPrefix + uid
	e, err := scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
	encryptPwd := hex.EncodeToString(e)
	if err != nil {
		return encryptPwd, err
	}
	return encryptPwd, nil
}
