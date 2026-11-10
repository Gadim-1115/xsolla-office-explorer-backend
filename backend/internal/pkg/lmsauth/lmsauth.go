// Package lmsauth checks access tokens issued by the Xsolla School LMS backend,
// so a student signed in to the LMS is the same person in the office tour.
//
// LMS tokens are JWTs signed with HMAC-SHA256 using the LMS's JWT_ACCESS_SECRET.
// Their payload is {"sub": <user id>, "roles": [...], "exp": <unix seconds>}.
package lmsauth

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// progressNamespace is the fixed namespace for ProgressUserID. Never change it:
// every LMS student's stored progress is keyed by UUIDs derived from it.
var progressNamespace = mustParseUUID("8d3f1b2a-6c4e-4f7a-9b1d-2e5c7a9f0b13")

// ProgressUserID is the UUID under which an LMS student's office-tour progress is stored:
// a version-5 UUID (SHA-1 of the namespace and "lms:<id>"). The same LMS user always gets
// the same UUID, it cannot equal a guest's random version-4 UUID, and the existing
// user_progress.user_id UUID column needs no change.
func ProgressUserID(lmsUserID int64) string {
	h := sha1.New()
	h.Write(progressNamespace[:])
	h.Write([]byte("lms:" + strconv.FormatInt(lmsUserID, 10)))
	u := h.Sum(nil)[:16]
	u[6] = (u[6] & 0x0f) | 0x50 // version 5
	u[8] = (u[8] & 0x3f) | 0x80 // RFC 4122 variant
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:16])
}

func mustParseUUID(s string) [16]byte {
	var out [16]byte
	b, err := hex.DecodeString(strings.ReplaceAll(s, "-", ""))
	if err != nil || len(b) != 16 {
		panic("lmsauth: bad namespace UUID " + s)
	}
	copy(out[:], b)
	return out
}

var (
	ErrMalformed = errors.New("lms token is malformed")
	ErrSignature = errors.New("lms token signature does not match")
	ErrExpired   = errors.New("lms token has expired")
)

// Claims is the payload of an LMS access token.
type Claims struct {
	UserID int64    `json:"sub"`
	Roles  []string `json:"roles"`
	Exp    int64    `json:"exp"`
}

// Verify checks an LMS access token's algorithm, signature and expiry.
func Verify(token string, secret []byte, now time.Time) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 || len(secret) == 0 {
		return nil, ErrMalformed
	}

	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, ErrMalformed
	}
	var header struct {
		Alg string `json:"alg"`
	}
	if json.Unmarshal(headerJSON, &header) != nil || header.Alg != "HS256" {
		return nil, ErrMalformed
	}

	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, ErrMalformed
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(parts[0] + "." + parts[1]))
	if !hmac.Equal(sig, mac.Sum(nil)) {
		return nil, ErrSignature
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrMalformed
	}
	var c Claims
	if json.Unmarshal(payload, &c) != nil || c.UserID <= 0 {
		return nil, ErrMalformed
	}
	if now.Unix() >= c.Exp {
		return nil, ErrExpired
	}
	return &c, nil
}
