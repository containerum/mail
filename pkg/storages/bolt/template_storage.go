package bolt

import (
	"os"
	"time"

	"git.containerum.net/ch/mail-templater/pkg/models"
	"git.containerum.net/ch/mail-templater/pkg/mterrors"
	"git.containerum.net/ch/mail-templater/pkg/storages"
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
		log.WithError(err).WithError(errStorageOpenFailed)
		return nil, err
	}
	return &boltTemplateStorage{
		db:  db,
		log: log,
	}, nil
}

// Close
// Closes bolt storage
func (s *boltTemplateStorage) Close() error {
	return s.db.Close()
}

// GetTemplatesList
// Gets templates list
func (s *boltTemplateStorage) GetTemplatesList() (*models.TemplatesListResponse, error) {
	loge := s.log.WithField("name", "templates list")
	loge.Infoln("Trying to get list of all templates")

	resp := models.TemplatesListResponse{
		Templates: []models.TemplatesListEntry{},
	}
	err := s.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) (err error) {
			template, err := s.GetTemplates(string(name))
			if err != nil {
				return err
			}

			var versions []string
			for version := range template {
				versions = append(versions, version)
			}

			resp.Templates = append(resp.Templates, models.TemplatesListEntry{
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

// GetTemplates
// Gets all versions of specific template
func (s *boltTemplateStorage) GetTemplates(templateName string) (map[string]*models.Template, error) {
	loge := s.log.WithField("name", templateName)
	loge.Infoln("Trying to get all versions of template")

	templates := make(map[string]*models.Template)
	err := s.db.View(func(tx *bolt.Tx) error {
		loge.Debugln("Getting bucket")
		b := tx.Bucket([]byte(templateName))
		if b == nil {
			loge.Infoln("Cannot find bucket")
			return mterrors.ErrTemplateNotExist()
		}

		loge.Debugf("Iterating over bucket")
		err := b.ForEach(func(k, v []byte) error {
			loge.Debugf("Handling version %s", k)
			var value models.Template
			err := json.Unmarshal(v, &value)
			templates[string(k)] = &value
			return err
		})
		if err != nil {
			loge.WithError(err).Errorln("Iterating error")
			return mterrors.ErrUnableGetTemplatesList()
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// GetTemplate
// Gets specific version of specific template
func (s *boltTemplateStorage) GetTemplate(templateName, templateVersion string) (*models.Template, error) {
	loge := s.log.WithFields(logrus.Fields{
		"name":    templateName,
		"version": templateVersion,
	})
	loge.Infoln("Trying to get template")

	var templateValue models.Template
	err := s.db.View(func(tx *bolt.Tx) error {
		loge.Debugln("Getting bucket")
		b := tx.Bucket([]byte(templateName))
		if b == nil {
			loge.Infoln("Cannot find bucket")
			return mterrors.ErrTemplateNotExist() //models.ErrTemplateNotExists
		}

		loge.Debugln("Getting value")
		templateB := b.Get([]byte(templateVersion))
		if templateB == nil {
			loge.Infoln("Cannot find version")
			return mterrors.ErrTemplateVersionNotExist() //models.ErrTemplateVersionNotExists
		}
		return json.Unmarshal(templateB, &templateValue)
	})

	if err != nil {
		return nil, err
	}
	return &templateValue, nil
}

// GetTemplate
// Gets latest version of specific template
func (s *boltTemplateStorage) GetLatestVersionTemplate(templateName string) (*string, *models.Template, error) {
	loge := s.log.WithField("name", templateName)
	loge.Infoln("Trying to get latest version of template")

	var templateValue models.Template
	var templateVersion string

	err := s.db.View(func(tx *bolt.Tx) error {
		loge.Debugln("Getting bucket")
		b := tx.Bucket([]byte(templateName))
		if b == nil {
			loge.Infoln("Cannot find bucket")
			return mterrors.ErrTemplateNotExist() //models.ErrTemplateNotExists
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
			return err
		}

		loge.Debugf("Extracting latest version %v", latestVerStr)
		templateB := b.Get([]byte(latestVerStr))
		if templateB == nil {
			loge.Infof("Cannot find version %v", latestVerStr)
			return mterrors.ErrTemplateVersionNotExist() //models.ErrTemplateVersionNotExists
		}
		return json.Unmarshal(templateB, &templateValue)
	})
	if err != nil {
		return nil, nil, err
	}

	return &templateVersion, &templateValue, nil
}

// PutTemplate
// Saves template to db
func (s *boltTemplateStorage) PutTemplate(templateName, templateVersion, templateData, templateSubject string, new bool) error {
	loge := s.log.WithFields(logrus.Fields{
		"name":    templateName,
		"version": templateVersion,
	})
	loge.Debugln("Putting template to storage")
	err := s.db.Update(func(tx *bolt.Tx) error {
		loge.Debugln("Creating bucket")
		b, err := tx.CreateBucketIfNotExists([]byte(templateName))
		if err != nil {
			loge.WithError(err)
			return mterrors.ErrUnableSaveTemplate()
		}

		createdAt := time.Now().UTC()

		loge.Debugln("Putting kv data")
		value, _ := json.Marshal(&models.Template{
			Data:      templateData,
			CreatedAt: &createdAt,
			Subject:   templateSubject,
		})

		if new {
			if b.Get([]byte(templateVersion)) != nil {
				loge.Errorln("This version of template already exists:", templateName, templateVersion)
				return mterrors.ErrTemplateAlreadyExists()
			}
		}

		if err := b.Put([]byte(templateVersion), value); err != nil {
			loge.WithError(err).Errorln("Put kv data failed")
			return mterrors.ErrUnableSaveTemplate()
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// DeleteTemplate
// Deletes specific template version from db
func (s *boltTemplateStorage) DeleteTemplate(templateName, templateVersion string) error {
	loge := s.log.WithFields(logrus.Fields{
		"name":    templateName,
		"version": templateVersion,
	})
	loge.Infoln("Trying to delete template")

	err := s.db.Update(func(tx *bolt.Tx) error {
		loge.Debugln("Getting bucket")
		b := tx.Bucket([]byte(templateName))
		if b == nil {
			loge.Infoln("Cannot find bucket")
			return mterrors.ErrTemplateNotExist() //models.ErrTemplateNotExists
		}

		loge.Debugln("Deleting entry")
		// check if entry exists
		if v := b.Get([]byte(templateVersion)); v == nil {
			loge.Infoln("Cannot find version")
			return mterrors.ErrTemplateVersionNotExist()
		}
		if err := b.Delete([]byte(templateVersion)); err != nil {
			loge.WithError(err).Errorln("Version delete failed")
			return mterrors.ErrUnableDeleteTemplate()
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// DeleteTemplate
// Deletes template completely from db
func (s *boltTemplateStorage) DeleteTemplates(templateName string) error {
	loge := s.log.WithField("name", templateName)
	loge.Infoln("Trying to delete all versions of template")

	err := s.db.Update(func(tx *bolt.Tx) error {
		loge.Debugln("Deleting bucket")
		if err := tx.DeleteBucket([]byte(templateName)); err == bolt.ErrBucketNotFound {
			loge.WithError(err).Errorf("Bucket not found")
			return mterrors.ErrTemplateNotExist()
		} else if err != nil {
			loge.WithError(err).Errorln("Bucket delete failed")
			return mterrors.ErrUnableDeleteTemplate()
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
