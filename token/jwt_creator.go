package token

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTCreator struct {
	secret string
}

func NewJWTCreator(secret string) Creator {
	return &JWTCreator{
		secret: secret,
	}
}

func (creator *JWTCreator) Create(username string, duration time.Duration) (string, error) {
	payload := NewPayload(username, duration)

	jwtToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		payload,
	)

	return jwtToken.SignedString([]byte(creator.secret))
}

func (creator *JWTCreator) Verify(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(creator.secret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(
		token,
		&Payload{},
		keyFunc,
	)

	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(validationErr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}

		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
