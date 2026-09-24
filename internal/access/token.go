package access

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

type tokenClaims struct {
	Subject string    `json:"sub"`
	Expiry  time.Time `json:"exp"`
	Purpose string    `json:"purpose"`
}

type TokenSigner struct {
	secret []byte
	clock  func() time.Time
}

func NewTokenSigner(secret []byte, clock func() time.Time) *TokenSigner {
	if clock == nil {
		clock = time.Now
	}
	return &TokenSigner{secret: secret, clock: clock}
}

func (s *TokenSigner) Sign(claims tokenClaims) (string, error) {
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payloadB64 := base64.RawURLEncoding.EncodeToString(payload)
	signature := s.sign(payloadB64)
	return payloadB64 + "." + signature, nil
}

func (s *TokenSigner) Verify(token string) (tokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return tokenClaims{}, ErrInvalidToken
	}
	payloadB64, signature := parts[0], parts[1]
	if !hmac.Equal([]byte(signature), []byte(s.sign(payloadB64))) {
		return tokenClaims{}, ErrInvalidToken
	}
	payload, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return tokenClaims{}, ErrInvalidToken
	}
	var claims tokenClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return tokenClaims{}, ErrInvalidToken
	}
	if !claims.Expiry.After(s.clock().UTC()) {
		return tokenClaims{}, ErrExpiredToken
	}
	return claims, nil
}

func (s *TokenSigner) sign(payload string) string {
	h := hmac.New(sha256.New, s.secret)
	_, _ = h.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
