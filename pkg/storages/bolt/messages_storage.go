package bolt

import (
	"errors"
	"os"

	"git.containerum.net/ch/mail-templater/pkg/models"
	"git.containerum.net/ch/mail-templater/pkg/mtErrors"
	"git.containerum.net/ch/mail-templater/pkg/storages"
	"github.com/boltdb/bolt"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

type boltMessagesStorage struct {
	db  *bolt.DB
	log *logrus.Entry
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const boltMessagesStorageBucket = "messages"

const iterationFinished = "Iteration finished"

// NewBoltMessagesStorage BoltDB-based messages storage.
// Used for storing sent messages
func NewBoltMessagesStorage(file string, options *bolt.Options) (storages.MessagesStorage, error) {
	log := logrus.WithField("component", "messages_storage")
	log.Infof("Opening storage at %s with options %#v", file, options)
	db, err := bolt.Open(file, os.ModePerm, options)
	if err != nil {
		log.WithError(err).Errorln(errStorageOpenFailed)
		return nil, err
	}

	log.Infof("Creating bucket %s", boltMessagesStorageBucket)
	err = db.Update(func(tx *bolt.Tx) error {
		_, txerr := tx.CreateBucketIfNotExists([]byte(boltMessagesStorageBucket))
		return txerr
	})
	if err != nil {
		log.WithError(err).Errorln(errStorageOpenFailed)
		return nil, err
	}

	return &boltMessagesStorage{
		db:  db,
		log: log,
	}, nil
}

// Close
// Closes bolt storage
func (s *boltMessagesStorage) Close() error {
	return s.db.Close()
}

// GetMessageList
// Gets messages list from db
func (s *boltMessagesStorage) GetMessageList(page int, perPage int) (*models.MessageListResponse, error) {
	loge := s.log.WithField("name", "message list")
	loge.Infoln("Trying to get list of all messages")

	resp := models.MessageListResponse{
		Messages: []models.MessageListEntry{},
		MessageListQuery: &models.MessageListQuery{
			Page:    page,
			PerPage: perPage,
		},
	}
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(boltMessagesStorageBucket))

		startMessage := (page - 1) * perPage
		var messageNumber int

		err := b.ForEach(func(k, v []byte) error {

			if messageNumber >= startMessage+perPage {
				return errors.New(iterationFinished)
			}

			var value models.MessageListEntry

			if err := json.Unmarshal(v, &value); err != nil {
				loge.WithError(err).Errorln("Value unmarshal failed")
				return mtErrors.ErrUnableGetMessagesList()
			}

			if messageNumber >= startMessage {
				resp.Messages = append(resp.Messages, models.MessageListEntry{
					ID:           string(k),
					UserID:       value.UserID,
					TemplateName: value.TemplateName,
					CreatedAt:    value.CreatedAt,
				})
			}

			messageNumber++

			return nil
		})
		if err != nil {
			if err.Error() == iterationFinished {
				return nil
			}
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetMessage
// Gets message from db
func (s *boltMessagesStorage) GetMessage(id string) (*models.MessagesStorageValue, error) {
	loge := s.log.WithField("id", id)
	loge.Infof("Getting message")
	var value models.MessagesStorageValue
	err := s.db.View(func(tx *bolt.Tx) error {
		loge.Debugln("Get bucket")
		b := tx.Bucket([]byte(boltMessagesStorageBucket))

		loge.Debugln("Extract value from storage")
		valueB := b.Get([]byte(id))
		if valueB == nil {
			loge.Infoln("Cannot find value")
			return mtErrors.ErrMessageNotExist() //models.ErrMessageNotExists
		}

		if err := json.Unmarshal(valueB, &value); err != nil {
			loge.WithError(err).Errorln("Value unmarshal failed")
			return mtErrors.ErrUnableGetMessage()
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &value, nil
}

// PutMessage
// Saves message to db
func (s *boltMessagesStorage) PutMessage(id string, value *models.MessagesStorageValue) error {
	loge := s.log.WithFields(logrus.Fields{
		"id": id,
	})
	loge.Infof("Putting message")
	err := s.db.Update(func(tx *bolt.Tx) error {
		loge.Debugln("Get bucket")
		b := tx.Bucket([]byte(boltMessagesStorageBucket))

		loge.Debugln("Marshal json")
		valueB, err := json.Marshal(value)
		if err != nil {
			loge.WithError(err).Errorln("Error marshalling value")
			return mtErrors.ErrUnableSaveMessage()
		}
		return b.Put([]byte(id), valueB)
	})
	if err != nil {
		return err
	}
	return nil
}
