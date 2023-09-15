package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrorExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type Payload struct {
	ID uuid.UUID `json:"id"`
	Username string `json:"username"`
	IssueAt time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPaylload(username string, duration time.Duration) (*Payload,error) {
	id,err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID: id,
		Username: username,
		IssueAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload,err
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrorExpiredToken
	}
	return nil
}