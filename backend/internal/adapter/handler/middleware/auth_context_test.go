package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/xsolla-school/office-explorer/backend/internal/pkg/lmsauth"
)

var lmsSecret = []byte("test-lms-access-secret")

func lmsToken(t *testing.T, key []byte, sub int64, exp time.Time) string {
	t.Helper()
	h := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	p := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"sub":%d,"roles":["student"],"exp":%d}`, sub, exp.Unix())))
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(h + "." + p))
	return h + "." + p + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func TestStudentIdentity(t *testing.T) {
	const guest = "6f1c2a7e-1b2c-4d3e-8f90-123456789abc"
	valid := lmsToken(t, lmsSecret, 42, time.Now().Add(10*time.Minute))
	expired := lmsToken(t, lmsSecret, 42, time.Now().Add(-time.Minute))
	otherKey := lmsToken(t, []byte("not-the-lms-secret"), 42, time.Now().Add(10*time.Minute))

	tests := []struct {
		name       string
		secret     []byte
		bearer     string
		xUserID    string
		wantStatus int
		wantUser   string
	}{
		{"no headers: shared default guest", lmsSecret, "", "", 200, DefaultStudentUUID},
		{"guest UUID from X-User-ID", lmsSecret, "", guest, 200, guest},
		{"invalid X-User-ID falls back to default guest", lmsSecret, "", "not-a-uuid", 200, DefaultStudentUUID},
		{"valid LMS token wins over X-User-ID", lmsSecret, valid, guest, 200, lmsauth.ProgressUserID(42)},
		{"expired LMS token: 401", lmsSecret, expired, guest, 401, ""},
		{"token signed with another key: 401", lmsSecret, otherKey, guest, 401, ""},
		{"garbage token: 401", lmsSecret, "abc", guest, 401, ""},
		{"LMS sign-in off: token ignored, guest used", nil, valid, guest, 200, guest},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gotUser string
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotUser = GetUserID(r.Context())
				w.WriteHeader(http.StatusOK)
			})
			req := httptest.NewRequest(http.MethodGet, "/api/v1/explorer/init", nil)
			if tc.bearer != "" {
				req.Header.Set("Authorization", "Bearer "+tc.bearer)
			}
			if tc.xUserID != "" {
				req.Header.Set("X-User-ID", tc.xUserID)
			}
			rec := httptest.NewRecorder()
			StudentIdentity(next, tc.secret).ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body %s)", rec.Code, tc.wantStatus, rec.Body.String())
			}
			if gotUser != tc.wantUser {
				t.Fatalf("user = %q, want %q", gotUser, tc.wantUser)
			}
		})
	}
}
