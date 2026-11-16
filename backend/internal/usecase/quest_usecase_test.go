package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/xsolla-school/office-explorer/backend/internal/domain"
	"github.com/xsolla-school/office-explorer/backend/internal/domain/entity"
)

type mockZoneRepo struct {
	zones    map[string]entity.Zone
	hotspots []entity.Hotspot
}

func (m *mockZoneRepo) GetAllZones(_ context.Context) ([]entity.Zone, error) {
	res := make([]entity.Zone, 0, len(m.zones))
	for _, z := range m.zones {
		res = append(res, z)
	}
	return res, nil
}

func (m *mockZoneRepo) GetZoneByID(_ context.Context, zoneID string) (*entity.Zone, error) {
	z, ok := m.zones[zoneID]
	if !ok {
		return nil, domain.ErrZoneNotFound
	}
	return &z, nil
}

func (m *mockZoneRepo) GetAllHotspots(_ context.Context) ([]entity.Hotspot, error) {
	return m.hotspots, nil
}

func (m *mockZoneRepo) GetHotspotsByZoneID(_ context.Context, zoneID string) ([]entity.Hotspot, error) {
	res := make([]entity.Hotspot, 0)
	for _, h := range m.hotspots {
		if h.ZoneID != nil && *h.ZoneID == zoneID {
			res = append(res, h)
		}
	}
	return res, nil
}

func (m *mockZoneRepo) GetHotspotByID(_ context.Context, hotspotID string) (*entity.Hotspot, error) {
	for _, h := range m.hotspots {
		if h.ID == hotspotID {
			return &h, nil
		}
	}
	return nil, domain.ErrHotspotNotFound
}

type mockQuestRepo struct {
	quests map[string]entity.Quest
}

func (m *mockQuestRepo) GetQuestByZoneID(_ context.Context, zoneID string) (*entity.Quest, error) {
	q, ok := m.quests[zoneID]
	if !ok {
		return nil, domain.ErrQuestNotFound
	}
	return &q, nil
}

type mockProgressRepo struct {
	records map[string]*entity.UserProgress
}

func (m *mockProgressRepo) GetOrCreateProgress(_ context.Context, userID string) (*entity.UserProgress, error) {
	if p, ok := m.records[userID]; ok {
		return p, nil
	}
	p := entity.NewUserProgress(userID)
	m.records[userID] = p
	return p, nil
}

func (m *mockProgressRepo) SaveProgress(_ context.Context, progress *entity.UserProgress) error {
	m.records[progress.UserID] = progress
	return nil
}

func TestVerifyQuestSubmission_SuccessAndFailure(t *testing.T) {
	easy := entity.DifficultyEasy
	zoneID := entity.ZoneGameRoom1

	zoneRepo := &mockZoneRepo{
		zones: map[string]entity.Zone{
			zoneID: {
				ID:              zoneID,
				Name:            "Game Room — Sector 1 (Board Games)",
				IsLockedDefault: true,
				Difficulty:      &easy,
				MediaURL:        "/assets/rooms/gameroom_sector_1.jpg",
				Description:     "Board games area.",
			},
		},
	}

	questRepo := &mockQuestRepo{
		quests: map[string]entity.Quest{
			zoneID: {
				ID:               "q-1",
				ZoneID:           zoneID,
				Title:            "Swap Two Values",
				ProblemStatement: "Swap a and b",
				CodeBlocks: []entity.CodeBlock{
					{ID: "b1", Code: "temp := a"},
					{ID: "b2", Code: "a = b"},
					{ID: "b3", Code: "b = temp"},
				},
				CorrectOrder: []string{"b1", "b2", "b3"},
			},
		},
	}

	progressRepo := &mockProgressRepo{records: make(map[string]*entity.UserProgress)}
	uc := NewQuestUseCase(zoneRepo, questRepo, progressRepo)
	ctx := context.Background()
	userID := "00000000-0000-0000-0000-000000000001"

	_, err := uc.VerifyQuestSubmission(ctx, userID, zoneID, []string{"b3", "b1", "b2"})
	if !errors.Is(err, domain.ErrInvalidSequence) {
		t.Fatalf("expected ErrInvalidSequence, got %v", err)
	}

	res, err := uc.VerifyQuestSubmission(ctx, userID, zoneID, []string{"b1", "b2", "b3"})
	if err != nil {
		t.Fatalf("expected no error on valid sequence, got %v", err)
	}

	if res.UnlockedZone.Status != entity.ZoneStatusUnlocked {
		t.Fatalf("expected zone status UNLOCKED, got %s", res.UnlockedZone.Status)
	}

	if res.UnlockedZone.MediaURL == "" {
		t.Fatal("expected unmasked MediaURL after unlocking zone")
	}

	if res.ProgressMetrics.UnlockedZonesCount != 1 {
		t.Fatalf("expected 1 unlocked zone, got %d", res.ProgressMetrics.UnlockedZonesCount)
	}
}

func TestPuzzleMode_GetAndVerifySuccess(t *testing.T) {
	easy := entity.DifficultyEasy
	zoneID := entity.ZoneMeeting9

	zoneRepo := &mockZoneRepo{
		zones: map[string]entity.Zone{
			zoneID: {
				ID:              zoneID,
				Name:            "Meeting Room 9",
				IsLockedDefault: true,
				Difficulty:      &easy,
				MediaURL:        "/assets/rooms/meeting_9.jpg",
				Description:     "Presentation meeting room.",
			},
		},
	}

	questRepo := &mockQuestRepo{quests: map[string]entity.Quest{}}
	progressRepo := &mockProgressRepo{records: make(map[string]*entity.UserProgress)}
	uc := NewQuestUseCase(zoneRepo, questRepo, progressRepo)
	ctx := context.Background()
	userID := "00000000-0000-0000-0000-000000000002"

	view, err := uc.GetShuffledQuestWithMode(ctx, zoneID, QuestModePuzzle)
	if err != nil {
		t.Fatalf("expected puzzle quest, got %v", err)
	}
	if view.Mode != QuestModePuzzle {
		t.Fatalf("expected mode puzzle, got %s", view.Mode)
	}
	if len(view.CodeBlocks) != 4 {
		t.Fatalf("expected 4 puzzle blocks, got %d", len(view.CodeBlocks))
	}

	_, err = uc.VerifyQuestSubmissionWithMode(ctx, userID, zoneID, QuestModePuzzle, []string{"p4", "p3", "p2", "p1"})
	if !errors.Is(err, domain.ErrInvalidSequence) {
		t.Fatalf("expected ErrInvalidSequence for wrong puzzle order, got %v", err)
	}

	res, err := uc.VerifyQuestSubmissionWithMode(ctx, userID, zoneID, QuestModePuzzle, []string{"p1", "p2", "p3", "p4"})
	if err != nil {
		t.Fatalf("expected puzzle unlock success, got %v", err)
	}
	if res.UnlockedZone.Status != entity.ZoneStatusUnlocked {
		t.Fatalf("expected UNLOCKED zone, got %s", res.UnlockedZone.Status)
	}
}
