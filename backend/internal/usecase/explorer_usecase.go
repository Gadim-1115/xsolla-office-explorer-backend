package usecase

import (
	"context"

	"github.com/xsolla-school/office-explorer/backend/internal/domain"
	"github.com/xsolla-school/office-explorer/backend/internal/domain/entity"
)

type ZoneView struct {
	ID          string
	Name        string
	Status      entity.ZoneStatus
	Difficulty  *entity.Difficulty
	MediaURL    string
	Description string
	Hotspots    []HotspotView
}

type HotspotView struct {
	ID        string
	ZoneID    *string
	Category  entity.HotspotCategory
	XPercent  float64
	YPercent  float64
	Title     string
	Content   string
	IsClaimed bool
}

type ProgressMetrics struct {
	UserID               string
	UnlockedZonesCount   int
	TotalLockedZones     int
	FoundEggsCount       int
	TotalEasterEggs      int
	CompletionPercentage float64
	UnlockedZoneIDs      []string
	FoundEggIDs          []string
}

type MapStateResult struct {
	Zones           []ZoneView
	MapHotspots     []HotspotView
	ProgressMetrics ProgressMetrics
}

type ExplorerUseCase struct {
	zoneRepo     domain.ZoneRepository
	progressRepo domain.ProgressRepository
}

func NewExplorerUseCase(zoneRepo domain.ZoneRepository, progressRepo domain.ProgressRepository) *ExplorerUseCase {
	return &ExplorerUseCase{
		zoneRepo:     zoneRepo,
		progressRepo: progressRepo,
	}
}

func (uc *ExplorerUseCase) InitMapState(ctx context.Context, userID string) (*MapStateResult, error) {
	zones, err := uc.zoneRepo.GetAllZones(ctx)
	if err != nil {
		return nil, err
	}

	hotspots, err := uc.zoneRepo.GetAllHotspots(ctx)
	if err != nil {
		return nil, err
	}

	progress, err := uc.progressRepo.GetOrCreateProgress(ctx, userID)
	if err != nil {
		return nil, err
	}

	hotspotsByZone := make(map[string][]HotspotView)
	mapHotspots := make([]HotspotView, 0)
	totalEasterEggs := 0

	for _, h := range hotspots {
		if h.IsEasterEgg() {
			totalEasterEggs++
		}
		view := HotspotView{
			ID:        h.ID,
			ZoneID:    h.ZoneID,
			Category:  h.Category,
			XPercent:  h.XPercent,
			YPercent:  h.YPercent,
			Title:     h.Title,
			Content:   h.Content,
			IsClaimed: progress.HasFoundEgg(h.ID),
		}
		if h.ZoneID == nil {
			mapHotspots = append(mapHotspots, view)
		} else {
			hotspotsByZone[*h.ZoneID] = append(hotspotsByZone[*h.ZoneID], view)
		}
	}

	zoneViews := make([]ZoneView, 0, len(zones))
	totalLockedZones := 0

	for _, z := range zones {
		if z.IsLockedDefault {
			totalLockedZones++
		}
		unlocked := z.IsUnlockedFor(progress)
		zoneViews = append(zoneViews, BuildSanitizedZoneView(z, unlocked, hotspotsByZone[z.ID]))
	}

	metrics := buildProgressMetrics(progress, totalLockedZones, totalEasterEggs)

	return &MapStateResult{
		Zones:           zoneViews,
		MapHotspots:     mapHotspots,
		ProgressMetrics: metrics,
	}, nil
}

