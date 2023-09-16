package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	passeto *paseto.V2
	symmetricKey []byte
}

func NewPasetMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: key must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		passeto: paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker,nil
}

func(maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload,err := NewPaylload(username,duration)
	if err != nil {
		return "",err
	}

	return maker.passeto.Encrypt(maker.symmetricKey,payload, nil)
}

func(maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.passeto.Decrypt(token,maker.symmetricKey,payload,nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}