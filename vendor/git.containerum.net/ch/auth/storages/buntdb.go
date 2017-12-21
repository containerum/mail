package storages

import (
	"encoding/json"

	"crypto/md5"
	"encoding/hex"

	"git.containerum.net/ch/auth/token"
	"git.containerum.net/ch/auth/utils"
	"git.containerum.net/ch/grpc-proto-files/auth"
	"git.containerum.net/ch/grpc-proto-files/common"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
	"golang.org/x/net/context"
)

const (
	indexTokens = "tokens"
	indexUsers  = "users"
)

type BuntDBStorageConfig struct {
	File         string
	BuntDBConfig buntdb.Config
	TokenFactory token.IssuerValidator
}

// TokenStorage using BuntDB library
type BuntDBStorage struct {
	db     *buntdb.DB
	logger *logrus.Entry
	BuntDBStorageConfig
}

func NewBuntDBStorage(config BuntDBStorageConfig) (storage *BuntDBStorage, err error) {
	logger := logrus.WithField("component", "BuntDBStorage")
	logger.WithField("config", config).Info("Initializing BuntDBStorage")

	logger.Debugf("Opening file %s", config.File)
	db, err := buntdb.Open(config.File)
	if err != nil {
		return nil, err
	}

	logger.Debugf("Setting database config")
	if err := db.SetConfig(config.BuntDBConfig); err != nil {
		return nil, err
	}

	err = db.Update(func(tx *buntdb.Tx) error {
		logger.Debugf("Create index for tokens")
		if err := tx.CreateIndex(indexTokens, "*", buntdb.IndexJSON("platform"),
			buntdb.IndexJSON("fingerprint"), buntdb.IndexJSON("user_ip")); err != nil {
			return err
		}
		logger.Debugf("Create index for users")
		if err := tx.CreateIndex(indexUsers, "*", buntdb.IndexJSON("user_id.value")); err != nil {
			return err
		}
		return nil
	})
	return &BuntDBStorage{
		db:                  db,
		BuntDBStorageConfig: config,
		logger:              logger,
	}, err
}

type tokenOwnerIdentity struct {
	UserAgent, UserIp, Fingerprint string
}

func (s *BuntDBStorage) forTokensByIdentity(tx *buntdb.Tx,
	identity *tokenOwnerIdentity,
	iterator func(key, value string) bool) error {
	pivot, _ := json.Marshal(auth.StoredToken{
		Platform:    utils.ShortUserAgent(identity.UserAgent),
		UserIp:      identity.UserIp,
		Fingerprint: identity.Fingerprint,
	})
	s.logger.WithField("pivot", pivot).Debugf("Iterating by identity")
	return tx.AscendEqual(indexTokens, string(pivot), iterator)
}

func (s *BuntDBStorage) forTokensByUsers(tx *buntdb.Tx, userId *common.UUID, iterator func(key, value string) bool) error {
	pivot, _ := json.Marshal(auth.StoredToken{
		UserId: userId,
	})
	s.logger.WithField("pivot", pivot).Debugf("Iterating by user")
	return tx.AscendEqual(indexUsers, string(pivot), iterator)
}

func (s *BuntDBStorage) marshalRecord(st *auth.StoredToken) string {
	ret, _ := json.Marshal(st)
	s.logger.WithField("record", st).Debugf("Marshal record")
	return string(ret)
}

func (s *BuntDBStorage) unmarshalRecord(rawRecord string) *auth.StoredToken {
	ret := new(auth.StoredToken)
	json.Unmarshal([]byte(rawRecord), ret)
	s.logger.WithField("rawRecord", rawRecord).Debugf("Unmarshal record")
	return ret
}

func (s *BuntDBStorage) deleteTokenByIdentity(tx *buntdb.Tx, identity *tokenOwnerIdentity) error {
	s.logger.WithField("identity", identity).Debugf("Delete token by identity")

	var keysToDelete []string
	err := s.forTokensByIdentity(tx, identity, func(key, value string) bool {
		keysToDelete = append(keysToDelete, key)
		return true
	})
	if err != nil {
		return err
	}
	for _, v := range keysToDelete {
		if _, err := tx.Delete(v); err != nil {
			return err
		}
	}
	return nil
}

