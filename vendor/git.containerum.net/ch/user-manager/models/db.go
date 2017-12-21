package models

import (
	"errors"
	"time"

	"fmt"

	chutils "git.containerum.net/ch/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	migdrv "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"github.com/sirupsen/logrus"
)

type DB struct {
	conn *sqlx.DB // do not use it in select/exec operations
	log  *logrus.Entry
	qLog *chutils.SQLXQueryLogger
	eLog *chutils.SQLXExecLogger
}

var (
	ErrTransactionBegin    = errors.New("transaction begin error")
	ErrTransactionRollback = errors.New("transaction rollback error")
	ErrTransactionCommit   = errors.New("transaction commit error")
)

func DBConnect(pgConnStr string) (*DB, error) {
	log := logrus.WithField("component", "db")
	log.Infoln("Connecting to ", pgConnStr)
	conn, err := sqlx.Open("postgres", pgConnStr)
	if err != nil {
		log.WithError(err).Errorln("Postgres connection failed")
		return nil, err
	}

	ret := &DB{
		conn: conn,
		log:  log,
		qLog: chutils.NewSQLXQueryLogger(conn, log),
		eLog: chutils.NewSQLXExecLogger(conn, log),
	}

	m, err := ret.migrateUp("migrations")
	if err != nil {
		return nil, err
	}
	version, _, _ := m.Version()
	log.WithField("version", version).Infoln("Migrate up")

	return ret, nil
}

func (db *DB) migrateUp(path string) (*migrate.Migrate, error) {
	db.log.Infof("Running migrations")
	instance, err := migdrv.WithInstance(db.conn.DB, &migdrv.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+path, "clickhouse", instance)
	if err != nil {
		return nil, err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}
	return m, nil
}

func (db *DB) Transactional(f func(tx *DB) error) (err error) {
	start := time.Now().Format(time.ANSIC)
	e := db.log.WithField("transaction_at", start)
	e.Debugln("Begin transaction")
	tx, txErr := db.conn.Beginx()
	if txErr != nil {
		e.WithError(txErr).Errorln("Begin transaction error")
		return ErrTransactionBegin
	}

	arg := &DB{
		conn: db.conn,
		log:  e,
		eLog: chutils.NewSQLXExecLogger(tx, e),
		qLog: chutils.NewSQLXQueryLogger(tx, e),
	}

	// needed for recovering panics in transactions.
	defer func(dberr error) {
		// if panic recovered, try to rollback transaction
		if panicErr := recover(); panicErr != nil {
			dberr = fmt.Errorf("panic in transaction: %v", panicErr)
		}

		if dberr != nil {
			e.WithError(dberr).Debugln("Rollback transaction")
			if rerr := tx.Rollback(); rerr != nil {
				e.WithError(rerr).Errorln("Rollback error")
				err = ErrTransactionRollback
			}
			err = dberr // forward error with panic description
			return
		}

		e.Debugln("Commit transaction")
		if cerr := tx.Commit(); cerr != nil {
			e.WithError(cerr).Errorln("Commit error")
			err = ErrTransactionCommit
		}
	}(f(arg))

	return
}

func (db *DB) Close() error {
	return db.conn.Close()
}
