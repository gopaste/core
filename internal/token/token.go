package token

import (
	"fmt"
	"time"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type Maker interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(user *entity.User, expiry time.Duration) (string, *entity.Payload, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*entity.Payload, error)
}

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *PasetoMaker) CreateToken(user *entity.User, duration time.Duration) (string, *entity.Payload, error) {
	payload, err := entity.NewPayload(user.Name, user.ID, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token, payload, err
}

// VerifyToken checks if the token is valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*entity.Payload, error) {
	payload := &entity.Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, entity.TokenInvalidError
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
