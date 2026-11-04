package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xsolla-school/office-explorer/backend/internal/domain"
	"github.com/xsolla-school/office-explorer/backend/internal/domain/entity"
)

type QuestRepo struct {
	pool *pgxpool.Pool
}

func NewQuestRepo(pool *pgxpool.Pool) *QuestRepo {
	return &QuestRepo{pool: pool}
}

func (r *QuestRepo) GetQuestByZoneID(ctx context.Context, zoneID string) (*entity.Quest, error) {
	const query = `
		SELECT id::text, zone_id, title, problem_statement, code_blocks, correct_order
		FROM quests
		WHERE zone_id = $1
	`

	var q entity.Quest
	var blocksJSON []byte
	var orderJSON []byte

	err := r.pool.QueryRow(ctx, query, zoneID).Scan(
		&q.ID,
		&q.ZoneID,
		&q.Title,
		&q.ProblemStatement,
		&blocksJSON,
		&orderJSON,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrQuestNotFound
		}
		return nil, err
	}

	if err := json.Unmarshal(blocksJSON, &q.CodeBlocks); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(orderJSON, &q.CorrectOrder); err != nil {
		return nil, err
	}

	return &q, nil
}
