package entity

// Difficulty represents the algorithmic quest difficulty tier of a locked zone.
type Difficulty string

const (
	DifficultyEasy       Difficulty = "EASY"
	DifficultyMedium     Difficulty = "MEDIUM"
	DifficultyMediumHard Difficulty = "MEDIUM_HARD"
	DifficultyHard       Difficulty = "HARD"
)

// IsValid checks whether the difficulty value matches one of the allowed tiers.
func (d Difficulty) IsValid() bool {
	switch d {
	case DifficultyEasy, DifficultyMedium, DifficultyMediumHard, DifficultyHard:
		return true
	default:
		return false
	}
}

// ZoneStatus represents the runtime lock state of a zone for a specific student.
type ZoneStatus string

const (
	ZoneStatusLocked   ZoneStatus = "LOCKED"
	ZoneStatusUnlocked ZoneStatus = "UNLOCKED"
)

// Predefined SVG Zone IDs strictly matching the latest Xsolla.svg on the 31st-floor oval map.
const (
	ZoneGameRoom1   = "zone_gameroom_1"
	ZoneGameRoom2   = "zone_gameroom_2"
	ZoneGameRoom3   = "zone_gameroom_3"
	ZoneGameRoom4   = "zone_gameroom_4"
	ZoneMeeting1    = "zone_meeting_1"
	ZoneMeeting2    = "zone_meeting_2"
	ZoneMeeting3    = "zone_meeting_3"
	ZoneMeeting4    = "zone_meeting_4"
	ZoneMeeting5    = "zone_meeting_5"
	ZoneMeeting6    = "zone_meeting_6"
	ZoneMeeting7    = "zone_meeting_7"
	ZoneMeeting8    = "zone_meeting_8"
	ZoneMeeting9    = "zone_meeting_9"
	ZoneMeeting10   = "zone_meeting_10"
	ZoneMeeting11   = "zone_meeting_11"
	ZoneOpenSpace1  = "zone_lecturezone_open_space_1" // Xsolla School students workspace
	ZoneOpenSpace2  = "zone_open_space_2"             // Xsolla employees workspace
	ZoneSession     = "zone_session"
	ZoneKitchen     = "zone_kitchen"
	ZoneStock       = "zone_stock"
)

// Zone represents a physical area or room on the 31st-floor oval office map.
type Zone struct {
	ID              string
	Name            string
	IsLockedDefault bool
	Difficulty      *Difficulty
	MediaURL        string
	Description     string
}

// IsUnlockedFor determines if the zone is accessible given the student's progress.
func (z *Zone) IsUnlockedFor(progress *UserProgress) bool {
	if !z.IsLockedDefault {
		return true
	}
	if progress == nil {
		return false
	}
	return progress.IsZoneUnlocked(z.ID)
}
