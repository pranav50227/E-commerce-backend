package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"shared/auth"
	"shared/constants"
)

func TestVerifyJWT(t *testing.T) {
	jwtSecret := []byte(constants.DefaultJWTSecret)

	// Helper to generate a test token
	generateTestToken := func(userId string, expired bool) string {
		header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
		
		var exp int64
		if expired {
			exp = time.Now().Add(-1 * time.Hour).Unix()
		} else {
			exp = time.Now().Add(1 * time.Hour).Unix()
		}

		payloadData := map[string]interface{}{
			"userId": userId,
			"exp":    exp,
		}
		payloadJSON, _ := json.Marshal(payloadData)
		payload := base64.RawURLEncoding.EncodeToString(payloadJSON)

		h := hmac.New(sha256.New, jwtSecret)
		h.Write([]byte(header + "." + payload))
		signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

		return header + "." + payload + "." + signature
	}

	// Test Valid Token
	validToken := generateTestToken("user-123", false)
	userId, err := auth.VerifyJWT(validToken, jwtSecret)
	if err != nil {
		t.Fatalf("expected valid token to pass, got error: %v", err)
	}
	if userId != "user-123" {
		t.Errorf("expected userId user-123, got %s", userId)
	}

	// Test Expired Token
	expiredToken := generateTestToken("user-123", true)
	_, err = auth.VerifyJWT(expiredToken, jwtSecret)
	if err == nil || err.Error() != "token expired" {
		t.Errorf("expected token expired error, got: %v", err)
	}

	// Test Invalid Signature
	badToken := validToken + "corrupted"
	_, err = auth.VerifyJWT(badToken, jwtSecret)
	if err == nil {
		t.Error("expected error for invalid signature, got nil")
	}
}
