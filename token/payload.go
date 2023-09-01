package token

import (
	"errors"
	"time"
)

type Payload struct {
	Username  string
	CreateAt  time.Time
	ExpiredAt time.Time
}

var ErrExpiredToken = errors.New("token has expired")
var ErrInvalidToken = errors.New("token is invalid")

func NewPayload(username string, duration time.Duration) *Payload {
	createAt := time.Now()

	return &Payload{
		Username:  username,
		CreateAt:  createAt,
		ExpiredAt: createAt.Add(duration),
	}
}

func (payload *Payload) Valid() error {
	tokenIsExpired := time.Now().After(payload.ExpiredAt)
	if tokenIsExpired {
		return ErrExpiredToken
	}

	return nil
}
