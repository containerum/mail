package storages

import (
	"os"

	"time"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

type MessagesStorage struct {
	db  *bolt.DB
	log *logrus.Logger
}

type MessagesStorageValue struct {
	UserId       string            `json:"user_id"`
	TemplateName string            `json:"template_name"`
	Variables    map[string]string `json:"variables,omitempty"`
	CreatedAt    time.Time         `json:"created_at"` // UTC
	Message      string            `json:"message"`    // base64
}

func NewMessageStorage(file string, options *bolt.Options) (*MessagesStorage, error) {
	log := logrus.WithField("component", "messages_storage").Logger
	log.Infof("Opening storage at %s with options %#v", file, options)
	db, err := bolt.Open(file, os.ModePerm, options)
	if err != nil {
		log.WithError(err).Error("Failed to open storage")
		return nil, err
	}
	return &MessagesStorage{
		db:  db,
		log: log,
	}, nil
}
