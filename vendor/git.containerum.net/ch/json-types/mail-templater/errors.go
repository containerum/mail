package mail

import "git.containerum.net/ch/json-types/errors"

var (
	ErrTemplateNotExists        = &errors.NotFoundError{errors.New("Specified template not exists in storage")}
	ErrTemplateVersionNotExists = &errors.NotFoundError{errors.New("Specified version of template not exists in storage")}
	ErrMessageNotExists         = &errors.NotFoundError{errors.New("Specified message not exists in storage")}
)

var (
	ErrStorageOpenFailed       = &errors.InternalError{errors.New("Failed to open storage")}
	ErrCreateBucketFailed      = &errors.InternalError{errors.New("Create bucket failed")}
	ErrTemplatePutFailed       = &errors.InternalError{errors.New("Unable to save template")}
	ErrTemplateGetFailed       = &errors.InternalError{errors.New("Unable to get template")}
	ErrTemplateGetLatestFailed = &errors.InternalError{errors.New("Unable to get latest version of template")}
	ErrTemplateGetAllFailed    = &errors.InternalError{errors.New("Unable to get all versions of template")}
	ErrTemplateDeleteFailed    = &errors.InternalError{errors.New("Unable to delete template")}
	ErrTemplateDeleteAllFailed = &errors.InternalError{errors.New("Unable to delete all versions of template")}
	ErrTemplatesGetFailed      = &errors.InternalError{errors.New("Unable to get templates list")}
	ErrMessagePutFailed        = &errors.InternalError{errors.New("Unable to save message")}
	ErrMessageGetFailed        = &errors.InternalError{errors.New("Unable to get message")}
	ErrMessagesGetFailed       = &errors.InternalError{errors.New("Unable to get messages list")}
)
