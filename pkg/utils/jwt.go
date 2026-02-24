package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type TokenClaims struct {
	Subject string `json:"sub"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	Exp     int64  `json:"exp"`
	Iat     int64  `json:"iat"`
}

func GenerateAccessToken(secret string, claims TokenClaims) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("marshal token header: %w", err)
	}

	claims.Iat = time.Now().Unix()
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("marshal token claims: %w", err)
	}

	encode := base64.RawURLEncoding.EncodeToString
	headerSegment := encode(headerJSON)
	payloadSegment := encode(payloadJSON)
	unsigned := headerSegment + "." + payloadSegment

	mac := hmac.New(sha256.New, []byte(secret))
	if _, err := mac.Write([]byte(unsigned)); err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	signature := encode(mac.Sum(nil))
	return unsigned + "." + signature, nil
}
