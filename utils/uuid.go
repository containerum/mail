package utils

import (
	"crypto/rand"
	"fmt"

	"github.com/sirupsen/logrus"
)

func NewUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		logrus.WithField("component", "uuid").WithError(err).Error("Generate uuid failed")
		return ""
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
