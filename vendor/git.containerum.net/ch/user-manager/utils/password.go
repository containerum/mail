package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

const pwdIteration = 30
const keyLen = 32

func GenSalt(args ...string) string {
	timeSalt := strconv.FormatInt(time.Now().UnixNano(), 10)
	args = append(args, timeSalt)

	randomByteSalt := make([]byte, 10)
	rand.Read(randomByteSalt)
	args = append(args, string(randomByteSalt))

	resultSalt := make([]byte, 0)

	for i := len(args) - 1; i >= 0; i-- { // More random data goes first
		t := sha256.Sum256(append(resultSalt, []byte(args[i])...))
		resultSalt = t[:]
	}
	return base64.StdEncoding.EncodeToString(resultSalt)
}

func GetByteKey(pwd, salt string) []byte {
	return pbkdf2.Key([]byte(pwd), []byte(salt), pwdIteration, keyLen, sha256.New)
}

func GetKey(pwd, salt string) string {
	bKey := GetByteKey(pwd, salt)
	return base64.StdEncoding.EncodeToString(bKey)
}

func CheckPassword(pwd, salt, key string) bool {
	if key == GetKey(pwd, salt) {
		return true
	}
	return false
}
