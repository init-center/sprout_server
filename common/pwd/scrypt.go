package pwd

import (
	"encoding/hex"

	"golang.org/x/crypto/scrypt"
)

func Encrypt(password string, salt string) (string, error) {
	e, err := scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
	encryptPwd := hex.EncodeToString(e)
	if err != nil {
		return encryptPwd, err
	}
	return encryptPwd, nil
}
