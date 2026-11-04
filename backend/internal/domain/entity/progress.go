package entity

import (
	"math"
	"slices"
	"time"
)

// UserProgress is the aggregate root tracking a student's unlocked zones and claimed Easter eggs.
type UserProgress struct {
	UserID        string
	UnlockedZones []string
	FoundEggs     []string
	UpdatedAt     time.Time
}

// NewUserProgress initializes an empty progress state for a student.
func NewUserProgress(userID string) *UserProgress {
	return &UserProgress{
		UserID:        userID,
		UnlockedZones: make([]string, 0),
		FoundEggs:     make([]string, 0),
		UpdatedAt:     time.Now().UTC(),
	}
}

// IsZoneUnlocked checks whether the student has already unlocked a specific zone ID.
func (p *UserProgress) IsZoneUnlocked(zoneID string) bool {
	return slices.Contains(p.UnlockedZones, zoneID)
}

// UnlockZone idempotently records a newly unlocked zone. Returns true if state changed.
func (p *UserProgress) UnlockZone(zoneID string) bool {
	if p.IsZoneUnlocked(zoneID) {
		return false
	}
	p.UnlockedZones = append(p.UnlockedZones, zoneID)
	p.UpdatedAt = time.Now().UTC()
	return true
}

// HasFoundEgg checks whether the student has already claimed a specific Easter egg hotspot ID.
func (p *UserProgress) HasFoundEgg(hotspotID string) bool {
	return slices.Contains(p.FoundEggs, hotspotID)
}

// ClaimEasterEgg idempotently records a found Easter egg. Returns true if state changed.
func (p *UserProgress) ClaimEasterEgg(hotspotID string) bool {
	if p.HasFoundEgg(hotspotID) {
		return false
	}
	p.FoundEggs = append(p.FoundEggs, hotspotID)
	p.UpdatedAt = time.Now().UTC()
	return true
}

// CompletionPercentage calculates overall onboarding progress (0 to 100) across locked zones and Easter eggs.
func (p *UserProgress) CompletionPercentage(totalLockedZones, totalEasterEggs int) float64 {
	totalItems := totalLockedZones + totalEasterEggs
	if totalItems <= 0 {
		return 100.0
	}
	completed := len(p.UnlockedZones) + len(p.FoundEggs)
	if completed > totalItems {
		completed = totalItems
	}
	raw := (float64(completed) / float64(totalItems)) * 100.0
	return math.Round(raw*100) / 100
}
