package middleware

import (
	"context"
	"net/http"
	"regexp"
)

type contextKey string

const userIDContextKey contextKey = "student_user_id"

const DefaultStudentUUID = "00000000-0000-0000-0000-000000000001"

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func StudentIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" || !uuidRegex.MatchString(userID) {
			userID = DefaultStudentUUID
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
