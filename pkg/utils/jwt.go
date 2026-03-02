package utils

import (
	"errors"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UID string `json:"uid"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, key string, expires ...time.Duration) (string, error) {
	exp := 24 * time.Hour
	if len(expires) > 0 {
		exp = expires[0]
	}

	claims := JWTClaims{
		UID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(key))

	if err != nil {
		return "", errors.New("Failed to generate token")
	}

	return tokenString, nil
}

func GenerateAccessToken(UserID string, expires ...time.Duration) (string, error) {
	return GenerateToken(UserID, config.GetEnv().AccessKey, expires...)
}

func VerifyToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, errors.New("Invalid or expired token")
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Invalid token claims")
}
