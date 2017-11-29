package storages

import (
	"errors"
	"os"

	"time"

	"github.com/blang/semver"
	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

// TemplateStorage is a storage of email templates based on boltDB. It has versioning support.
type TemplateStorage struct {
	db  *bolt.DB
	log *logrus.Logger
}

type TemplateStorageValue struct {
	Data      string    `json:"data"`
	Subject   string    `json:"template_subject"`
	CreatedAt time.Time `json:"created_at"` // UTC
}

var (
	ErrTemplateNotExists = errors.New("specified template not exists in storage")
	ErrVersionNotExists  = errors.New("specified version not exists in storage")
)

func NewTemplateStorage(file string, options *bolt.Options) (*TemplateStorage, error) {
	log := logrus.WithField("component", "template_storage").Logger
	log.Infof("Opening storage at %s with options %#v", file, options)
	db, err := bolt.Open(file, os.ModePerm, options)
	if err != nil {
		log.WithError(err).Error("Failed to open storage")
		return nil, err
	}
	return &TemplateStorage{
		db:  db,
		log: log,
	}, err
}

// PutTemplate puts template to storage. If template with specified name and version exists it will be overwritten.
func (s *TemplateStorage) PutTemplate(templateName, templateVersion, templateData, templateSubject string) error {
	loge := s.log.WithFields(logrus.Fields{
		"name":    templateName,
		"version": templateVersion,
	})
	loge.Info("Putting template to storage")
	return s.db.Update(func(tx *bolt.Tx) error {
		loge.Debug("Creating bucket")
		b, err := tx.CreateBucketIfNotExists([]byte(templateName))
		if err != nil {
			loge.WithError(err).Error("Bucket create failed")
			return err
		}

		loge.Debug("Putting kv data")
		value, _ := json.Marshal(&TemplateStorageValue{
			Data:      templateData,
			CreatedAt: time.Now().UTC(),
			Subject:   templateSubject,
		})
		if err := b.Put([]byte(templateVersion), value); err != nil {
			loge.WithError(err).Error("Put kv data failed")
			return err
		}

		return nil
	})
}

// GetTemplate returns specified version of template.
func (s *TemplateStorage) GetTemplate(templateName, templateVersion string) (*TemplateStorageValue, error) {
	loge := s.log.WithFields(logrus.Fields{
		"name":    templateName,
		"version": templateVersion,
	})
	loge.Info("Trying to get template")

	var templateValue TemplateStorageValue
	err := s.db.View(func(tx *bolt.Tx) error {
		loge.Debug("Getting bucket")
		b := tx.Bucket([]byte(templateName))
		if b == nil {
			loge.Info("Cannot find bucket")
			return ErrTemplateNotExists
		}

		loge.Debug("Getting value")
		templateB := b.Get([]byte(templateVersion))
		if templateB == nil {
			loge.Info("Cannot find version")
			return ErrVersionNotExists
		}
		return json.Unmarshal(templateB, &templateValue)
	})

	return &templateValue, err
}

// GetLatestVersionTemplate returns latest version of template and it`s value using semver to compare versions.
func (s *TemplateStorage) GetLatestVersionTemplate(templateName string) (string, *TemplateStorageValue, error) {
	loge := s.log.WithField("name", templateName)
	loge.Info("Trying to get latest version of template")

	var templateValue TemplateStorageValue
	var templateVersion string

	err := s.db.View(func(tx *bolt.Tx) error {
		loge.Debug("Getting bucket")
		b := tx.Bucket([]byte(templateName))
		if b == nil {
			loge.Info("Cannot find bucket")
			return ErrTemplateNotExists
		}

		loge.Debugf("Iterating over bucket")
		var latestVer semver.Version
		err := b.ForEach(func(k, v []byte) error {
			loge.Debugf("Handling version %s", k)
			ver, err := semver.Parse(string(k))
			if err != nil {
				return nil // skip non-semver keys
			}
			if ver.GT(latestVer) {
				latestVer = ver
			}
			return nil
		})
		if err != nil {
			loge.WithError(err).Error("Iterating error")
		}

		templateVersion = latestVer.String()
		loge.Debugf("Extracting latest version")
		templateB := b.Get([]byte(templateVersion))
		if templateB == nil {
			loge.Info("Cannot find version")
			return ErrVersionNotExists
		}
		return json.Unmarshal(templateB, &templateValue)
	})

	return templateVersion, &templateValue, err
}

// GetTemplates returns all versions of templates in map (key is version, value is template).
func (s *TemplateStorage) GetTemplates(templateName string) (map[string]*TemplateStorageValue, error) {
	loge := s.log.WithField("name", templateName)
	loge.Info("Trying to get all versions of template")

	templates := make(map[string]*TemplateStorageValue)
	err := s.db.View(func(tx *bolt.Tx) error {
		loge.Debug("Getting bucket")
		b := tx.Bucket([]byte(templateName))
		if b == nil {
			loge.Info("Cannot find bucket")
			return ErrTemplateNotExists
		}

		loge.Debugf("Iterating over bucket")
		err := b.ForEach(func(k, v []byte) error {
			loge.Debugf("Handling version %s", k)
			var value TemplateStorageValue
			err := json.Unmarshal(v, &value)
			templates[string(k)] = &value
			return err
		})
		if err != nil {
			loge.WithError(err).Error("Iterating error")
		}

		return nil
	})

	return templates, err
}

// DeleteTemplate deletes specified version of template. Returns nil on successful delete.
func (s *TemplateStorage) DeleteTemplate(templateName, templateVersion string) error {
	loge := s.log.WithFields(logrus.Fields{
		"name":    templateName,
		"version": templateVersion,
	})
	loge.Info("Trying to delete template")

	return s.db.Update(func(tx *bolt.Tx) error {
		loge.Debug("Getting bucket")
		b := tx.Bucket([]byte(templateName))
		if b == nil {
			loge.Info("Cannot find bucket")
			return ErrTemplateNotExists
		}

		loge.Debug("Deleting entry")
		// check if entry exists
		if v := b.Get([]byte(templateVersion)); v == nil {
			loge.Info("Cannot find version")
			return ErrVersionNotExists
		}
		if err := b.Delete([]byte(templateVersion)); err != nil {
			loge.WithError(err).Error("Version delete failed")
			return err
		}

		return nil
	})
}

// DeleteTemplates deletes all versions of template. Returns nil on successful delete.
func (s *TemplateStorage) DeleteTemplates(templateName string) error {
	loge := s.log.WithField("name", templateName)
	loge.Info("Trying to delete all versions of template")

	return s.db.Update(func(tx *bolt.Tx) error {
		loge.Debug("Deleting bucket")
		if err := tx.DeleteBucket([]byte(templateName)); err == bolt.ErrBucketNotFound {
			loge.WithError(err).Errorf("Bucket not found")
			return ErrTemplateNotExists
		} else if err != nil {
			loge.WithError(err).Error("Bucket delete failed")
			return err
		}

		return nil
	})
}

func (s *TemplateStorage) Close() error {
	return s.db.Close()
}
