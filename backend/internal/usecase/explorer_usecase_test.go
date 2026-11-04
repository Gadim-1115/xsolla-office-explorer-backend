package usecase

import (
	"context"
	"testing"

	"github.com/xsolla-school/office-explorer/backend/internal/domain/entity"
)

func TestInitMapState_MasksLockedZoneMediaAndClaimsEasterEgg(t *testing.T) {
	easy := entity.DifficultyEasy
	lockedID := entity.ZoneGameRoom1
	unlockedID := entity.ZoneKitchen

	zoneRepo := &mockZoneRepo{
		zones: map[string]entity.Zone{
			lockedID: {
				ID:              lockedID,
				Name:            "Game Room 1",
				IsLockedDefault: true,
				Difficulty:      &easy,
				MediaURL:        "/assets/rooms/secret.jpg",
				Description:     "Secret description",
			},
			unlockedID: {
				ID:              unlockedID,
				Name:            "Kitchen",
				IsLockedDefault: false,
				MediaURL:        "/assets/rooms/kitchen.jpg",
				Description:     "Kitchen description",
			},
		},
		hotspots: []entity.Hotspot{
			{
				ID:       "hotspot_egg_kitchen",
				ZoneID:   &unlockedID,
				Category: entity.CategoryEasterEgg,
				XPercent: 50,
				YPercent: 50,
				Title:    "Hidden Logo",
				Content:  "Found!",
			},
		},
	}

	progressRepo := &mockProgressRepo{records: make(map[string]*entity.UserProgress)}
	uc := NewExplorerUseCase(zoneRepo, progressRepo)
	ctx := context.Background()
	userID := "00000000-0000-0000-0000-000000000001"

	state, err := uc.InitMapState(ctx, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, z := range state.Zones {
		if z.ID == lockedID && z.MediaURL != "" {
			t.Fatal("locked zone MediaURL must be masked")
		}
		if z.ID == unlockedID && z.MediaURL == "" {
			t.Fatal("unlocked zone MediaURL must not be masked")
		}
	}

	metrics, err := uc.ClaimEasterEgg(ctx, userID, "hotspot_egg_kitchen")
	if err != nil {
		t.Fatalf("unexpected error claiming egg: %v", err)
	}

	if metrics.FoundEggsCount != 1 {
		t.Fatalf("expected 1 found egg, got %d", metrics.FoundEggsCount)
	}

	if metrics.CompletionPercentage != 50.0 {
		t.Fatalf("expected 50%% completion, got %v", metrics.CompletionPercentage)
	}
}
