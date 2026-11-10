package middleware

import (
	"errors"
	"net/http"

	"github.com/xsolla-school/office-explorer/backend/internal/adapter/handler/dto"
	"github.com/xsolla-school/office-explorer/backend/internal/domain"
)

func HandleDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrZoneNotFound):
		dto.WriteError(w, http.StatusNotFound, "ZONE_NOT_FOUND", err.Error())
	case errors.Is(err, domain.ErrQuestNotFound):
		dto.WriteError(w, http.StatusNotFound, "QUEST_NOT_FOUND", err.Error())
	case errors.Is(err, domain.ErrHotspotNotFound):
		dto.WriteError(w, http.StatusNotFound, "HOTSPOT_NOT_FOUND", err.Error())
	case errors.Is(err, domain.ErrZoneLocked):
		dto.WriteError(w, http.StatusForbidden, "ZONE_LOCKED", err.Error())
	case errors.Is(err, domain.ErrInvalidSequence):
		dto.WriteError(w, http.StatusUnprocessableEntity, "INVALID_SEQUENCE", err.Error())
	case errors.Is(err, domain.ErrNotAnEasterEgg):
		dto.WriteError(w, http.StatusBadRequest, "NOT_AN_EASTER_EGG", err.Error())
	case errors.Is(err, domain.ErrInvalidInput):
		dto.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		dto.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
	default:
		dto.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected internal server error occurred")
	}
}

func CORS(allowedOrigin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-User-ID")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
