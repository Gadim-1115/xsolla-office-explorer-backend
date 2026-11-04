package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xsolla-school/office-explorer/backend/internal/domain/entity"
)

type ProgressRepo struct {
	pool *pgxpool.Pool
}

func NewProgressRepo(pool *pgxpool.Pool) *ProgressRepo {
	return &ProgressRepo{pool: pool}
}

func (r *ProgressRepo) GetOrCreateProgress(ctx context.Context, userID string) (*entity.UserProgress, error) {
	const selectQuery = `
		SELECT user_id::text, unlocked_zones, found_eggs, updated_at
		FROM user_progress
		WHERE user_id = $1::uuid
	`

	var p entity.UserProgress
	var unlockedJSON []byte
	var eggsJSON []byte

	err := r.pool.QueryRow(ctx, selectQuery, userID).Scan(
		&p.UserID,
		&unlockedJSON,
		&eggsJSON,
		&p.UpdatedAt,
	)
	if err == nil {
		if err := json.Unmarshal(unlockedJSON, &p.UnlockedZones); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(eggsJSON, &p.FoundEggs); err != nil {
			return nil, err
		}
		return &p, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	newProgress := entity.NewUserProgress(userID)
	if err := r.SaveProgress(ctx, newProgress); err != nil {
		return nil, err
	}

	return newProgress, nil
}

func (r *ProgressRepo) SaveProgress(ctx context.Context, progress *entity.UserProgress) error {
	unlockedJSON, err := json.Marshal(progress.UnlockedZones)
	if err != nil {
		return err
	}

	eggsJSON, err := json.Marshal(progress.FoundEggs)
	if err != nil {
		return err
	}

	const upsertQuery = `
		INSERT INTO user_progress (user_id, unlocked_zones, found_eggs, updated_at)
		VALUES ($1::uuid, $2::jsonb, $3::jsonb, $4)
		ON CONFLICT (user_id) DO UPDATE
		SET unlocked_zones = EXCLUDED.unlocked_zones,
		    found_eggs     = EXCLUDED.found_eggs,
		    updated_at     = EXCLUDED.updated_at
	`

	_, err = r.pool.Exec(ctx, upsertQuery, progress.UserID, unlockedJSON, eggsJSON, progress.UpdatedAt)
	return err
}
