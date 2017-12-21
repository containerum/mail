package models

import (
	"crypto/sha512"
	"encoding/hex"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type LinkType string

const (
	LinkTypeConfirm   LinkType = "confirm"
	LinkTypePwdChange LinkType = "pwd_change"
	LinkTypeDelete    LinkType = "delete"
)

type Link struct {
	Link      string
	Type      LinkType
	CreatedAt time.Time
	ExpiredAt time.Time
	IsActive  bool
	SentAt    pq.NullTime

	User *User
}

const linkQueryColumnsWithUser = "links.link, links.type, links.created_at, links.expired_at, links.is_active, links.sent_at, " +
	"users.id, users.login, users.password_hash, users.salt, users.role, users.is_active, users.is_deleted, users.is_in_blacklist"
const linkQueryColumns = "link, type, created_at, expired_at, is_active, sent_at"

func (db *DB) CreateLink(linkType LinkType, lifeTime time.Duration, user *User) (*Link, error) {
	now := time.Now().UTC()
	ret := &Link{
		Link:      strings.ToUpper(hex.EncodeToString(sha512.New().Sum([]byte(user.ID)))),
		User:      user,
		Type:      linkType,
		CreatedAt: now,
		ExpiredAt: now.Add(lifeTime),
		IsActive:  true,
	}
	db.log.WithFields(logrus.Fields{
		"user":          user.Login,
		"creation_time": now.Format(time.ANSIC),
	}).Debugln("Create activation link")
	_, err := db.eLog.Exec("INSERT INTO links (link, type, created_at, expired_at, is_active, user_id) VALUES "+
		"($1, $2, $3, $4, $5, $6)", ret.Link, ret.Type, ret.CreatedAt, ret.ExpiredAt, ret.IsActive, ret.User.ID)
	return ret, err
}

func (db *DB) GetLinkForUser(linkType LinkType, user *User) (*Link, error) {
	db.log.Debugln("Get link", linkType, "for", user.Login)
	rows, err := db.qLog.Queryx("SELECT "+linkQueryColumns+" FROM links "+
		"WHERE user_id = $1 AND type = $2 AND is_active AND expired_at > NOW()", user.ID, linkType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, rows.Err()
	}
	link := Link{User: user}
	err = rows.Scan(&link.Link, &link.Type, &link.CreatedAt, &link.ExpiredAt, &link.IsActive, &link.SentAt)

	return &link, err
}

func (db *DB) GetLinkFromString(strLink string) (*Link, error) {
	db.log.Debugln("Get link", strLink)
	rows, err := db.qLog.Queryx("SELECT "+linkQueryColumnsWithUser+" FROM links "+
		"JOIN users ON links.user_id = users.id "+
		"WHERE link = $1 AND links.is_active AND links.expired_at > NOW()", strLink)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, rows.Err()
	}
	defer rows.Close()
	link := Link{User: &User{}}
	err = rows.Scan(&link.Link, &link.Type, &link.CreatedAt, &link.ExpiredAt, &link.IsActive, &link.SentAt,
		&link.User.ID, &link.User.Login, &link.User.PasswordHash, &link.User.Salt, &link.User.Role,
		&link.User.IsActive, &link.User.IsDeleted, &link.User.IsInBlacklist)

	return &link, err
}

func (db *DB) UpdateLink(link *Link) error {
	db.log.Debugf("Update link %#v", link)
	_, err := db.eLog.Exec("UPDATE links set type = $2, expired_at = $3, is_active = $4, sent_at = $5 "+
		"WHERE link = $1", link.Link, link.Type, link.ExpiredAt, link.IsActive, link.SentAt)
	return err
}

func (db *DB) GetUserLinks(user *User) ([]Link, error) {
	db.log.Debugln("Get links for", user.Login)
	var ret []Link
	rows, err := db.qLog.Queryx("SELECT "+linkQueryColumns+" FROM links "+
		"WHERE user_id = $1 AND is_active AND expired_at > NOW()", user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		link := Link{User: user}
		err := rows.Scan(&link.Link, &link.Type, &link.CreatedAt, &link.ExpiredAt, &link.IsActive, &link.SentAt)
		if err != nil {
			return nil, err
		}
		ret = append(ret, link)
	}

	return ret, rows.Err()
}
