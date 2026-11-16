package usecase

import (
	"strings"

	"github.com/xsolla-school/office-explorer/backend/internal/domain/entity"
)

const (
	QuestModeAlgorithm = "algorithm"
	QuestModePuzzle    = "puzzle"
)

func NormalizeQuestMode(raw string) string {
	if strings.EqualFold(strings.TrimSpace(raw), QuestModePuzzle) {
		return QuestModePuzzle
	}
	return QuestModeAlgorithm
}

var puzzleQuestCatalog = map[string]entity.Quest{
	entity.ZoneGameRoom1: {
		ID:     "puzzle-gameroom-1",
		ZoneID: entity.ZoneGameRoom1,
		Title:  "Sector 1 Puzzle (Easy): Board Game Seating",
		ProblemStatement: "Four teammates sit in a row from Seat 1 (top) to Seat 4 (bottom). Arrange them using these clues:\n" +
			"• Ayla sits in the 1st seat (top).\n" +
			"• Tural sits in the 4th seat (bottom).\n" +
			"• Nigar sits directly above Tural.\n" +
			"• Murad sits directly below Ayla.",
		CodeBlocks: []entity.CodeBlock{
			{ID: "p1", Code: "Ayla — Board Game Host"},
			{ID: "p2", Code: "Murad — Strategy Lead"},
			{ID: "p3", Code: "Nigar — Card Tactician"},
			{ID: "p4", Code: "Tural — Dice Master"},
		},
		CorrectOrder: []string{"p1", "p2", "p3", "p4"},
	},
	entity.ZoneGameRoom2: {
		ID:     "puzzle-gameroom-2",
		ZoneID: entity.ZoneGameRoom2,
		Title:  "Sector 2 Puzzle (Medium): Arcade Number Cipher",
		ProblemStatement: "Unlock the Foosball & Air Hockey zone by ordering the 5 number tiles from top to bottom.\n" +
			"Pattern rule: Start at 2, then strictly alternate ×3 and +4 (Start 2 → ×3 → +4 → ×3 → +4).",
		CodeBlocks: []entity.CodeBlock{
			{ID: "p1", Code: "2  (Starting Number)"},
			{ID: "p2", Code: "6"},
			{ID: "p3", Code: "10"},
			{ID: "p4", Code: "30"},
			{ID: "p5", Code: "34"},
		},
		CorrectOrder: []string{"p1", "p2", "p3", "p4", "p5"},
	},
	entity.ZoneGameRoom3: {
		ID:     "puzzle-gameroom-3",
		ZoneID: entity.ZoneGameRoom3,
		Title:  "Sector 3 Puzzle (Medium-Hard): Tournament Ranking",
		ProblemStatement: "Five teams finished 1st to 5th in the Ping-Pong Championship. Order them from 1st Place (top) to 5th Place (bottom):\n" +
			"• Violet Team finished 1st overall.\n" +
			"• Red Team finished last (5th place).\n" +
			"• Amber Team finished directly below Violet Team and directly above Green Team.\n" +
			"• Blue Team finished directly above Red Team.",
		CodeBlocks: []entity.CodeBlock{
			{ID: "p1", Code: "Violet Team"},
			{ID: "p2", Code: "Amber Team"},
			{ID: "p3", Code: "Green Team"},
			{ID: "p4", Code: "Blue Team"},
			{ID: "p5", Code: "Red Team"},
		},
		CorrectOrder: []string{"p1", "p2", "p3", "p4", "p5"},
	},
	entity.ZoneGameRoom4: {
		ID:     "puzzle-gameroom-4",
		ZoneID: entity.ZoneGameRoom4,
		Title:  "Sector 4 Puzzle (Hard): Crescent Elevator Cipher",
		ProblemStatement: "Arrange the 6 Crescent Tower stops in exact express-elevator order from First Stop (top) to Last Stop (bottom):\n" +
			"• The elevator starts at Floor 10 and ends at Floor 31 (Xsolla Office).\n" +
			"• Floor 18 is visited immediately after Floor 10.\n" +
			"• Floor 25 is visited immediately before Floor 31.\n" +
			"• Between Floor 18 and Floor 25, it stops at Floor 22 first, then returns down to Floor 14 before climbing to Floor 25.",
		CodeBlocks: []entity.CodeBlock{
			{ID: "p1", Code: "Floor 10 — Main Lobby Gate"},
			{ID: "p2", Code: "Floor 18 — Sky Lounge"},
			{ID: "p3", Code: "Floor 22 — Co-Working Deck"},
			{ID: "p4", Code: "Floor 14 — Badge Checkpoint"},
			{ID: "p5", Code: "Floor 25 — Atrium Transfer"},
			{ID: "p6", Code: "Floor 31 — Xsolla Office"},
		},
		CorrectOrder: []string{"p1", "p2", "p3", "p4", "p5", "p6"},
	},
	entity.ZoneMeeting1: {
		ID:     "puzzle-meeting-1",
		ZoneID: entity.ZoneMeeting1,
		Title:  "Meeting Room 1 Puzzle (Easy): Morning Arrival",
		ProblemStatement: "Arrange the 4 steps of arriving at Xsolla Office in chronological order from first (top) to last (bottom).",
		CodeBlocks: []entity.CodeBlock{
			{ID: "p1", Code: "Scan your access badge at the 1st-floor turnstiles"},
			{ID: "p2", Code: "Take the North elevator up to the 31st floor"},
			{ID: "p3", Code: "Enter through the glass doors into Open Space 1"},
			{ID: "p4", Code: "Walk past the desks to reach Meeting Room 1"},
		},
		CorrectOrder: []string{"p1", "p2", "p3", "p4"},
	},
	entity.ZoneMeeting5: {
		ID:     "puzzle-meeting-5",
		ZoneID: entity.ZoneMeeting5,
		Title:  "Meeting Room 5 Puzzle (Easy): Daily Schedule Order",
		ProblemStatement: "Order Meeting Room 5's four daily bookings chronologically from earliest in the day (top) to latest in the evening (bottom).",
		CodeBlocks: []entity.CodeBlock{
			{ID: "p1", Code: "Quarter past 9 in the morning — Daily Kickoff"},
			{ID: "p2", Code: "Half past Noon — Team Lunch & Sync"},
			{ID: "p3", Code: "15:45 (Afternoon) — Design Walkthrough"},
			{ID: "p4", Code: "Ten past 6 in the evening — Evening Wrap-Up"},
		},
		CorrectOrder: []string{"p1", "p2", "p3", "p4"},
	},
	entity.ZoneMeeting9: {
		ID:     "puzzle-meeting-9",
		ZoneID: entity.ZoneMeeting9,
		Title:  "Meeting Room 9 Puzzle (Easy): Office Word-Chain",
		ProblemStatement: "Connect the 4 word-chain tiles from top to bottom so that the LAST word of each tile matches the FIRST word of the next tile.",
		CodeBlocks: []entity.CodeBlock{
			{ID: "p1", Code: "COFFEE  →  BREAK"},
			{ID: "p2", Code: "BREAK  →  ROOM"},
			{ID: "p3", Code: "ROOM  →  KEY"},
			{ID: "p4", Code: "KEY  →  CARD"},
		},
		CorrectOrder: []string{"p1", "p2", "p3", "p4"},
	},
}

func getPuzzleQuestForZone(zoneID string) (entity.Quest, bool) {
	q, ok := puzzleQuestCatalog[zoneID]
	return q, ok
}
