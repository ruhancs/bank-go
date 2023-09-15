package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("secret key size must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey: secretKey},nil 
}

func(jwtMake *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload,err := NewPaylload(username,duration)
	if err != nil {
		return "",err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(jwtMake.secretKey))
}

func(jwtMake *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(jwtMake.secretKey),nil
	}
	jwtToken,err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr,ok := err.(*jwt.ValidationError)
		//verificar se o erro do token e de expiracao
		if ok && errors.Is(verr.Inner,ErrorExpiredToken) {
			return nil, ErrorExpiredToken
		}
		return nil, ErrInvalidToken
	}

	//converter jwtToken.Claims em payload
	payload,ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}