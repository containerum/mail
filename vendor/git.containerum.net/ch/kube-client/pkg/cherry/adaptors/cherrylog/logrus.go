package cherrylog

import (
	"git.containerum.net/ch/kube-client/pkg/cherry"
	"github.com/sirupsen/logrus"
)

// LogrusAdapter -- log origin and returning errors through logrus
type LogrusAdapter struct {
	*logrus.Entry
}

func (a *LogrusAdapter) Log(origin error, returning *cherry.Err) {
	a.WithError(origin).Errorf("returning %+v", returning)
}

// NewLogrusAdpater -- for more convenient usage
func NewLogrusAdapter(e *logrus.Entry) *LogrusAdapter {
	return &LogrusAdapter{Entry: e}
}
