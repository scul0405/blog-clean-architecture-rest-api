package paseto

import (
	"errors"
	"github.com/o1egl/paseto"
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	Email     string    `json:"email"`
	ID        string    `json:"id"`
	IssueAt   time.Time `json:"issue_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func GeneratePASETOToken(user *models.User, config *config.Config) (string, error) {
	payload := &Payload{
		Email:     user.Email,
		ID:        user.UserID.String(),
		IssueAt:   time.Now(),
		ExpiredAt: time.Now().Add(time.Minute * 60),
	}

	token, err := paseto.NewV2().Encrypt([]byte(config.Server.SymmetricKey), payload, nil)
	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyPASETOToken(token string, config *config.Config) (*Payload, error) {
	payload := &Payload{}

	err := paseto.NewV2().Decrypt(token, []byte(config.Server.SymmetricKey), payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if err = payload.IsValid(); err != nil {
		return nil, err
	}

	return payload, nil
}

func (p *Payload) IsValid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}

	return nil
}