func (s *BuntDBStorage) deleteTokenByUser(tx *buntdb.Tx, userId *common.UUID) error {
	s.logger.WithField("userId", userId).Debugf("Delete token by user")

	var keysToDelete []string
	err := s.forTokensByUsers(tx, userId, func(key, value string) bool {
		keysToDelete = append(keysToDelete, key)
		return true
	})
	if err != nil {
		return err
	}
	for _, v := range keysToDelete {
		if _, err := tx.Delete(v); err != nil {
			return err
		}
	}
	return nil
}

func (s *BuntDBStorage) CreateToken(ctx context.Context, req *auth.CreateTokenRequest) (*auth.CreateTokenResponse, error) {
	logger := s.logger.WithField("request", req)

	logger.Info("Creating token")
	var accessToken, refreshToken *token.IssuedToken
	err := s.db.Update(func(tx *buntdb.Tx) error {
		// remove already exist tokens
		logger.Debug("Remove already exist tokens")
		if err := s.deleteTokenByIdentity(tx, &tokenOwnerIdentity{
			UserAgent:   req.UserAgent,
			UserIp:      req.UserIp,
			Fingerprint: req.Fingerprint,
		}); err != nil {
			return err
		}

		// issue tokens
		var err error
		userIdHash := md5.Sum([]byte(req.UserId.Value))
		logger.WithField("userIDHash", userIdHash).Debug("Issue tokens")
		accessToken, refreshToken, err = s.TokenFactory.IssueTokens(token.ExtensionFields{
			UserIDHash: hex.EncodeToString(userIdHash[:]),
			Role:       req.UserRole.String(),
		})
		if err != nil {
			return err
		}

		// store tokens
		logger.WithField("accessToken", accessToken).
			WithField("refreshToken", refreshToken).
			Debug("Store tokens")
		_, _, err = tx.Set(refreshToken.Id.Value,
			s.marshalRecord(token.RequestToRecord(req, refreshToken)),
			&buntdb.SetOptions{
				Expires: true,
				TTL:     refreshToken.LifeTime,
			})
		return err
	})

	if err != nil {
		return nil, err
	}

	return &auth.CreateTokenResponse{
		AccessToken:  accessToken.Value,
		RefreshToken: refreshToken.Value,
	}, nil
}

func (s *BuntDBStorage) CheckToken(ctx context.Context, req *auth.CheckTokenRequest) (*auth.CheckTokenResponse, error) {
	logger := s.logger.WithField("request", req)

	logger.Infof("Validating token")
	valid, err := s.TokenFactory.ValidateToken(req.AccessToken)
	if err != nil || !valid.Valid || valid.Kind != token.KindAccess { // only access tokens may be checked
		return nil, ErrInvalidToken
	}
	var rec *auth.StoredToken
	logger.Debugf("Find record in storage")
	err = s.db.View(func(tx *buntdb.Tx) error {
		rawRec, err := tx.Get(valid.Id.Value)
		if err != nil {
			return err
		}
		rec = s.unmarshalRecord(rawRec)
		return nil
	})
	if err != nil || rec.UserIp != req.UserIp || rec.Fingerprint != req.FingerPrint {
		return nil, ErrTokenNotOwnedBySender
	}

	return &auth.CheckTokenResponse{
		Access: &auth.ResourcesAccess{
			Namespace: token.DecodeAccessObjects(rec.UserNamespace),
			Volume:    token.DecodeAccessObjects(rec.UserVolume),
		},
		UserId:      rec.UserId,
		UserRole:    rec.UserRole,
		TokenId:     rec.TokenId,
		PartTokenId: rec.PartTokenId,
	}, nil
}

