package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/xsolla-school/office-explorer/backend/internal/adapter/handler/dto"
	"github.com/xsolla-school/office-explorer/backend/internal/adapter/handler/middleware"
	"github.com/xsolla-school/office-explorer/backend/internal/adapter/storage/postgres"
	"github.com/xsolla-school/office-explorer/backend/internal/domain"
)

type StreakHandler struct {
	streakRepo *postgres.StreakRepo
}

func NewStreakHandler(streakRepo *postgres.StreakRepo) *StreakHandler {
	return &StreakHandler{streakRepo: streakRepo}
}

type completeStreakRequest struct {
	UserID    string  `json:"user_id"`
	UserName  string  `json:"user_name"`
	UserRole  string  `json:"user_role"`
	DateKey   string  `json:"date_key"`
	AvatarURL *string `json:"avatar_url"`
}

func (h *StreakHandler) GetStreak(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))
	if userID == "" {
		userID = middleware.GetUserID(r.Context())
	}
	rec, err := h.streakRepo.GetUserStreak(r.Context(), userID)
	if err != nil {
		middleware.HandleDomainError(w, err)
		return
	}
	dto.WriteSuccess(w, http.StatusOK, rec)
}

func (h *StreakHandler) CompleteDailyBug(w http.ResponseWriter, r *http.Request) {
	var req completeStreakRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleDomainError(w, domain.ErrInvalidInput)
		return
	}
	req.UserID = strings.TrimSpace(req.UserID)
	req.UserName = strings.TrimSpace(req.UserName)
	if req.UserID == "" || req.UserName == "" {
		middleware.HandleDomainError(w, domain.ErrInvalidInput)
		return
	}

	var cleanAvatar *string
	if req.AvatarURL != nil {
		trimmed := strings.TrimSpace(*req.AvatarURL)
		if trimmed != "" {
			cleanAvatar = &trimmed
		}
	}

	rec, err := h.streakRepo.CompleteDailyBug(
		r.Context(),
		req.UserID,
		req.UserName,
		strings.TrimSpace(req.UserRole),
		strings.TrimSpace(req.DateKey),
		cleanAvatar,
	)
	if err != nil {
		middleware.HandleDomainError(w, err)
		return
	}
	dto.WriteSuccess(w, http.StatusOK, rec)
}

func (h *StreakHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	rows, err := h.streakRepo.ListTopStreaks(r.Context(), 10)
	if err != nil {
		middleware.HandleDomainError(w, err)
		return
	}
	if rows == nil {
		rows = []postgres.StreakRecord{}
	}
	dto.WriteSuccess(w, http.StatusOK, rows)
}

type syncLeaderboardRequest struct {
	PruneMissing bool                   `json:"prune_missing"`
	Users        []postgres.LmsSyncUser `json:"users"`
}

func (h *StreakHandler) SyncLeaderboardUsers(w http.ResponseWriter, r *http.Request) {
	var req syncLeaderboardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleDomainError(w, domain.ErrInvalidInput)
		return
	}

	cleaned := make([]postgres.LmsSyncUser, 0, len(req.Users))
	for _, u := range req.Users {
		uid := strings.TrimSpace(u.UserID)
		if uid == "" || !strings.HasPrefix(uid, "lms-") {
			continue
		}
		var avatar *string
		if u.AvatarURL != nil {
			trimmed := strings.TrimSpace(*u.AvatarURL)
			if trimmed != "" {
				avatar = &trimmed
			}
		}
		cleaned = append(cleaned, postgres.LmsSyncUser{
			UserID:    uid,
			UserName:  strings.TrimSpace(u.UserName),
			UserRole:  strings.TrimSpace(u.UserRole),
			AvatarURL: avatar,
		})
	}

	if err := h.streakRepo.SyncLmsUsers(r.Context(), cleaned, req.PruneMissing); err != nil {
		middleware.HandleDomainError(w, err)
		return
	}
	dto.WriteSuccess(w, http.StatusOK, map[string]bool{"synced": true})
}

func (h *StreakHandler) DeleteStreak(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimSpace(r.PathValue("userId"))
	if userID == "" {
		middleware.HandleDomainError(w, domain.ErrInvalidInput)
		return
	}
	if err := h.streakRepo.DeleteUserStreak(r.Context(), userID); err != nil {
		middleware.HandleDomainError(w, err)
		return
	}
	dto.WriteSuccess(w, http.StatusOK, map[string]bool{"deleted": true})
}
