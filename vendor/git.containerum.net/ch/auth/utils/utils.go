package utils

import (
	"encoding/base64"

	"crypto/rand"
	"fmt"

	"strings"

	"git.containerum.net/ch/grpc-proto-files/common"
	"github.com/mssola/user_agent"
)

// ShortUserAgent generates short user agent from normal user agent using base64
func ShortUserAgent(userAgent string) string {
	ua := user_agent.New(userAgent)
	platform := ua.Platform()
	engine, _ := ua.Engine()
	os := ua.OS()
	browser, _ := ua.Browser()
	toEncode := strings.Join([]string{platform, os, engine, browser}, " ")
	return base64.StdEncoding.EncodeToString([]byte(toEncode))
}

func NewUUID() *common.UUID {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}

	return &common.UUID{
		Value: fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]),
	}
}

func UUIDEquals(a, b *common.UUID) bool {
	return a == b || a != nil && b != nil && a.Value == b.Value
}

func UUIDFromString(value string) *common.UUID {
	return &common.UUID{
		Value: value,
	}
}
