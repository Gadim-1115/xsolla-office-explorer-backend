package dto

import (
	"encoding/json"
	"net/http"

	"github.com/xsolla-school/office-explorer/backend/internal/domain/entity"
	"github.com/xsolla-school/office-explorer/backend/internal/usecase"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Envelope struct {
	Success bool      `json:"success"`
	Data    any       `json:"data"`
	Error   *APIError `json:"error"`
}

func WriteSuccess(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{
		Success: true,
		Data:    data,
		Error:   nil,
	})
}

func WriteError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{
		Success: false,
		Data:    nil,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

type HotspotViewModel struct {
	ID        string                 `json:"id"`
	ZoneID    *string                `json:"zone_id"`
	Category  entity.HotspotCategory `json:"category"`
	XPercent  float64                `json:"x_percent"`
	YPercent  float64                `json:"y_percent"`
	Title     string                 `json:"title"`
	Content   string                 `json:"content"`
	IsClaimed bool                   `json:"is_claimed"`
}

type ZoneViewModel struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Status      entity.ZoneStatus  `json:"status"`
	Difficulty  *entity.Difficulty `json:"difficulty"`
	MediaURL    string             `json:"media_url,omitempty"`
	Description string             `json:"description,omitempty"`
	Hotspots    []HotspotViewModel `json:"hotspots"`
}

type ProgressViewModel struct {
	UserID               string   `json:"user_id"`
	UnlockedZonesCount   int      `json:"unlocked_zones_count"`
	TotalLockedZones     int      `json:"total_locked_zones"`
	FoundEggsCount       int      `json:"found_eggs_count"`
	TotalEasterEggs      int      `json:"total_easter_eggs"`
	CompletionPercentage float64  `json:"completion_percentage"`
	UnlockedZoneIDs      []string `json:"unlocked_zone_ids"`
	FoundEggIDs          []string `json:"found_egg_ids"`
}

type InitMapResponse struct {
	Zones       []ZoneViewModel    `json:"zones"`
	MapHotspots []HotspotViewModel `json:"map_hotspots"`
	Progress    ProgressViewModel  `json:"progress"`
}

type QuestResponse struct {
	ID               string             `json:"id"`
	ZoneID           string             `json:"zone_id"`
	ZoneName         string             `json:"zone_name"`
	Difficulty       *entity.Difficulty `json:"difficulty"`
	Mode             string             `json:"mode"`
	Title            string             `json:"title"`
	ProblemStatement string             `json:"problem_statement"`
	CodeBlocks       []entity.CodeBlock `json:"code_blocks"`
}

type VerifyQuestResponse struct {
	UnlockedZone ZoneViewModel     `json:"unlocked_zone"`
	Progress     ProgressViewModel `json:"progress"`
}

func ToZoneViewModel(v usecase.ZoneView) ZoneViewModel {
	hotspots := make([]HotspotViewModel, 0, len(v.Hotspots))
	for _, h := range v.Hotspots {
		hotspots = append(hotspots, ToHotspotViewModel(h))
	}
	return ZoneViewModel{
		ID:          v.ID,
		Name:        v.Name,
		Status:      v.Status,
		Difficulty:  v.Difficulty,
		MediaURL:    v.MediaURL,
		Description: v.Description,
		Hotspots:    hotspots,
	}
}

func ToHotspotViewModel(h usecase.HotspotView) HotspotViewModel {
	return HotspotViewModel{
		ID:        h.ID,
		ZoneID:    h.ZoneID,
		Category:  h.Category,
		XPercent:  h.XPercent,
		YPercent:  h.YPercent,
		Title:     h.Title,
		Content:   h.Content,
		IsClaimed: h.IsClaimed,
	}
}

func ToProgressViewModel(m usecase.ProgressMetrics) ProgressViewModel {
	return ProgressViewModel{
		UserID:               m.UserID,
		UnlockedZonesCount:   m.UnlockedZonesCount,
		TotalLockedZones:     m.TotalLockedZones,
		FoundEggsCount:       m.FoundEggsCount,
		TotalEasterEggs:      m.TotalEasterEggs,
		CompletionPercentage: m.CompletionPercentage,
		UnlockedZoneIDs:      m.UnlockedZoneIDs,
		FoundEggIDs:          m.FoundEggIDs,
	}
}
