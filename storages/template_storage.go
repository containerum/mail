package storages

import (
	"os"

	"time"

	"strings"

	mttypes "git.containerum.net/ch/json-types/mail-templater"
	"github.com/blang/semver"
	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

// TemplateStorage is a storage of email templates based on boltDB. It has versioning support.
type TemplateStorage struct {
	db  *bolt.DB
	log *logrus.Entry
}

func NewTemplateStorage(file string, options *bolt.Options) (*TemplateStorage, error) {
	log := logrus.WithField("component", "template_storage")
	log.Infof("Opening storage at %s with options %#v", file, options)
	db, err := bolt.Open(file, os.ModePerm, options)
	if err != nil {
		log.WithError(err).Errorln("Failed to open storage")
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
	loge.Infoln("Putting template to storage")
	return s.db.Update(func(tx *bolt.Tx) error {
		loge.Debugln("Creating bucket")
		b, err := tx.CreateBucketIfNotExists([]byte(templateName))
		if err != nil {
			loge.WithError(err).Errorln("Bucket create failed")
			return err
		}

		loge.Debugln("Putting kv data")
		value, _ := json.Marshal(&mttypes.TemplateStorageValue{
			Data:      templateData,
			CreatedAt: time.Now().UTC(),
			Subject:   templateSubject,
		})
		if err := b.Put([]byte(templateVersion), value); err != nil {
			loge.WithError(err).Errorln("Put kv data failed")
			return err
		}

		return nil
	})
}

// GetTemplate returns specified version of template.
func (s *TemplateStorage) GetTemplate(templateName, templateVersion string) (*mttypes.TemplateStorageValue, error) {
	loge := s.log.WithFields(logrus.Fields{
		"name":    templateName,
		"version": templateVersion,
	})
	loge.Infoln("Trying to get template")

	var templateValue mttypes.TemplateStorageValue
	err := s.db.View(func(tx *bolt.Tx) error {
		loge.Debugln("Getting bucket")
		b := tx.Bucket([]byte(templateName))
		if b == nil {
			loge.Infoln("Cannot find bucket")
			return mttypes.ErrTemplateNotExists
		}

		loge.Debugln("Getting value")
		templateB := b.Get([]byte(templateVersion))
		if templateB == nil {
			loge.Infoln("Cannot find version")
			return mttypes.ErrVersionNotExists
		}
		return json.Unmarshal(templateB, &templateValue)
	})

	return &templateValue, err
}

// GetLatestVersionTemplate returns latest version of template and it`s value using semver to compare versions.
func (s *TemplateStorage) GetLatestVersionTemplate(templateName string) (string, *mttypes.TemplateStorageValue, error) {
	loge := s.log.WithField("name", templateName)
	loge.Infoln("Trying to get latest version of template")

	var templateValue mttypes.TemplateStorageValue
	var templateVersion string

	err := s.db.View(func(tx *bolt.Tx) error {
		loge.Debugln("Getting bucket")
		b := tx.Bucket([]byte(templateName))
		if b == nil {
			loge.Infoln("Cannot find bucket")
			return mttypes.ErrTemplateNotExists
		}

		loge.Debugf("Iterating over bucket")
		var latestVer semver.Version
		err := b.ForEach(func(k, v []byte) error {
			loge.Debugf("Handling version %s", k)
			ver, err := semver.Parse(strings.TrimPrefix(string(k), "v")) // make it working if version starts with "v"
			if err != nil {
				loge.WithError(err).Debugf("skipping %s", k)
				return nil // skip non-semver keys
			}
			if ver.GT(latestVer) {
				latestVer = ver
			}
			return nil
		})
		if err != nil {
			loge.WithError(err).Errorln("Iterating error")
		}

		templateVersion = "v" + latestVer.String()
		loge.Debugf("Extracting latest version %v", templateVersion)
		templateB := b.Get([]byte(templateVersion))
		if templateB == nil {
			loge.Infof("Cannot find version %v", templateVersion)
			return mttypes.ErrVersionNotExists
		}
		return json.Unmarshal(templateB, &templateValue)
	})

	return templateVersion, &templateValue, err
}

// GetTemplates returns all versions of templates in map (key is version, value is template).
func (s *TemplateStorage) GetTemplates(templateName string) (map[string]*mttypes.TemplateStorageValue, error) {
	loge := s.log.WithField("name", templateName)
	loge.Infoln("Trying to get all versions of template")

	templates := make(map[string]*mttypes.TemplateStorageValue)
	err := s.db.View(func(tx *bolt.Tx) error {
		loge.Debugln("Getting bucket")
		b := tx.Bucket([]byte(templateName))
		if b == nil {
			loge.Infoln("Cannot find bucket")
			return mttypes.ErrTemplateNotExists
		}

		loge.Debugf("Iterating over bucket")
		err := b.ForEach(func(k, v []byte) error {
			loge.Debugf("Handling version %s", k)
			var value mttypes.TemplateStorageValue
			err := json.Unmarshal(v, &value)
			templates[string(k)] = &value
			return err
		})
		if err != nil {
			loge.WithError(err).Errorln("Iterating error")
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
	loge.Infoln("Trying to delete template")

	return s.db.Update(func(tx *bolt.Tx) error {
		loge.Debugln("Getting bucket")
		b := tx.Bucket([]byte(templateName))
		if b == nil {
			loge.Infoln("Cannot find bucket")
			return mttypes.ErrTemplateNotExists
		}

		loge.Debugln("Deleting entry")
		// check if entry exists
		if v := b.Get([]byte(templateVersion)); v == nil {
			loge.Infoln("Cannot find version")
			return mttypes.ErrVersionNotExists
		}
		if err := b.Delete([]byte(templateVersion)); err != nil {
			loge.WithError(err).Errorln("Version delete failed")
			return err
		}

		return nil
	})
}

// DeleteTemplates deletes all versions of template. Returns nil on successful delete.
func (s *TemplateStorage) DeleteTemplates(templateName string) error {
	loge := s.log.WithField("name", templateName)
	loge.Infoln("Trying to delete all versions of template")

	return s.db.Update(func(tx *bolt.Tx) error {
		loge.Debugln("Deleting bucket")
		if err := tx.DeleteBucket([]byte(templateName)); err == bolt.ErrBucketNotFound {
			loge.WithError(err).Errorf("Bucket not found")
			return mttypes.ErrTemplateNotExists
		} else if err != nil {
			loge.WithError(err).Errorln("Bucket delete failed")
			return err
		}

		return nil
	})
}

func (s *TemplateStorage) Close() error {
	return s.db.Close()
}
