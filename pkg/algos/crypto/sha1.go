package ucrypto

import (
	"crypto/sha1"
	"encoding/hex"
)

func NewSHA1(data []byte) (string, error) {
	hash := sha1.New()
	if _, err := hash.Write(data); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
