package storages

import (
	"os"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

type MessageStorage struct {
	db  *bolt.DB
	log *logrus.Logger
}

func NewMessageStorage(file string, options *bolt.Options) (*MessageStorage, error) {
	log := logrus.WithField("component", "template_storage").Logger
	log.Infof("Opening storage at %s with options %#v", file, options)
	db, err := bolt.Open(file, os.ModePerm, options)
	if err != nil {
		log.WithError(err).Error("Failed to open storage")
		return nil, err
	}
	return &MessageStorage{
		db:  db,
		log: log,
	}, nil
}
