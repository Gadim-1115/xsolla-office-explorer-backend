package handler

import (
	"net/http"

	"github.com/xsolla-school/office-explorer/backend/internal/adapter/handler/middleware"
)

func NewRouter(
	explorerHandler *ExplorerHandler,
	questHandler *QuestHandler,
	allowedOrigin string,
	streakHandler ...*StreakHandler,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/explorer/init", explorerHandler.InitMap)
	mux.HandleFunc("GET /api/v1/explorer/zones/{zoneId}", explorerHandler.GetZoneDetails)
	mux.HandleFunc("POST /api/v1/explorer/easter-eggs/{hotspotId}/claim", explorerHandler.ClaimEasterEgg)

	mux.HandleFunc("GET /api/v1/explorer/quests/{zoneId}", questHandler.GetQuest)
	mux.HandleFunc("POST /api/v1/explorer/quests/{zoneId}/verify", questHandler.VerifyQuest)

	if len(streakHandler) > 0 && streakHandler[0] != nil {
		sh := streakHandler[0]
		mux.HandleFunc("GET /api/v1/explorer/streak", sh.GetStreak)
		mux.HandleFunc("POST /api/v1/explorer/streak/complete", sh.CompleteDailyBug)
		mux.HandleFunc("GET /api/v1/explorer/leaderboard", sh.GetLeaderboard)
	}

	withIdentity := middleware.StudentIdentity(mux)
	return middleware.CORS(allowedOrigin, withIdentity)
}
