package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base32"
	"encoding/hex"
)

// GenerateRandomString - generates random string
func GenerateRandomString(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return base32.StdEncoding.EncodeToString(b)
}

// GetMD5Hash - makes hash of string
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
