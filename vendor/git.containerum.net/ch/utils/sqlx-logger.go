package utils

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type SQLXQueryLogger struct {
	q sqlx.Queryer
	l *logrus.Entry
}

func NewSQLXQueryLogger(queryer sqlx.Queryer, entry *logrus.Entry) *SQLXQueryLogger {
	return &SQLXQueryLogger{
		q: queryer,
		l: entry,
	}
}

func (q *SQLXQueryLogger) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := q.q.Query(query, args...)
	q.l.WithField("query", query).WithError(err).Debugln(args...)
	return rows, err
}

func (q *SQLXQueryLogger) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	rows, err := q.q.Queryx(query, args...)
	q.l.WithField("query", query).WithError(err).Debugln(args...)
	return rows, err
}

func (q *SQLXQueryLogger) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	row := q.q.QueryRowx(query, args...)
	q.l.WithField("query", query).Debugln(args...)
	return row
}

type SQLXExecLogger struct {
	e sqlx.Execer
	l *logrus.Entry
}

func NewSQLXExecLogger(execer sqlx.Execer, entry *logrus.Entry) *SQLXExecLogger {
	return &SQLXExecLogger{
		e: execer,
		l: entry,
	}
}

func (e *SQLXExecLogger) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := e.e.Exec(query, args...)
	e.l.WithField("query", query).WithError(err).Debugln(args...)
	return result, err
}
