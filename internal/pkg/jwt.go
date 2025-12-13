package pkg

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWT struct {
	secretKey       string
	accessTokenExp  uint
	refreshTokenExp uint
}

func InitJWT(secretKey string, accessTokenExp, refreshTokenExp uint) *JWT {
	return &JWT{
		secretKey,
		accessTokenExp,
		refreshTokenExp,
	}
}

func (j JWT) GenerateAccessToken(userID uint, email string, username string) (string, error) {
	claims := jwt.MapClaims{
		"id":       userID,
		"email":    email,
		"username": username,
		"exp":      time.Now().Add(time.Duration(j.accessTokenExp) * time.Minute).Unix(),
		"iat":      time.Now().Unix(),
		"type":     "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secretKey))
}

func (j JWT) GenerateRefreshToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"id":   userID,
		"exp":  time.Now().Add(time.Duration(j.refreshTokenExp) * 24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
		"type": "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secretKey))
}

func (j JWT) VerifyToken(tokenString string, expectedType string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	if claims["type"] != expectedType {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}
