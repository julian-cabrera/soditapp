package repository

import (
	"context"
	"database/sql"
	"fmt"

	"soditapp/backend/model"
)

// ClientRepository handles database operations for clients.
type ClientRepository struct {
	db *sql.DB
}

// NewClientRepository creates a new ClientRepository.
func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

// Create inserts a new client and returns it.
func (r *ClientRepository) Create(ctx context.Context, name, address string, zoneID int) (model.Client, error) {
	var c model.Client
	var addr sql.NullString
	var unregDate sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`INSERT INTO client (name, address, zone_id)
		 VALUES ($1, $2, $3)
		 RETURNING id, name, address, zone_id, create_date, unregister_date`,
		name, address, zoneID,
	).Scan(&c.ID, &c.Name, &addr, &c.ZoneID, &c.CreateDate, &unregDate)
	if err != nil {
		return c, fmt.Errorf("creating client: %w", err)
	}

	if addr.Valid {
		c.Address = addr.String
	}
	if unregDate.Valid {
		c.UnregisterDate = &unregDate.Time
	}

	return c, nil
}

// GetByID returns a client by its ID.
func (r *ClientRepository) GetByID(ctx context.Context, id int) (model.Client, error) {
	var c model.Client
	var addr sql.NullString
	var unregDate sql.NullTime

	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, address, zone_id, create_date, unregister_date FROM client WHERE id = $1",
		id,
	).Scan(&c.ID, &c.Name, &addr, &c.ZoneID, &c.CreateDate, &unregDate)
	if err != nil {
		return c, fmt.Errorf("getting client %d: %w", id, err)
	}

	if addr.Valid {
		c.Address = addr.String
	}
	if unregDate.Valid {
		c.UnregisterDate = &unregDate.Time
	}

	return c, nil
}

// GetByZone returns active (not unregistered) clients in a zone, ordered by name.
func (r *ClientRepository) GetByZone(ctx context.Context, zoneID int) ([]model.Client, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, address, zone_id, create_date, unregister_date
		 FROM client
		 WHERE zone_id = $1 AND unregister_date IS NULL
		 ORDER BY name`,
		zoneID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing clients for zone %d: %w", zoneID, err)
	}
	defer rows.Close()

	return scanClients(rows)
}

// GetAll returns all active (not unregistered) clients ordered by name.
func (r *ClientRepository) GetAll(ctx context.Context) ([]model.Client, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, address, zone_id, create_date, unregister_date
		 FROM client
		 WHERE unregister_date IS NULL
		 ORDER BY name`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing clients: %w", err)
	}
	defer rows.Close()

	return scanClients(rows)
}

// Update modifies a client's name, address, and zone.
func (r *ClientRepository) Update(ctx context.Context, client model.Client) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE client SET name = $1, address = $2, zone_id = $3 WHERE id = $4",
		client.Name, client.Address, client.ZoneID, client.ID,
	)
	if err != nil {
		return fmt.Errorf("updating client %d: %w", client.ID, err)
	}
	return nil
}

// Unregister soft-deletes a client by setting unregister_date to today.
func (r *ClientRepository) Unregister(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE client SET unregister_date = CURRENT_DATE WHERE id = $1",
		id,
	)
	if err != nil {
		return fmt.Errorf("unregistering client %d: %w", id, err)
	}
	return nil
}

// Delete removes a client by its ID.
func (r *ClientRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM client WHERE id = $1",
		id,
	)
	if err != nil {
		return fmt.Errorf("deleting client %d: %w", id, err)
	}
	return nil
}

// BulkUpdateZone moves a set of clients to a new zone in a single transaction.
func (r *ClientRepository) BulkUpdateZone(ctx context.Context, clientIDs []int, targetZoneID int) error {
	if len(clientIDs) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning bulk update zone tx: %w", err)
	}
	defer tx.Rollback()

	// Build: UPDATE client SET zone_id = $1 WHERE id IN ($2, $3, ...)
	placeholders := make([]string, len(clientIDs))
	args := make([]interface{}, len(clientIDs)+1)
	args[0] = targetZoneID
	for i, id := range clientIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = id
	}
	query := fmt.Sprintf(
		"UPDATE client SET zone_id = $1 WHERE id IN (%s)",
		joinStrings(placeholders, ", "),
	)

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("bulk updating zone: %w", err)
	}

	return tx.Commit()
}

// joinStrings joins a slice of strings with a separator.
func joinStrings(parts []string, sep string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += sep
		}
		result += p
	}
	return result
}

// scanClients reads client rows and handles nullable columns.
func scanClients(rows *sql.Rows) ([]model.Client, error) {
	var clients []model.Client
	for rows.Next() {
		var c model.Client
		var addr sql.NullString
		var unregDate sql.NullTime

		if err := rows.Scan(&c.ID, &c.Name, &addr, &c.ZoneID, &c.CreateDate, &unregDate); err != nil {
			return nil, fmt.Errorf("scanning client: %w", err)
		}

		if addr.Valid {
			c.Address = addr.String
		}
		if unregDate.Valid {
			c.UnregisterDate = &unregDate.Time
		}

		clients = append(clients, c)
	}
	return clients, rows.Err()
}
