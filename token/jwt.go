package token

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func (maker JWTMaker) generateToken(username string, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, role, duration)
	if err != nil {
		return "", nil, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

func (maker JWTMaker) verifyToken(token string) (*Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(maker.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, fmt.Errorf("token is invalid")
	}

	return payload, nil
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("secret key must be at least %d char", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}
