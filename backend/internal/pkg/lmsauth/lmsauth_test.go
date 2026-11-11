package lmsauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
	"testing"
	"time"
)

var secret = []byte("test-lms-access-secret")

// sign builds a token the same way the LMS backend does (internal/auth/jwt.go there).
func sign(t *testing.T, key []byte, header, payload string) string {
	t.Helper()
	h := base64.RawURLEncoding.EncodeToString([]byte(header))
	p := base64.RawURLEncoding.EncodeToString([]byte(payload))
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(h + "." + p))
	return h + "." + p + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// editPayload swaps a token's payload and keeps the old signature.
func editPayload(token, payload string) string {
	first := strings.Index(token, ".")
	last := strings.LastIndex(token, ".")
	return token[:first+1] + base64.RawURLEncoding.EncodeToString([]byte(payload)) + token[last:]
}

func TestVerify(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	const hs256 = `{"alg":"HS256","typ":"JWT"}`
	valid := sign(t, secret, hs256, `{"sub":5,"roles":["student"],"exp":1800000600}`)

	c, err := Verify(valid, secret, now)
	if err != nil {
		t.Fatalf("valid token rejected: %v", err)
	}
	if c.UserID != 5 || len(c.Roles) != 1 || c.Roles[0] != "student" {
		t.Fatalf("unexpected claims: %+v", c)
	}

	tests := []struct {
		name  string
		token string
		key   []byte
		want  error
	}{
		{"wrong secret", valid, []byte("other"), ErrSignature},
		{"expired", sign(t, secret, hs256, `{"sub":5,"roles":["student"],"exp":1799999999}`), secret, ErrExpired},
		{"payload edited after signing (sub 5 -> 1)", editPayload(valid, `{"sub":1,"roles":["admin"],"exp":1800000600}`), secret, ErrSignature},
		{"alg none", sign(t, secret, `{"alg":"none"}`, `{"sub":5,"exp":1800000600}`), secret, ErrMalformed},
		{"no user id", sign(t, secret, hs256, `{"roles":["admin"],"exp":1800000600}`), secret, ErrMalformed},
		{"two parts only", "abc.def", secret, ErrMalformed},
		{"server has no secret", valid, nil, ErrMalformed},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Verify(tc.token, tc.key, now); !errors.Is(err, tc.want) {
				t.Fatalf("got %v, want %v", err, tc.want)
			}
		})
	}
}

func TestProgressUserID(t *testing.T) {
	// Reference values from Python: uuid.uuid5(uuid.UUID("8d3f1b2a-6c4e-4f7a-9b1d-2e5c7a9f0b13"), "lms:42")
	if got, want := ProgressUserID(42), "4088c92a-cbbf-5859-aab6-061359277f70"; got != want {
		t.Fatalf("ProgressUserID(42) = %s, want %s", got, want)
	}
	if got, want := ProgressUserID(4), "90a45148-b623-50c4-b501-a97b370aaf37"; got != want {
		t.Fatalf("ProgressUserID(4) = %s, want %s", got, want)
	}
	if ProgressUserID(42) != ProgressUserID(42) {
		t.Fatal("same LMS user must always get the same UUID")
	}
	if ProgressUserID(42) == ProgressUserID(43) {
		t.Fatal("different LMS users must get different UUIDs")
	}
}
