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
	SolvedAt       *string `json:"solved_at"`
}

type LmsSyncUser struct {
	UserID    string  `json:"user_id"`
	UserName  string  `json:"user_name"`
	UserRole  string  `json:"user_role"`
	AvatarURL *string `json:"avatar_url"`
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
		`SELECT user_id, user_name, user_role, streak, last_solved_date, avatar_url,
		        TO_CHAR(updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"')
		 FROM daily_bug_streaks
		 WHERE user_id = $1`,
		userID,
	).Scan(&rec.UserID, &rec.UserName, &rec.UserRole, &rec.Streak, &rec.LastSolvedDate, &rec.AvatarURL, &rec.SolvedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return &StreakRecord{
			UserID:         userID,
			UserName:       "",
			UserRole:       "Student",
			Streak:         0,
			LastSolvedDate: nil,
			AvatarURL:      nil,
			SolvedAt:       nil,
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
	if userRole == "" {
		if existing.UserRole != "" {
			userRole = existing.UserRole
		} else {
			userRole = "Student"
		}
	}

	if existing.LastSolvedDate != nil && *existing.LastSolvedDate == dateKey {
		nameChanged := userName != "" && existing.UserName != userName
		roleChanged := userRole != "" && existing.UserRole != userRole
		avatarChanged := resolvedAvatar != nil && (existing.AvatarURL == nil || *existing.AvatarURL != *resolvedAvatar)
		if nameChanged || roleChanged || avatarChanged {
			if userName != "" {
				existing.UserName = userName
			}
			existing.UserRole = userRole
			existing.AvatarURL = resolvedAvatar
			_, _ = r.pool.Exec(
				ctx,
				`UPDATE daily_bug_streaks
				 SET user_name = $2,
				     user_role = $3,
				     avatar_url = $4
				 WHERE user_id = $1`,
				userID, existing.UserName, existing.UserRole, existing.AvatarURL,
			)
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

	var solvedAt string
	err = r.pool.QueryRow(
		ctx,
		`INSERT INTO daily_bug_streaks (user_id, user_name, user_role, streak, last_solved_date, avatar_url, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, NOW())
		 ON CONFLICT (user_id) DO UPDATE
		 SET user_name = EXCLUDED.user_name,
		     user_role = EXCLUDED.user_role,
		     streak = EXCLUDED.streak,
		     last_solved_date = EXCLUDED.last_solved_date,
		     avatar_url = COALESCE(EXCLUDED.avatar_url, daily_bug_streaks.avatar_url),
		     updated_at = NOW()
		 RETURNING TO_CHAR(updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"')`,
		userID, userName, userRole, nextStreak, dateKey, resolvedAvatar,
	).Scan(&solvedAt)
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
		SolvedAt:       &solvedAt,
	}, nil
}

func (r *StreakRepo) ListTopStreaks(ctx context.Context, limit int) ([]StreakRecord, error) {
	if limit <= 0 {
		limit = 10
	}
	rows, err := r.pool.Query(
		ctx,
		`SELECT user_id, user_name, user_role, streak, last_solved_date, avatar_url,
		        TO_CHAR(updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"')
		 FROM daily_bug_streaks
		 WHERE streak > 0 AND user_id NOT LIKE 'seed-%'
		 ORDER BY streak DESC, last_solved_date DESC NULLS LAST, updated_at ASC
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
			&rec.SolvedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, rec)
	}
	return result, rows.Err()
}

func (r *StreakRepo) DeleteUserStreak(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM daily_bug_streaks WHERE user_id = $1`, userID)
	return err
}

func (r *StreakRepo) SyncLmsUsers(ctx context.Context, users []LmsSyncUser, pruneMissing bool) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if pruneMissing {
		activeIDs := make([]string, 0, len(users))
		for _, u := range users {
			if u.UserID != "" {
				activeIDs = append(activeIDs, u.UserID)
			}
		}
		if len(activeIDs) == 0 {
			if _, err := tx.Exec(ctx, `DELETE FROM daily_bug_streaks WHERE user_id LIKE 'lms-%'`); err != nil {
				return err
			}
		} else {
			if _, err := tx.Exec(
				ctx,
				`DELETE FROM daily_bug_streaks WHERE user_id LIKE 'lms-%' AND NOT (user_id = ANY($1))`,
				activeIDs,
			); err != nil {
				return err
			}
		}
	}

	for _, u := range users {
		if u.UserID == "" || u.UserName == "" {
			continue
		}
		if _, err := tx.Exec(
			ctx,
			`UPDATE daily_bug_streaks
			 SET user_name = $2,
			     user_role = COALESCE(NULLIF($3, ''), user_role),
			     avatar_url = COALESCE(NULLIF($4, ''), avatar_url)
			 WHERE user_id = $1`,
			u.UserID,
			u.UserName,
			u.UserRole,
			u.AvatarURL,
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
