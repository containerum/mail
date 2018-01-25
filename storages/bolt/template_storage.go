package bolt

import (
	"os"

	"time"

	mttypes "git.containerum.net/ch/json-types/mail-templater"
	"git.containerum.net/ch/mail-templater/storages"
	"github.com/blang/semver"
	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

type boltTemplateStorage struct {
	db  *bolt.DB
	log *logrus.Entry
}

// NewBoltTemplateStorage returns BoltDB-based template storage.
// Supports template tagging and versioning with "semantic versioning 2.0" standard
func NewBoltTemplateStorage(file string, options *bolt.Options) (storages.TemplateStorage, error) {
	log := logrus.WithField("component", "template_storage")
	log.Infof("Opening storage at %s with options %#v", file, options)
	db, err := bolt.Open(file, os.ModePerm, options)
	if err != nil {
		log.WithError(err).Errorln("Failed to open storage")
		return nil, err
	}
	return &boltTemplateStorage{
		db:  db,
		log: log,
	}, err
}

func (s *boltTemplateStorage) PutTemplate(templateName, templateVersion, templateData, templateSubject string) error {
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

func (s *boltTemplateStorage) GetTemplate(templateName, templateVersion string) (*mttypes.TemplateStorageValue, error) {
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

func (s *boltTemplateStorage) GetLatestVersionTemplate(templateName string) (string, *mttypes.TemplateStorageValue, error) {
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
		var latestVerStr string
		err := b.ForEach(func(k, v []byte) error {
			loge.Debugf("Handling version %s", k)
			ver, err := semver.ParseTolerant(string(k))
			if err != nil {
				loge.WithError(err).Debugf("skipping %s", k)
				return nil // skip non-semver keys
			}
			if ver.GT(latestVer) {
				latestVer = ver
				latestVerStr = string(k)
			}
			return nil
		})
		if err != nil {
			loge.WithError(err).Errorln("Iterating error")
		}

		loge.Debugf("Extracting latest version %v", latestVerStr)
		templateB := b.Get([]byte(latestVerStr))
		if templateB == nil {
			loge.Infof("Cannot find version %v", latestVerStr)
			return mttypes.ErrVersionNotExists
		}
		return json.Unmarshal(templateB, &templateValue)
	})

	return templateVersion, &templateValue, err
}

func (s *boltTemplateStorage) GetTemplates(templateName string) (map[string]*mttypes.TemplateStorageValue, error) {
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

func (s *boltTemplateStorage) DeleteTemplate(templateName, templateVersion string) error {
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

func (s *boltTemplateStorage) DeleteTemplates(templateName string) error {
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

func (s *boltTemplateStorage) GetTemplatesList() (*mttypes.TemplatesListResponse, error) {
	loge := s.log.WithField("name", "templates list")
	loge.Infoln("Trying to get list of all templates")

	resp := mttypes.TemplatesListResponse{
		Templates: []mttypes.TemplatesListEntry{},
	}
	err := s.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) (err error) {
			template, err := s.GetTemplates(string(name))
			if err != nil {
				return err
			}

			var versions []string
			for version, _ := range template {
				versions = append(versions, version)
			}

			resp.Templates = append(resp.Templates, mttypes.TemplatesListEntry{
				Name:     string(name),
				Versions: versions,
			})
			return nil

		})
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *boltTemplateStorage) Close() error {
	return s.db.Close()
}
