package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenUtils interface {
	GenerateAccessToken(userID uint64) (string, error)
	GenerateRefreshToken(userID uint64) (string, error)
	ParseRefreshToken(tokenString string) (*Claims, error)
}

type tokenUtils struct {
	AccessSecret  []byte
	RefreshSecret []byte
}

func NewTokenUtils(accessSecret []byte, refreshSecret []byte) TokenUtils {
	return &tokenUtils{
		AccessSecret:  accessSecret,
		RefreshSecret: refreshSecret,
	}
}

func (tu *tokenUtils) GenerateAccessToken(userID uint64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tu.AccessSecret)
}

func (tu *tokenUtils) GenerateRefreshToken(userID uint64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tu.RefreshSecret)
}

func (tu *tokenUtils) ParseRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return tu.RefreshSecret, nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
