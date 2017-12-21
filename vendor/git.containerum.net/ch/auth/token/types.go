package token

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"git.containerum.net/ch/auth/utils"
	"git.containerum.net/ch/grpc-proto-files/auth"
	"git.containerum.net/ch/grpc-proto-files/common"
)

type Kind int

const (
	KindAccess Kind = iota
	KindRefresh
)

type ExtensionFields struct {
	UserIDHash  string `json:"userID,omitempty"`
	Role        string `json:"role,omitempty"`
	PartTokenId string `json:"partTokenID,omitempty"`
}

type IssuedToken struct {
	Value    string
	Id       *common.UUID
	LifeTime time.Duration
}

// Issuer is interface for creating access and refresh tokens.
type Issuer interface {
	IssueTokens(extensionFields ExtensionFields) (accessToken, refreshToken *IssuedToken, err error)
}

type ValidationResult struct {
	Valid bool
	Id    *common.UUID
	Kind  Kind
}

// Validator is interface for validating tokens
type Validator interface {
	ValidateToken(token string) (result *ValidationResult, err error)
}

type IssuerValidator interface {
	Issuer
	Validator
}

func EncodeAccessObjects(req []*auth.AccessObject) string {
	ret, _ := json.Marshal(req)
	return base64.StdEncoding.EncodeToString(ret)
}

func DecodeAccessObjects(value string) (ret []*auth.AccessObject) {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return make([]*auth.AccessObject, 0)
	}
	err = json.Unmarshal(decoded, &ret)
	if err != nil {
		return make([]*auth.AccessObject, 0)
	}
	return
}

func RequestToRecord(req *auth.CreateTokenRequest, token *IssuedToken) *auth.StoredToken {
	return &auth.StoredToken{
		TokenId:       token.Id,
		UserAgent:     req.UserAgent,
		Platform:      utils.ShortUserAgent(req.UserAgent),
		Fingerprint:   req.Fingerprint,
		UserId:        req.UserId,
		UserRole:      req.UserRole,
		UserNamespace: EncodeAccessObjects(req.Access.Namespace),
		UserVolume:    EncodeAccessObjects(req.Access.Volume),
		RwAccess:      req.RwAccess,
		UserIp:        req.UserIp,
		PartTokenId:   req.PartTokenId,
	}
}
