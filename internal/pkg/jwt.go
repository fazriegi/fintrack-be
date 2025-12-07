package pkg

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWT struct {
	secretKey string
	expHour   uint16
}

func InitJWT(secretKey string, expHour uint16) *JWT {
	return &JWT{
		secretKey,
		expHour,
	}
}

func (s JWT) GenerateJWTToken(id uint, email, username string) (string, error) {
	exp := time.Duration(s.expHour) * time.Hour
	claims := jwt.MapClaims{
		"id":       id,
		"email":    email,
		"username": username,
		"exp":      time.Now().Add(exp).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.secretKey))
}

func (s JWT) VerifyJWTTOken(tokenString string) (any, error) {
	errResponse := errors.New("invalid or expired token")
	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errResponse
		}

		return []byte(s.secretKey), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return nil, errResponse
	}

	return token.Claims.(jwt.MapClaims), nil
}