func (s *BuntDBStorage) ExtendToken(ctx context.Context, req *auth.ExtendTokenRequest) (*auth.ExtendTokenResponse, error) {
	logger := s.logger.WithField("request", req)

	logger.Info("Extend token")

	// validate received token
	logger.Debugf("Validate token")
	valid, err := s.TokenFactory.ValidateToken(req.RefreshToken)
	if err != nil || !valid.Valid || valid.Kind != token.KindRefresh { // user must send refresh token
		return nil, ErrInvalidToken
	}

	var accessToken, refreshToken *token.IssuedToken
	err = s.db.Update(func(tx *buntdb.Tx) error {
		// identify token owner
		logger.Debugf("Identify token owner")
		rawRec, err := tx.Get(valid.Id.Value)
		if err != nil {
			return err
		}
		rec := s.unmarshalRecord(rawRec)
		if rec.Fingerprint != req.Fingerprint {
			return ErrTokenNotOwnedBySender
		}

		// remove old tokens
		logger.WithField("record", rec).Debugf("Delete old token")
		if err := s.deleteTokenByIdentity(tx, &tokenOwnerIdentity{
			UserAgent:   rec.UserAgent,
			UserIp:      rec.UserIp,
			Fingerprint: rec.Fingerprint,
		}); err != nil {
			return err
		}

		// issue new tokens
		userIdHash := md5.Sum([]byte(rec.UserId.Value))
		logger.WithField("userIdHash", userIdHash).Debug("Issue new tokens")
		accessToken, refreshToken, err = s.TokenFactory.IssueTokens(token.ExtensionFields{
			UserIDHash: hex.EncodeToString(userIdHash[:]),
			Role:       rec.UserRole.String(),
		})
		if err != nil {
			return err
		}
		refreshTokenRecord := *rec
		refreshTokenRecord.TokenId = refreshToken.Id

		// store new tokens
		logger.WithField("record", refreshTokenRecord).Debug("Store new tokens")
		_, _, err = tx.Set(refreshToken.Id.Value,
			s.marshalRecord(&refreshTokenRecord),
			&buntdb.SetOptions{
				Expires: true,
				TTL:     refreshToken.LifeTime,
			})
		return err
	})

	if err != nil {
		return nil, err
	}

	return &auth.ExtendTokenResponse{
		AccessToken:  accessToken.Value,
		RefreshToken: refreshToken.Value,
	}, nil
}

func (*BuntDBStorage) UpdateAccess(context.Context, *auth.UpdateAccessRequest) (*empty.Empty, error) {
	panic("implement me")
}

func (s *BuntDBStorage) GetUserTokens(ctx context.Context, req *auth.GetUserTokensRequest) (*auth.GetUserTokensResponse, error) {
	logger := s.logger.WithField("request", req)

	logger.Infof("Get user tokens")
	resp := new(auth.GetUserTokensResponse)
	err := s.db.View(func(tx *buntdb.Tx) error {
		return s.forTokensByUsers(tx, req.UserId, func(key, value string) bool {
			rec := s.unmarshalRecord(value)
			resp.Tokens = append(resp.Tokens, &auth.StoredTokenForUser{
				TokenId:   rec.TokenId,
				UserAgent: rec.UserAgent,
				Ip:        rec.UserIp,
				// CreatedAt is not stored in db
			})
			return true
		})
	})
	return resp, err
}

func (s *BuntDBStorage) DeleteToken(ctx context.Context, req *auth.DeleteTokenRequest) (*empty.Empty, error) {
	logger := s.logger.WithField("request", req)

	logger.Infof("Delete token")
	return new(empty.Empty), s.db.Update(func(tx *buntdb.Tx) error {
		value, err := tx.Delete(req.TokenId.Value)
		if err != nil {
			return err
		}
		rec := s.unmarshalRecord(value)
		if !utils.UUIDEquals(rec.UserId, req.UserId) {
			err = ErrTokenNotOwnedBySender
		}
		return err
	})
}

func (s *BuntDBStorage) DeleteUserTokens(ctx context.Context, req *auth.DeleteUserTokensRequest) (*empty.Empty, error) {
	logger := s.logger.WithField("request", req)

	logger.Infof("Delete user tokens")
	return new(empty.Empty), s.db.Update(func(tx *buntdb.Tx) error {
		return s.deleteTokenByUser(tx, req.UserId)
	})
}

// Implement Closer interface
func (s *BuntDBStorage) Close() error {
	s.logger.Info("Closing database")
	return s.db.Close()
}
