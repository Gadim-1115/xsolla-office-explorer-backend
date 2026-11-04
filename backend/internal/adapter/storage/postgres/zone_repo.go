package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xsolla-school/office-explorer/backend/internal/domain"
	"github.com/xsolla-school/office-explorer/backend/internal/domain/entity"
)

type ZoneRepo struct {
	pool *pgxpool.Pool
}

func NewZoneRepo(pool *pgxpool.Pool) *ZoneRepo {
	return &ZoneRepo{pool: pool}
}

func (r *ZoneRepo) GetAllZones(ctx context.Context) ([]entity.Zone, error) {
	const query = `
		SELECT id, name, is_locked_default, difficulty, media_url, description
		FROM zones
		ORDER BY id ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	zones := make([]entity.Zone, 0)
	for rows.Next() {
		var z entity.Zone
		var diffStr *string
		if err := rows.Scan(&z.ID, &z.Name, &z.IsLockedDefault, &diffStr, &z.MediaURL, &z.Description); err != nil {
			return nil, err
		}
		if diffStr != nil {
			d := entity.Difficulty(*diffStr)
			z.Difficulty = &d
		}
		zones = append(zones, z)
	}

	return zones, rows.Err()
}

func (r *ZoneRepo) GetZoneByID(ctx context.Context, zoneID string) (*entity.Zone, error) {
	const query = `
		SELECT id, name, is_locked_default, difficulty, media_url, description
		FROM zones
		WHERE id = $1
	`

	var z entity.Zone
	var diffStr *string
	err := r.pool.QueryRow(ctx, query, zoneID).Scan(
		&z.ID,
		&z.Name,
		&z.IsLockedDefault,
		&diffStr,
		&z.MediaURL,
		&z.Description,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrZoneNotFound
		}
		return nil, err
	}

	if diffStr != nil {
		d := entity.Difficulty(*diffStr)
		z.Difficulty = &d
	}

	return &z, nil
}

func (r *ZoneRepo) GetAllHotspots(ctx context.Context) ([]entity.Hotspot, error) {
	const query = `
		SELECT id, zone_id, category, x_percent, y_percent, title, content
		FROM hotspots
		ORDER BY id ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hotspots := make([]entity.Hotspot, 0)
	for rows.Next() {
		var h entity.Hotspot
		var catStr string
		if err := rows.Scan(&h.ID, &h.ZoneID, &catStr, &h.XPercent, &h.YPercent, &h.Title, &h.Content); err != nil {
			return nil, err
		}
		h.Category = entity.HotspotCategory(catStr)
		hotspots = append(hotspots, h)
	}

	return hotspots, rows.Err()
}

func (r *ZoneRepo) GetHotspotsByZoneID(ctx context.Context, zoneID string) ([]entity.Hotspot, error) {
	const query = `
		SELECT id, zone_id, category, x_percent, y_percent, title, content
		FROM hotspots
		WHERE zone_id = $1
		ORDER BY id ASC
	`

	rows, err := r.pool.Query(ctx, query, zoneID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hotspots := make([]entity.Hotspot, 0)
	for rows.Next() {
		var h entity.Hotspot
		var catStr string
		if err := rows.Scan(&h.ID, &h.ZoneID, &catStr, &h.XPercent, &h.YPercent, &h.Title, &h.Content); err != nil {
			return nil, err
		}
		h.Category = entity.HotspotCategory(catStr)
		hotspots = append(hotspots, h)
	}

	return hotspots, rows.Err()
}

func (r *ZoneRepo) GetHotspotByID(ctx context.Context, hotspotID string) (*entity.Hotspot, error) {
	const query = `
		SELECT id, zone_id, category, x_percent, y_percent, title, content
		FROM hotspots
		WHERE id = $1
	`

	var h entity.Hotspot
	var catStr string
	err := r.pool.QueryRow(ctx, query, hotspotID).Scan(
		&h.ID,
		&h.ZoneID,
		&catStr,
		&h.XPercent,
		&h.YPercent,
		&h.Title,
		&h.Content,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrHotspotNotFound
		}
		return nil, err
	}

	h.Category = entity.HotspotCategory(catStr)
	return &h, nil
}
