package mail

import "git.containerum.net/ch/json-types/errors"

var (
	ErrMessageNotExists = errors.New("message not exists")
	ErrTemplateNotExists = errors.New("specified template not exists in storage")
	ErrVersionNotExists  = errors.New("specified version not exists in storage")
)