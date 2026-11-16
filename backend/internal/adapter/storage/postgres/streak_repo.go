package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StreakRecord struct {
	UserID         string  `json:"user_id"`
	UserName       string  `json:"user_name"`
	UserRole       string  `json:"user_role"`
	Streak         int     `json:"streak"`
	LastSolvedDate *string `json:"last_solved_date"`
	AvatarURL      *string `json:"avatar_url"`
}

type StreakRepo struct {
	pool *pgxpool.Pool
}

func NewStreakRepo(ctx context.Context, pool *pgxpool.Pool) *StreakRepo {
	repo := &StreakRepo{pool: pool}
	_ = repo.ensureSchema(ctx)
	return repo
}

func (r *StreakRepo) ensureSchema(ctx context.Context) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS daily_bug_streaks (
	user_id VARCHAR(128) PRIMARY KEY,
	user_name VARCHAR(255) NOT NULL,
	user_role VARCHAR(64) NOT NULL DEFAULT 'Student',
	streak INT NOT NULL DEFAULT 0,
	last_solved_date VARCHAR(16),
	avatar_url TEXT,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
ALTER TABLE daily_bug_streaks ADD COLUMN IF NOT EXISTS avatar_url TEXT;
DELETE FROM daily_bug_streaks WHERE user_id LIKE 'seed-%';`
	_, err := r.pool.Exec(ctx, ddl)
	return err
}

func (r *StreakRepo) GetUserStreak(ctx context.Context, userID string) (*StreakRecord, error) {
	var rec StreakRecord
	err := r.pool.QueryRow(
		ctx,
		`SELECT user_id, user_name, user_role, streak, last_solved_date, avatar_url
		 FROM daily_bug_streaks
		 WHERE user_id = $1`,
		userID,
	).Scan(&rec.UserID, &rec.UserName, &rec.UserRole, &rec.Streak, &rec.LastSolvedDate, &rec.AvatarURL)
	if errors.Is(err, pgx.ErrNoRows) {
		return &StreakRecord{
			UserID:         userID,
			UserName:       "",
			UserRole:       "Student",
			Streak:         0,
			LastSolvedDate: nil,
			AvatarURL:      nil,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *StreakRepo) CompleteDailyBug(
	ctx context.Context,
	userID, userName, userRole, dateKey string,
	avatarURL *string,
) (*StreakRecord, error) {
	existing, err := r.GetUserStreak(ctx, userID)
	if err != nil {
		return nil, err
	}

	resolvedAvatar := avatarURL
	if resolvedAvatar == nil || *resolvedAvatar == "" {
		resolvedAvatar = existing.AvatarURL
	}

	if existing.LastSolvedDate != nil && *existing.LastSolvedDate == dateKey {
		if resolvedAvatar != nil && (existing.AvatarURL == nil || *existing.AvatarURL != *resolvedAvatar) {
			_, _ = r.pool.Exec(
				ctx,
				`UPDATE daily_bug_streaks SET avatar_url = $2, updated_at = NOW() WHERE user_id = $1`,
				userID, *resolvedAvatar,
			)
			existing.AvatarURL = resolvedAvatar
		}
		return existing, nil
	}

	parsedDate, err := time.Parse("2006-01-02", dateKey)
	if err != nil {
		parsedDate = time.Now()
		dateKey = parsedDate.Format("2006-01-02")
	}
	yesterdayKey := parsedDate.AddDate(0, 0, -1).Format("2006-01-02")

	nextStreak := 1
	if existing.LastSolvedDate != nil && *existing.LastSolvedDate == yesterdayKey && existing.Streak > 0 {
		nextStreak = existing.Streak + 1
	}

	if userRole == "" {
		userRole = "Student"
	}

	_, err = r.pool.Exec(
		ctx,
		`INSERT INTO daily_bug_streaks (user_id, user_name, user_role, streak, last_solved_date, avatar_url, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, NOW())
		 ON CONFLICT (user_id) DO UPDATE
		 SET user_name = EXCLUDED.user_name,
		     user_role = EXCLUDED.user_role,
		     streak = EXCLUDED.streak,
		     last_solved_date = EXCLUDED.last_solved_date,
		     avatar_url = COALESCE(EXCLUDED.avatar_url, daily_bug_streaks.avatar_url),
		     updated_at = NOW()`,
		userID, userName, userRole, nextStreak, dateKey, resolvedAvatar,
	)
	if err != nil {
		return nil, err
	}

	return &StreakRecord{
		UserID:         userID,
		UserName:       userName,
		UserRole:       userRole,
		Streak:         nextStreak,
		LastSolvedDate: &dateKey,
		AvatarURL:      resolvedAvatar,
	}, nil
}

func (r *StreakRepo) ListTopStreaks(ctx context.Context, limit int) ([]StreakRecord, error) {
	if limit <= 0 {
		limit = 10
	}
	rows, err := r.pool.Query(
		ctx,
		`SELECT user_id, user_name, user_role, streak, last_solved_date, avatar_url
		 FROM daily_bug_streaks
		 WHERE streak > 0 AND user_id NOT LIKE 'seed-%'
		 ORDER BY streak DESC, updated_at DESC
		 LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []StreakRecord
	for rows.Next() {
		var rec StreakRecord
		if err := rows.Scan(
			&rec.UserID,
			&rec.UserName,
			&rec.UserRole,
			&rec.Streak,
			&rec.LastSolvedDate,
			&rec.AvatarURL,
		); err != nil {
			return nil, err
		}
		result = append(result, rec)
	}
	return result, rows.Err()
}
