package token

import (
	"time"

	"errors"

	"git.containerum.net/ch/auth/utils"
	"git.containerum.net/ch/grpc-proto-files/common"
)

var _ IssuerValidator = &MockIssuerValidator{}

type mockTokenRecord struct {
	IssuedAt time.Time
}

type MockIssuerValidator struct {
	returnedLifeTime time.Duration
	issuedTokens     map[string]mockTokenRecord
}

func NewMockIssuerValidator(returnedLifeTime time.Duration) *MockIssuerValidator {
	return &MockIssuerValidator{
		returnedLifeTime: returnedLifeTime,
		issuedTokens:     make(map[string]mockTokenRecord),
	}
}

func (m *MockIssuerValidator) IssueTokens(extensionFields ExtensionFields) (accessToken, refreshToken *IssuedToken, err error) {
	tokenId := utils.NewUUID()
	accessToken = &IssuedToken{
		Value:    "a" + tokenId.Value,
		LifeTime: m.returnedLifeTime,
		Id:       tokenId,
	}
	m.issuedTokens[tokenId.Value] = mockTokenRecord{
		IssuedAt: time.Now(),
	}
	refreshToken = &IssuedToken{
		Value:    "r" + tokenId.Value,
		LifeTime: m.returnedLifeTime,
		Id:       tokenId,
	}
	return
}

func (m *MockIssuerValidator) ValidateToken(token string) (result *ValidationResult, err error) {
	rec, present := m.issuedTokens[token[1:]]
	var kind Kind
	switch token[0] {
	case 'a':
		kind = KindAccess
	case 'r':
		kind = KindRefresh
	default:
		return nil, errors.New("invalid token received")
	}
	return &ValidationResult{
		Valid: present && time.Now().Before(rec.IssuedAt.Add(m.returnedLifeTime)),
		Kind:  kind,
		Id:    &common.UUID{Value: token[1:]},
	}, nil
}
