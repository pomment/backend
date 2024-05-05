package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"github.com/pomment/pomment/config"
)

func GetMailHash(text string) string {
	if config.Content.Avatar.UseSha256 {
		hash := sha256.Sum256([]byte(text))
		return hex.EncodeToString(hash[:])
	}
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
