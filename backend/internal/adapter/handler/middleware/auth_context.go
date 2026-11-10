package middleware

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/xsolla-school/office-explorer/backend/internal/adapter/handler/dto"
	"github.com/xsolla-school/office-explorer/backend/internal/pkg/lmsauth"
)

type contextKey string

const userIDContextKey contextKey = "student_user_id"

const DefaultStudentUUID = "00000000-0000-0000-0000-000000000001"

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// StudentIdentity decides whose progress a request reads and changes:
//
//  1. Authorization: Bearer <LMS access token> -> a fixed UUID for that LMS user
//     (lmsauth.ProgressUserID), when lmsSecret is set
//     (LMS_JWT_SECRET = the LMS backend's JWT_ACCESS_SECRET). A bad or expired token gets 401,
//     so the frontend can refresh it or ask the student to sign in again.
//  2. X-User-ID: a per-browser guest UUID.
//  3. Otherwise the shared default guest.
//
// Without lmsSecret, Authorization is ignored and everyone is a guest, as before.
func StudentIdentity(next http.Handler, lmsSecret []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := ""

		if bearer, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer "); ok && len(lmsSecret) > 0 {
			claims, err := lmsauth.Verify(strings.TrimSpace(bearer), lmsSecret, time.Now())
			if err != nil {
				dto.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "LMS sign-in is invalid or expired; sign in again")
				return
			}
			userID = lmsauth.ProgressUserID(claims.UserID)
		}

		if userID == "" {
			userID = r.Header.Get("X-User-ID")
			if userID == "" || !uuidRegex.MatchString(userID) {
				userID = DefaultStudentUUID
			}
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) string {
	if val, ok := ctx.Value(userIDContextKey).(string); ok && val != "" {
		return val
	}
	return DefaultStudentUUID
}
