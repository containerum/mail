package bolt

import (
	"os"

	mttypes "git.containerum.net/ch/json-types/mail-templater"
	"git.containerum.net/ch/mail-templater/storages"
	"github.com/boltdb/bolt"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

type boltMessagesStorage struct {
	db  *bolt.DB
	log *logrus.Entry
}

var _ storages.MessagesStorage = &boltMessagesStorage{}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const boltMessagesStorageBucket = "messages"

func NewBoltMessagesStorage(file string, options *bolt.Options) (storages.MessagesStorage, error) {
	log := logrus.WithField("component", "messages_storage")
	log.Infof("Opening storage at %s with options %#v", file, options)
	db, err := bolt.Open(file, os.ModePerm, options)
	if err != nil {
		log.WithError(err).Errorln("Failed to open storage")
		return nil, err
	}

	log.Infof("Creating bucket %s", boltMessagesStorageBucket)
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(boltMessagesStorageBucket))
		return err
	})
	if err != nil {
		log.WithError(err).Errorln("Create bucket failed")
		return nil, err
	}

	return &boltMessagesStorage{
		db:  db,
		log: log,
	}, nil
}

func (s *boltMessagesStorage) PutValue(id string, value *mttypes.MessagesStorageValue) error {
	loge := s.log.WithFields(logrus.Fields{
		"id":    id,
		"value": value,
	})
	loge.Infof("Putting value")
	return s.db.Update(func(tx *bolt.Tx) error {
		loge.Debugln("Get bucket")
		b := tx.Bucket([]byte(boltMessagesStorageBucket))

		loge.Debugln("Marshal json")
		valueB, err := json.Marshal(value)
		if err != nil {
			loge.WithError(err).Errorln("Error marshalling value")
		}
		return b.Put([]byte(id), valueB)
	})
}

func (s *boltMessagesStorage) GetValue(id string) (*mttypes.MessagesStorageValue, error) {
	loge := s.log.WithField("id", id)
	loge.Infof("Getting value")
	var value mttypes.MessagesStorageValue
	err := s.db.View(func(tx *bolt.Tx) error {
		loge.Debugln("Get bucket")
		b := tx.Bucket([]byte(boltMessagesStorageBucket))

		loge.Debugln("Extract value from storage")
		valueB := b.Get([]byte(id))
		if valueB == nil {
			loge.Infoln("Cannot find value")
			return mttypes.ErrMessageNotExists
		}

		if err := json.Unmarshal(valueB, &value); err != nil {
			loge.WithError(err).Errorln("Value unmarshal failed")
			return err
		}

		return nil
	})
	return &value, err
}

func (s *boltMessagesStorage) Close() error {
	return s.db.Close()
}