func (uc *ExplorerUseCase) GetZoneDetails(ctx context.Context, userID, zoneID string) (*ZoneView, error) {
	zone, err := uc.zoneRepo.GetZoneByID(ctx, zoneID)
	if err != nil {
		return nil, err
	}

	progress, err := uc.progressRepo.GetOrCreateProgress(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !zone.IsUnlockedFor(progress) {
		return nil, domain.ErrZoneLocked
	}

	hotspots, err := uc.zoneRepo.GetHotspotsByZoneID(ctx, zoneID)
	if err != nil {
		return nil, err
	}

	hotspotViews := make([]HotspotView, 0, len(hotspots))
	for _, h := range hotspots {
		hotspotViews = append(hotspotViews, HotspotView{
			ID:        h.ID,
			ZoneID:    h.ZoneID,
			Category:  h.Category,
			XPercent:  h.XPercent,
			YPercent:  h.YPercent,
			Title:     h.Title,
			Content:   h.Content,
			IsClaimed: progress.HasFoundEgg(h.ID),
		})
	}

	view := BuildSanitizedZoneView(*zone, true, hotspotViews)
	return &view, nil
}

func (uc *ExplorerUseCase) ResetProgress(ctx context.Context, userID string) (*MapStateResult, error) {
	fresh := entity.NewUserProgress(userID)
	if err := uc.progressRepo.SaveProgress(ctx, fresh); err != nil {
		return nil, err
	}
	return uc.InitMapState(ctx, userID)
}

func (uc *ExplorerUseCase) ClaimEasterEgg(ctx context.Context, userID, hotspotID string) (*ProgressMetrics, error) {
	hotspot, err := uc.zoneRepo.GetHotspotByID(ctx, hotspotID)
	if err != nil {
		return nil, err
	}

	if !hotspot.IsEasterEgg() {
		return nil, domain.ErrNotAnEasterEgg
	}

	progress, err := uc.progressRepo.GetOrCreateProgress(ctx, userID)
	if err != nil {
		return nil, err
	}

	if hotspot.ZoneID != nil {
		zone, err := uc.zoneRepo.GetZoneByID(ctx, *hotspot.ZoneID)
		if err != nil {
			return nil, err
		}
		if !zone.IsUnlockedFor(progress) {
			return nil, domain.ErrZoneLocked
		}
	}

	if progress.ClaimEasterEgg(hotspot.ID) {
		if err := uc.progressRepo.SaveProgress(ctx, progress); err != nil {
			return nil, err
		}
	}

	zones, err := uc.zoneRepo.GetAllZones(ctx)
	if err != nil {
		return nil, err
	}
	hotspots, err := uc.zoneRepo.GetAllHotspots(ctx)
	if err != nil {
		return nil, err
	}

	totalLockedZones := 0
	for _, z := range zones {
		if z.IsLockedDefault {
			totalLockedZones++
		}
	}
	totalEasterEggs := 0
	for _, h := range hotspots {
		if h.IsEasterEgg() {
			totalEasterEggs++
		}
	}

	metrics := buildProgressMetrics(progress, totalLockedZones, totalEasterEggs)
	return &metrics, nil
}

func BuildSanitizedZoneView(z entity.Zone, unlocked bool, hotspots []HotspotView) ZoneView {
	if !unlocked {
		return ZoneView{
			ID:          z.ID,
			Name:        z.Name,
			Status:      entity.ZoneStatusLocked,
			Difficulty:  z.Difficulty,
			MediaURL:    "",
			Description: "",
			Hotspots:    []HotspotView{},
		}
	}

	if hotspots == nil {
		hotspots = []HotspotView{}
	}

	return ZoneView{
		ID:          z.ID,
		Name:        z.Name,
		Status:      entity.ZoneStatusUnlocked,
		Difficulty:  z.Difficulty,
		MediaURL:    z.MediaURL,
		Description: z.Description,
		Hotspots:    hotspots,
	}
}

func buildProgressMetrics(progress *entity.UserProgress, totalLockedZones, totalEasterEggs int) ProgressMetrics {
	return ProgressMetrics{
		UserID:               progress.UserID,
		UnlockedZonesCount:   len(progress.UnlockedZones),
		TotalLockedZones:     totalLockedZones,
		FoundEggsCount:       len(progress.FoundEggs),
		TotalEasterEggs:      totalEasterEggs,
		CompletionPercentage: progress.CompletionPercentage(totalLockedZones, totalEasterEggs),
		UnlockedZoneIDs:      progress.UnlockedZones,
		FoundEggIDs:          progress.FoundEggs,
	}
}
