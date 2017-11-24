package storages

import (
	"os"

	"time"

	"errors"

	"github.com/boltdb/bolt"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

type MessagesStorage struct {
	db  *bolt.DB
	log *logrus.Logger
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const messagesStorageBucket = "messages"

type MessagesStorageValue struct {
	UserId       string            `json:"user_id"`
	TemplateName string            `json:"template_name"`
	Variables    map[string]string `json:"variables,omitempty"`
	CreatedAt    time.Time         `json:"created_at"` // UTC
	Message      string            `json:"message"`    // base64
}

var ErrMessageNotExists = errors.New("message not exists")

func NewMessagesStorage(file string, options *bolt.Options) (*MessagesStorage, error) {
	log := logrus.WithField("component", "messages_storage").Logger
	log.Infof("Opening storage at %s with options %#v", file, options)
	db, err := bolt.Open(file, os.ModePerm, options)
	if err != nil {
		log.WithError(err).Error("Failed to open storage")
		return nil, err
	}

	log.Infof("Creating bucket %s", messagesStorageBucket)
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(messagesStorageBucket))
		return err
	})
	if err != nil {
		log.WithError(err).Error("Create bucket failed")
		return nil, err
	}

	return &MessagesStorage{
		db:  db,
		log: log,
	}, nil
}

// PutValue puts MessageStorageValue to storage.
// If message with specified id already exists in storage it will be overwritten.
func (s *MessagesStorage) PutValue(id string, value *MessagesStorageValue) error {
	loge := s.log.WithFields(logrus.Fields{
		"id":    id,
		"value": value,
	})
	loge.Infof("Putting value")
	return s.db.Update(func(tx *bolt.Tx) error {
		loge.Debug("Get bucket")
		b := tx.Bucket([]byte(messagesStorageBucket))

		loge.Debug("Marshal json")
		valueB, err := json.Marshal(value)
		if err != nil {
			loge.WithError(err).Error("Error marshalling value")
		}
		return b.Put([]byte(id), valueB)
	})
}

// GetValue returns value by specified ID.
func (s *MessagesStorage) GetValue(id string) (*MessagesStorageValue, error) {
	loge := s.log.WithField("id", id)
	loge.Infof("Getting value")
	var value MessagesStorageValue
	err := s.db.View(func(tx *bolt.Tx) error {
		loge.Debug("Get bucket")
		b := tx.Bucket([]byte(messagesStorageBucket))

		loge.Debug("Extract value from storage")
		valueB := b.Get([]byte(id))
		if valueB == nil {
			loge.Info("Cannot find value")
			return ErrMessageNotExists
		}

		if err := json.Unmarshal(valueB, &value); err != nil {
			loge.WithError(err).Error("Value unmarshal failed")
			return err
		}

		return nil
	})
	return &value, err
}

func (s *MessagesStorage) Close() error {
	return s.db.Close()
}
