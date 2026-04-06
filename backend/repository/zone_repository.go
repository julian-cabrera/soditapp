package repository

import (
	"context"
	"database/sql"
	"fmt"

	"soditapp/backend/model"
)

// ZoneRepository handles database operations for zones.
type ZoneRepository struct {
	db *sql.DB
}

// NewZoneRepository creates a new ZoneRepository.
func NewZoneRepository(db *sql.DB) *ZoneRepository {
	return &ZoneRepository{db: db}
}

// Create inserts a new zone and returns it.
func (r *ZoneRepository) Create(ctx context.Context, name string) (model.Zone, error) {
	var z model.Zone
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO zone (name) VALUES ($1) RETURNING id, name",
		name,
	).Scan(&z.ID, &z.Name)
	if err != nil {
		return z, fmt.Errorf("creating zone: %w", err)
	}
	return z, nil
}

// GetByID returns a zone by its ID.
func (r *ZoneRepository) GetByID(ctx context.Context, id int) (model.Zone, error) {
	var z model.Zone
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name FROM zone WHERE id = $1",
		id,
	).Scan(&z.ID, &z.Name)
	if err != nil {
		return z, fmt.Errorf("getting zone %d: %w", id, err)
	}
	return z, nil
}

// GetAll returns all zones ordered by name.
func (r *ZoneRepository) GetAll(ctx context.Context) ([]model.Zone, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM zone ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("listing zones: %w", err)
	}
	defer rows.Close()

	var zones []model.Zone
	for rows.Next() {
		var z model.Zone
		if err := rows.Scan(&z.ID, &z.Name); err != nil {
			return nil, fmt.Errorf("scanning zone: %w", err)
		}
		zones = append(zones, z)
	}
	return zones, rows.Err()
}

// Update modifies a zone's name.
func (r *ZoneRepository) Update(ctx context.Context, id int, name string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE zone SET name = $1 WHERE id = $2",
		name, id,
	)
	if err != nil {
		return fmt.Errorf("updating zone %d: %w", id, err)
	}
	return nil
}

// Delete removes a zone by its ID.
func (r *ZoneRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM zone WHERE id = $1",
		id,
	)
	if err != nil {
		return fmt.Errorf("deleting zone %d: %w", id, err)
	}
	return nil
}
