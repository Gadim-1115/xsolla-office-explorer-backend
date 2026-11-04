package domain

import (
	"context"

	"github.com/xsolla-school/office-explorer/backend/internal/domain/entity"
)

// ZoneRepository defines the persistence contract for office zones and hotspots.
type ZoneRepository interface {
	GetAllZones(ctx context.Context) ([]entity.Zone, error)
	GetZoneByID(ctx context.Context, zoneID string) (*entity.Zone, error)
	GetAllHotspots(ctx context.Context) ([]entity.Hotspot, error)
	GetHotspotsByZoneID(ctx context.Context, zoneID string) ([]entity.Hotspot, error)
	GetHotspotByID(ctx context.Context, hotspotID string) (*entity.Hotspot, error)
}

// QuestRepository defines the persistence contract for algorithmic quests.
type QuestRepository interface {
	GetQuestByZoneID(ctx context.Context, zoneID string) (*entity.Quest, error)
}

// ProgressRepository defines the persistence contract for student exploration progress.
type ProgressRepository interface {
	GetOrCreateProgress(ctx context.Context, userID string) (*entity.UserProgress, error)
	SaveProgress(ctx context.Context, progress *entity.UserProgress) error
}
