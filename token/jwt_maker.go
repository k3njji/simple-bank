package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)


const minSizeSK=32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error){
	if len(secretKey) < minSizeSK {
		return nil, fmt.Errorf("invalid key size, must be %d characters", minSizeSK)
	}

	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker)CreateToken(username string, duration time.Duration) (string, error){
	payload, err := NewPayload(username, duration)
	if err!=nil{
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
    keyFunc := func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }

        return []byte(maker.secretKey), nil
    }

    jwtToken, err := jwt.ParseWithClaims(
        token,
        &Payload{},
        keyFunc,
    )

    if err != nil {
        return nil, err
    }

    payload, ok := jwtToken.Claims.(*Payload)
    if !ok || !jwtToken.Valid {
        return nil, ErrInvalidToken
    }

    return payload, nil
}