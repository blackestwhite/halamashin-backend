package utils

import (
	"crypto/sha1"
	"encoding/base64"
)

func HashString(stringToHash string) string {
	salt := "salt is bad for health"
	hasher := sha1.New()
	hasher.Write([]byte(stringToHash + salt))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
