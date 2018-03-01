package bolt

import (
	"errors"
	"os"

	mttypes "git.containerum.net/ch/json-types/mail-templater"
	cherry "git.containerum.net/ch/kube-client/pkg/cherry/mail-templater"
	"git.containerum.net/ch/mail-templater/pkg/model"
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
		log.WithError(err).Errorln(model.ErrStorageOpenFailed)
		return nil, err
	}

	log.Infof("Creating bucket %s", boltMessagesStorageBucket)
	err = db.Update(func(tx *bolt.Tx) error {
		_, txerr := tx.CreateBucketIfNotExists([]byte(boltMessagesStorageBucket))
		return txerr
	})
	if err != nil {
		log.WithError(err).Errorln(model.ErrStorageOpenFailed)
		return nil, err
	}

	return &boltMessagesStorage{
		db:  db,
		log: log,
	}, nil
}

func (s *boltMessagesStorage) Close() error {
	return s.db.Close()
}

func (s *boltMessagesStorage) PutMessage(id string, value *mttypes.MessagesStorageValue) error {
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
			return cherry.ErrUnableSaveMessage()
		}
		return b.Put([]byte(id), valueB)
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *boltMessagesStorage) GetMessage(id string) (*mttypes.MessagesStorageValue, error) {
	loge := s.log.WithField("id", id)
	loge.Infof("Getting message")
	var value mttypes.MessagesStorageValue
	err := s.db.View(func(tx *bolt.Tx) error {
		loge.Debugln("Get bucket")
		b := tx.Bucket([]byte(boltMessagesStorageBucket))

		loge.Debugln("Extract value from storage")
		valueB := b.Get([]byte(id))
		if valueB == nil {
			loge.Infoln("Cannot find value")
			return cherry.ErrMessageNotExist() //mttypes.ErrMessageNotExists
		}

		if err := json.Unmarshal(valueB, &value); err != nil {
			loge.WithError(err).Errorln("Value unmarshal failed")
			return cherry.ErrUnableGetMessage()
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func (s *boltMessagesStorage) GetMessageList(page int, perPage int) (*mttypes.MessageListResponse, error) {
	loge := s.log.WithField("name", "message list")
	loge.Infoln("Trying to get list of all messages")

	resp := mttypes.MessageListResponse{
		Messages: []mttypes.MessageListEntry{},
		MessageListQuery: &mttypes.MessageListQuery{
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

			var value mttypes.MessageListEntry

			if err := json.Unmarshal(v, &value); err != nil {
				loge.WithError(err).Errorln("Value unmarshal failed")
				return cherry.ErrUnableGetMessagesList()
			}

			if messageNumber >= startMessage {
				resp.Messages = append(resp.Messages, mttypes.MessageListEntry{
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
		loge.WithError(err) //.Errorln(mttypes.ErrMessagesGetFailed)
		return nil, err
	}

	return &resp, nil
}
