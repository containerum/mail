package models

import (
	"time"

	chutils "git.containerum.net/ch/user-manager/utils"
)

type Token struct {
	Token     string
	CreatedAt time.Time
	IsActive  bool
	SessionID string

	User *User
}

const tokenQueryColumnsWithUser = "tokens.token, tokens.created_at, tokens.is_active, tokens.session_id, " +
	"users.id, users.login, users.password_hash, users.salt, users.role, users.is_active, users.is_deleted, users.is_in_blacklist"
const tokenQueryColumns = "token, created_at, is_active, session_id"

func (db *DB) GetTokenObject(token string) (*Token, error) {
	db.log.Debugln("Get token object", token)
	rows, err := db.qLog.Queryx("SELECT "+tokenQueryColumnsWithUser+" FROM tokens "+
		"JOIN users ON tokens.user_id = users.id WHERE tokens.token = $1 AND tokens.is_active", token)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, rows.Err()
	}
	defer rows.Close()
	ret := Token{User: &User{}}
	err = rows.Scan(&ret.Token, &ret.CreatedAt, &ret.IsActive, &ret.SessionID,
		&ret.User.ID, &ret.User.Login, &ret.User.PasswordHash, &ret.User.Salt, &ret.User.Role,
		&ret.User.IsActive, &ret.User.IsDeleted, &ret.User.IsInBlacklist)
	return &ret, err
}

func (db *DB) CreateToken(user *User, sessionID string) (*Token, error) {
	db.log.Debugln("Generate one-time token for", user.Login)
	ret := &Token{
		Token:     chutils.GenSalt(user.ID, user.Login),
		User:      user,
		IsActive:  true,
		SessionID: sessionID,
		CreatedAt: time.Now().UTC(),
	}
	_, err := db.eLog.Exec("INSERT INTO tokens (token, user_id, is_active, session_id, created_at) "+
		"VALUES ($1, $2, $3, $4, $5)", ret.Token, ret.User.ID, ret.IsActive, ret.SessionID, ret.CreatedAt)
	return ret, err
}

func (db *DB) GetTokenBySessionID(sessionID string) (*Token, error) {
	db.log.Debugln("Get token by session id ", sessionID)
	rows, err := db.qLog.Queryx("SELECT "+tokenQueryColumnsWithUser+" FROM tokens "+
		"JOIN users ON tokens.user_id = users.id WHERE tokens.session_id = $1 and tokens.is_active", sessionID)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, rows.Err()
	}
	defer rows.Close()
	ret := Token{User: &User{}}
	err = rows.Scan(&ret.Token, &ret.CreatedAt, &ret.IsActive, &ret.SessionID,
		&ret.User.ID, &ret.User.Login, &ret.User.PasswordHash, &ret.User.Salt, &ret.User.Role,
		&ret.User.IsActive, &ret.User.IsDeleted, &ret.User.IsInBlacklist)

	return &ret, err
}

func (db *DB) DeleteToken(token string) error {
	db.log.Debugln("Remove token", token)
	_, err := db.eLog.Exec("DELETE FROM tokens WHERE token = $1", token)
	return err
}

func (db *DB) UpdateToken(token *Token) error {
	db.log.Debugln("Update token", token.Token)
	_, err := db.eLog.Exec("UPDATE tokens SET is_active = $2, session_id = $3 WHERE token = $1",
		token.Token, token.IsActive, token.SessionID)
	return err
}
