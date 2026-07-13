package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// GenerateToken creates a JWT access token signed with HMAC-SHA256.
func GenerateToken(userId string, secret []byte) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadData := map[string]interface{}{
		"userId": userId,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}
	payloadJSON, err := json.Marshal(payloadData)
	if err != nil {
		return "", err
	}
	payload := base64.RawURLEncoding.EncodeToString(payloadJSON)

	h := hmac.New(sha256.New, secret)
	h.Write([]byte(header + "." + payload))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return header + "." + payload + "." + signature, nil
}

// VerifyJWT validates a JWT token and returns the userId if successful.
func VerifyJWT(token string, secret []byte) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", errors.New("invalid token format")
	}

	header, payload, signature := parts[0], parts[1], parts[2]
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(header + "." + payload))
	expectedSignature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return "", errors.New("invalid signature")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return "", err
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return "", err
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return "", errors.New("expiration claim missing")
	}

	if time.Now().Unix() > int64(expFloat) {
		return "", errors.New("token expired")
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		return "", errors.New("userId claim missing")
	}

	return userId, nil
}
