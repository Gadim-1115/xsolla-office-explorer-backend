package usecase

import (
	"context"

	"github.com/xsolla-school/office-explorer/backend/internal/domain"
	"github.com/xsolla-school/office-explorer/backend/internal/domain/entity"
	"github.com/xsolla-school/office-explorer/backend/internal/pkg/shuffler"
)

type SanitizedQuestView struct {
	ID               string
	ZoneID           string
	ZoneName         string
	Difficulty       *entity.Difficulty
	Mode             string
	Title            string
	ProblemStatement string
	CodeBlocks       []entity.CodeBlock
}

type QuestVerificationResult struct {
	UnlockedZone    ZoneView
	ProgressMetrics ProgressMetrics
}

type QuestUseCase struct {
	zoneRepo     domain.ZoneRepository
	questRepo    domain.QuestRepository
	progressRepo domain.ProgressRepository
}

func NewQuestUseCase(
	zoneRepo domain.ZoneRepository,
	questRepo domain.QuestRepository,
	progressRepo domain.ProgressRepository,
) *QuestUseCase {
	return &QuestUseCase{
		zoneRepo:     zoneRepo,
		questRepo:    questRepo,
		progressRepo: progressRepo,
	}
}

func (uc *QuestUseCase) resolveQuest(ctx context.Context, zoneID, mode string) (*entity.Quest, string, error) {
	normalized := NormalizeQuestMode(mode)
	if normalized == QuestModePuzzle {
		if puzzle, ok := getPuzzleQuestForZone(zoneID); ok {
			return &puzzle, QuestModePuzzle, nil
		}
	}
	quest, err := uc.questRepo.GetQuestByZoneID(ctx, zoneID)
	if err != nil {
		return nil, "", err
	}
	return quest, QuestModeAlgorithm, nil
}

func (uc *QuestUseCase) GetShuffledQuest(ctx context.Context, zoneID string) (*SanitizedQuestView, error) {
	return uc.GetShuffledQuestWithMode(ctx, zoneID, QuestModeAlgorithm)
}

func (uc *QuestUseCase) GetShuffledQuestWithMode(ctx context.Context, zoneID, mode string) (*SanitizedQuestView, error) {
	zone, err := uc.zoneRepo.GetZoneByID(ctx, zoneID)
	if err != nil {
		return nil, err
	}

	quest, resolvedMode, err := uc.resolveQuest(ctx, zoneID, mode)
	if err != nil {
		return nil, err
	}

	shuffledBlocks := shuffler.ShuffleCodeBlocks(quest.CodeBlocks, quest.CorrectOrder)

	return &SanitizedQuestView{
		ID:               quest.ID,
		ZoneID:           quest.ZoneID,
		ZoneName:         zone.Name,
		Difficulty:       zone.Difficulty,
		Mode:             resolvedMode,
		Title:            quest.Title,
		ProblemStatement: quest.ProblemStatement,
		CodeBlocks:       shuffledBlocks,
	}, nil
}

func (uc *QuestUseCase) VerifyQuestSubmission(
	ctx context.Context,
	userID string,
	zoneID string,
	submittedOrder []string,
) (*QuestVerificationResult, error) {
	return uc.VerifyQuestSubmissionWithMode(ctx, userID, zoneID, QuestModeAlgorithm, submittedOrder)
}

func (uc *QuestUseCase) VerifyQuestSubmissionWithMode(
	ctx context.Context,
	userID string,
	zoneID string,
	mode string,
	submittedOrder []string,
) (*QuestVerificationResult, error) {
	if len(submittedOrder) == 0 {
		return nil, domain.ErrInvalidInput
	}

	zone, err := uc.zoneRepo.GetZoneByID(ctx, zoneID)
	if err != nil {
		return nil, err
	}

	quest, _, err := uc.resolveQuest(ctx, zoneID, mode)
	if err != nil {
		return nil, err
	}

	if !quest.VerifySequence(submittedOrder) {
		return nil, domain.ErrInvalidSequence
	}

	progress, err := uc.progressRepo.GetOrCreateProgress(ctx, userID)
	if err != nil {
		return nil, err
	}

	if progress.UnlockZone(zone.ID) {
		if err := uc.progressRepo.SaveProgress(ctx, progress); err != nil {
			return nil, err
		}
	}

	hotspots, err := uc.zoneRepo.GetHotspotsByZoneID(ctx, zone.ID)
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

	allZones, err := uc.zoneRepo.GetAllZones(ctx)
	if err != nil {
		return nil, err
	}
	allHotspots, err := uc.zoneRepo.GetAllHotspots(ctx)
	if err != nil {
		return nil, err
	}

	totalLockedZones := 0
	for _, z := range allZones {
		if z.IsLockedDefault {
			totalLockedZones++
		}
	}
	totalEasterEggs := 0
	for _, h := range allHotspots {
		if h.IsEasterEgg() {
			totalEasterEggs++
		}
	}

	unlockedZoneView := BuildSanitizedZoneView(*zone, true, hotspotViews)
	metrics := buildProgressMetrics(progress, totalLockedZones, totalEasterEggs)

	return &QuestVerificationResult{
		UnlockedZone:    unlockedZoneView,
		ProgressMetrics: metrics,
	}, nil
}
