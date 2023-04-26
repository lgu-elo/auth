//go:generate mockgen -source jwt.go -destination ./../../helper/kwt.go -package kwtmock
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type (
	claim struct {
		jwt.StandardClaims
		SecData SecData
	}

	SecData struct {
		ID       string
		Username string
	}
)

var (
	securityHash      = []byte("OAgEnqEGIzIpEbmzfWH0ZhFKUS1BbYEEexgSvhAPVaZR6DZU")
	ErrTokenIsExpired = errors.New("token expired")
)

const Expiration = 24 * time.Hour * 7

func GenJWT(id, username string) (string, error) {
	claims := &claim{
		SecData: SecData{
			ID:       id,
			Username: username,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(Expiration).Unix(),
			Issuer:    "",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(securityHash)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ExtractData(signed string) (*SecData, error) {
	token, err := jwt.ParseWithClaims(signed, &claim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(securityHash), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*claim)
	if !ok {
		return nil, errors.New("cannot check jwt")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, ErrTokenIsExpired
	}

	return &claims.SecData, nil
}

func DropJWT() (string, error) {
	claims := &claim{
		SecData: SecData{
			ID:       "",
			Username: "",
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Unix(),
			Issuer:    "",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(securityHash)
	if err != nil {
		return "", err
	}

	return signed, nil
}
