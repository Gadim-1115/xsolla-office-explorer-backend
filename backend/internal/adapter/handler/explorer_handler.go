package handler

import (
	"net/http"

	"github.com/xsolla-school/office-explorer/backend/internal/adapter/handler/dto"
	"github.com/xsolla-school/office-explorer/backend/internal/adapter/handler/middleware"
	"github.com/xsolla-school/office-explorer/backend/internal/domain"
	"github.com/xsolla-school/office-explorer/backend/internal/usecase"
)

type ExplorerHandler struct {
	explorerUC *usecase.ExplorerUseCase
}

func NewExplorerHandler(explorerUC *usecase.ExplorerUseCase) *ExplorerHandler {
	return &ExplorerHandler{explorerUC: explorerUC}
}

func (h *ExplorerHandler) InitMap(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	result, err := h.explorerUC.InitMapState(r.Context(), userID)
	if err != nil {
		middleware.HandleDomainError(w, err)
		return
	}

	zones := make([]dto.ZoneViewModel, 0, len(result.Zones))
	for _, z := range result.Zones {
		zones = append(zones, dto.ToZoneViewModel(z))
	}

	mapHotspots := make([]dto.HotspotViewModel, 0, len(result.MapHotspots))
	for _, hs := range result.MapHotspots {
		mapHotspots = append(mapHotspots, dto.ToHotspotViewModel(hs))
	}

	resp := dto.InitMapResponse{
		Zones:       zones,
		MapHotspots: mapHotspots,
		Progress:    dto.ToProgressViewModel(result.ProgressMetrics),
	}

	dto.WriteSuccess(w, http.StatusOK, resp)
}

func (h *ExplorerHandler) GetZoneDetails(w http.ResponseWriter, r *http.Request) {
	zoneID := r.PathValue("zoneId")
	if zoneID == "" {
		middleware.HandleDomainError(w, domain.ErrInvalidInput)
		return
	}

	userID := middleware.GetUserID(r.Context())

	zoneView, err := h.explorerUC.GetZoneDetails(r.Context(), userID, zoneID)
	if err != nil {
		middleware.HandleDomainError(w, err)
		return
	}

	dto.WriteSuccess(w, http.StatusOK, dto.ToZoneViewModel(*zoneView))
}

func (h *ExplorerHandler) ClaimEasterEgg(w http.ResponseWriter, r *http.Request) {
	hotspotID := r.PathValue("hotspotId")
	if hotspotID == "" {
		middleware.HandleDomainError(w, domain.ErrInvalidInput)
		return
	}

	userID := middleware.GetUserID(r.Context())

	metrics, err := h.explorerUC.ClaimEasterEgg(r.Context(), userID, hotspotID)
	if err != nil {
		middleware.HandleDomainError(w, err)
		return
	}

	dto.WriteSuccess(w, http.StatusOK, dto.ToProgressViewModel(*metrics))
}

func (h *ExplorerHandler) ResetSession(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	result, err := h.explorerUC.ResetProgress(r.Context(), userID)
	if err != nil {
		middleware.HandleDomainError(w, err)
		return
	}

	zones := make([]dto.ZoneViewModel, 0, len(result.Zones))
	for _, z := range result.Zones {
		zones = append(zones, dto.ToZoneViewModel(z))
	}

	mapHotspots := make([]dto.HotspotViewModel, 0, len(result.MapHotspots))
	for _, hs := range result.MapHotspots {
		mapHotspots = append(mapHotspots, dto.ToHotspotViewModel(hs))
	}

	resp := dto.InitMapResponse{
		Zones:       zones,
		MapHotspots: mapHotspots,
		Progress:    dto.ToProgressViewModel(result.ProgressMetrics),
	}

	dto.WriteSuccess(w, http.StatusOK, resp)
}
