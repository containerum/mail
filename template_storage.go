package ch_mail_templater

import (
	"errors"
	"os"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

// TemplateStorage is a storage of email templates based on boltDB. It has versioning support.
type TemplateStorage struct {
	db *bolt.DB
}

var log = logrus.WithField("component", "template_storage").Logger

var (
	ErrTemplateNotExists = errors.New("specified template not exists in storage")
	ErrVersionNotExists  = errors.New("specified version not exists in storage")
)

func NewTemplateStorage(file string, options *bolt.Options) (*TemplateStorage, error) {
	log.Infof("Opening storage at %s with options %#v", file, options)
	db, err := bolt.Open(file, os.ModePerm, options)
	if err != nil {
		log.WithError(err).Error("Failed to open storage")
		return nil, err
	}
	return &TemplateStorage{
		db: db,
	}, err
}

// PutTemplate puts template to storage. If template with specified name and version exists it will be overwritten.
func (s *TemplateStorage) PutTemplate(templateName, templateVersion, templateData string) error {
	loge := log.WithFields(logrus.Fields{
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
		if err := b.Put([]byte(templateVersion), []byte(templateData)); err != nil {
			loge.WithError(err).Error("Put kv data failed")
			return err
		}

		return nil
	})
}

// GetTemplate returns specified version of template.
func (s *TemplateStorage) GetTemplate(templateName, templateVersion string) (string, error) {
	loge := log.WithFields(logrus.Fields{
		"name":    templateName,
		"version": templateVersion,
	})
	loge.Info("Trying to get template")

	var template string
	err := s.db.View(func(tx *bolt.Tx) error {
		loge.Debug("Getting bucket")
		b := tx.Bucket([]byte(templateName))
		if b == nil {
			loge.Error("Cannot find bucket")
			return ErrTemplateNotExists
		}

		loge.Debug("Getting value")
		templateB := b.Get([]byte(templateVersion))
		if templateB == nil {
			loge.Error("Cannot find version")
			return ErrVersionNotExists
		}
		template = string(templateB)

		return nil
	})

	return template, err
}

// GetTemplates returns all versions of templates in map (key is version, value is template).
func (s *TemplateStorage) GetTemplates(templateName string) (map[string]string, error) {
	loge := log.WithField("name", templateName)
	loge.Info("Trying to get all versions of template")

	templates := make(map[string]string)
	err := s.db.View(func(tx *bolt.Tx) error {
		loge.Debug("Getting bucket")
		b := tx.Bucket([]byte(templateName))
		if b == nil {
			loge.Error("Cannot find bucket")
			return ErrTemplateNotExists
		}

		loge.Debugf("Iterating over bucket")
		err := b.ForEach(func(k, v []byte) error {
			loge.Debugf("Handling version %s", k)
			templates[string(k)] = string(v)
			return nil
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
	loge := log.WithFields(logrus.Fields{
		"name":    templateName,
		"version": templateVersion,
	})
	loge.Info("Trying to delete template")

	return s.db.Update(func(tx *bolt.Tx) error {
		loge.Debug("Getting bucket")
		b := tx.Bucket([]byte(templateName))
		if b == nil {
			loge.Error("Cannot find bucket")
			return ErrTemplateNotExists
		}

		loge.Debug("Deleting entry")
		// check if entry exists
		if v := b.Get([]byte(templateVersion)); v == nil {
			loge.Error("Cannot find version")
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
	loge := log.WithField("name", templateName)
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
