package handler

import (
	"encoding/json"
	"net/http"

	"github.com/xsolla-school/office-explorer/backend/internal/adapter/handler/dto"
	"github.com/xsolla-school/office-explorer/backend/internal/adapter/handler/middleware"
	"github.com/xsolla-school/office-explorer/backend/internal/domain"
	"github.com/xsolla-school/office-explorer/backend/internal/usecase"
)

type QuestHandler struct {
	questUC *usecase.QuestUseCase
}

func NewQuestHandler(questUC *usecase.QuestUseCase) *QuestHandler {
	return &QuestHandler{questUC: questUC}
}

func (h *QuestHandler) GetQuest(w http.ResponseWriter, r *http.Request) {
	zoneID := r.PathValue("zoneId")
	if zoneID == "" {
		middleware.HandleDomainError(w, domain.ErrInvalidInput)
		return
	}

	mode := r.URL.Query().Get("mode")
	questView, err := h.questUC.GetShuffledQuestWithMode(r.Context(), zoneID, mode)
	if err != nil {
		middleware.HandleDomainError(w, err)
		return
	}

	resp := dto.QuestResponse{
		ID:               questView.ID,
		ZoneID:           questView.ZoneID,
		ZoneName:         questView.ZoneName,
		Difficulty:       questView.Difficulty,
		Mode:             questView.Mode,
		Title:            questView.Title,
		ProblemStatement: questView.ProblemStatement,
		CodeBlocks:       questView.CodeBlocks,
	}

	dto.WriteSuccess(w, http.StatusOK, resp)
}

func (h *QuestHandler) VerifyQuest(w http.ResponseWriter, r *http.Request) {
	zoneID := r.PathValue("zoneId")
	if zoneID == "" {
		middleware.HandleDomainError(w, domain.ErrInvalidInput)
		return
	}

	var req dto.VerifyQuestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleDomainError(w, domain.ErrInvalidInput)
		return
	}

	if err := req.Validate(); err != nil {
		middleware.HandleDomainError(w, err)
		return
	}

	userID := middleware.GetUserID(r.Context())

	result, err := h.questUC.VerifyQuestSubmissionWithMode(r.Context(), userID, zoneID, req.Mode, req.SubmittedOrder)
	if err != nil {
		middleware.HandleDomainError(w, err)
		return
	}

	resp := dto.VerifyQuestResponse{
		UnlockedZone: dto.ToZoneViewModel(result.UnlockedZone),
		Progress:     dto.ToProgressViewModel(result.ProgressMetrics),
	}

	dto.WriteSuccess(w, http.StatusOK, resp)
}
