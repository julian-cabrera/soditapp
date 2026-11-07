package repository

import (
	"context"
	"database/sql"
	"fmt"

	"soditapp/backend/model"

	"github.com/google/uuid"
)

// ClientRepository handles database operations for clients.
type ClientRepository struct {
	db *sql.DB
}

// NewClientRepository creates a new ClientRepository.
func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

const DefaultPriorityGap = 1024

// Create inserts a new client and returns it. New clients are appended to the end
// of their zone's ordering using a gap-based priority (max + defaultPriorityGap).
func (r *ClientRepository) Create(ctx context.Context, name, address, zoneRefUUID string) (model.Client, error) {
	var c model.Client
	var addr sql.NullString
	var unregDate sql.NullString

	max, err := r.MaxPriorityInZone(ctx, zoneRefUUID)
	if err != nil {
		return c, fmt.Errorf("getting max priority: %w", err)
	}
	newPriority := max + DefaultPriorityGap

	refUUID, err := uuid.NewV7()
	if err != nil {
		return c, fmt.Errorf("generating UUID: %w", err)
	}

	err = r.db.QueryRowContext(ctx,
		`INSERT INTO client (name, address, zone_id, priority_order, reference_uuid)
		 VALUES (?, ?, ?, ?, ?)
		 RETURNING id, name, address, zone_id, create_date, unregister_date, priority_order, reference_uuid`,
		name, address, zoneRefUUID, newPriority, refUUID.String(),
	).Scan(&c.ID, &c.Name, &addr, &c.ZoneID, &c.CreateDate, &unregDate, &c.PriorityOrder, &c.ReferenceUUID)
	if err != nil {
		return c, fmt.Errorf("creating client: %w", err)
	}

	if addr.Valid {
		c.Address = addr.String
	}
	if unregDate.Valid {
		c.UnregisterDate = &unregDate.String
	}

	return c, nil
}

// Upsert inserts a new client or updates an existing one (by name).
// If inserting, assigns a priority_order as append.
func (r *ClientRepository) Upsert(ctx context.Context, name, address, zoneRefUUID string) (model.Client, error) {
	var c model.Client
	var addr sql.NullString
	var unregDate sql.NullString

	max, err := r.MaxPriorityInZone(ctx, zoneRefUUID)
	if err != nil {
		return c, fmt.Errorf("getting max priority: %w", err)
	}
	newPriority := max + DefaultPriorityGap

	refUUID, err := uuid.NewV7()
	if err != nil {
		return c, fmt.Errorf("generating UUID: %w", err)
	}

	err = r.db.QueryRowContext(ctx,
		`INSERT INTO client (name, address, zone_id, priority_order, reference_uuid)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT (name) DO UPDATE SET
		   address = EXCLUDED.address,
		   zone_id = EXCLUDED.zone_id
		 RETURNING id, name, address, zone_id, create_date, unregister_date, priority_order, reference_uuid`,
		name, address, zoneRefUUID, newPriority, refUUID.String(),
	).Scan(&c.ID, &c.Name, &addr, &c.ZoneID, &c.CreateDate, &unregDate, &c.PriorityOrder, &c.ReferenceUUID)
	if err != nil {
		return c, fmt.Errorf("upserting client: %w", err)
	}

	if addr.Valid {
		c.Address = addr.String
	}
	if unregDate.Valid {
		c.UnregisterDate = &unregDate.String
	}

	return c, nil
}

// GetByID returns a client by its ID, including priority_order.
func (r *ClientRepository) GetByID(ctx context.Context, id int) (model.Client, error) {
	var c model.Client
	var addr sql.NullString
	var unregDate sql.NullString

	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, address, zone_id, create_date, unregister_date, priority_order, reference_uuid FROM client WHERE id = ?",
		id,
	).Scan(&c.ID, &c.Name, &addr, &c.ZoneID, &c.CreateDate, &unregDate, &c.PriorityOrder, &c.ReferenceUUID)
	if err != nil {
		return c, fmt.Errorf("getting client %d: %w", id, err)
	}

	if addr.Valid {
		c.Address = addr.String
	}
	if unregDate.Valid {
		c.UnregisterDate = &unregDate.String
	}

	return c, nil
}

// MaxPriorityInZone returns the current maximum priority_order for a zone (0 if none).
func (r *ClientRepository) MaxPriorityInZone(ctx context.Context, zoneRefUUID string) (int, error) {
	var max sql.NullInt64
	err := r.db.QueryRowContext(ctx, "SELECT COALESCE(MAX(priority_order), 0) FROM client WHERE zone_id = ?", zoneRefUUID).Scan(&max)
	if err != nil {
		return 0, fmt.Errorf("getting max priority for zone %s: %w", zoneRefUUID, err)
	}
	if max.Valid {
		return int(max.Int64), nil
	}
	return 0, nil
}

// UpdatePriority sets the priority_order for a client.
func (r *ClientRepository) UpdatePriority(ctx context.Context, id int, priority int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE client SET priority_order = ? WHERE id = ?", priority, id)
	if err != nil {
		return fmt.Errorf("updating priority for client %d: %w", id, err)
	}
	return nil
}

// ReindexZone rebuilds priority_order values for all clients in the zone using the provided step.
// It assigns priority = (index+1) * step in current order.
func (r *ClientRepository) ReindexZone(ctx context.Context, zoneRefUUID string, step int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin reindex tx: %w", err)
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, "SELECT id FROM client WHERE zone_id = ? AND unregister_date IS NULL ORDER BY priority_order ASC, id ASC", zoneRefUUID)
	if err != nil {
		return fmt.Errorf("selecting clients for reindex: %w", err)
	}
	defer rows.Close()

	ids := make([]int, 0)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("scanning id for reindex: %w", err)
		}
		ids = append(ids, id)
	}

	for i, id := range ids {
		newPriority := (i + 1) * step
		if _, err := tx.ExecContext(ctx, "UPDATE client SET priority_order = ? WHERE id = ?", newPriority, id); err != nil {
			return fmt.Errorf("updating priority during reindex: %w", err)
		}
	}

	return tx.Commit()
}

// GetByZone returns active (not unregistered) clients in a zone, ordered by priority_order.
func (r *ClientRepository) GetByZone(ctx context.Context, zoneRefUUID string) ([]model.Client, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, address, zone_id, create_date, unregister_date, priority_order, reference_uuid
		 FROM client
		 WHERE zone_id = ? AND unregister_date IS NULL
		 ORDER BY priority_order ASC`,
		zoneRefUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing clients for zone %s: %w", zoneRefUUID, err)
	}
	defer rows.Close()

	return scanClients(rows)
}

// GetAll returns all active (not unregistered) clients ordered by priority_order.
func (r *ClientRepository) GetAll(ctx context.Context) ([]model.Client, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, address, zone_id, create_date, unregister_date, priority_order, reference_uuid
		 FROM client
		 WHERE unregister_date IS NULL
		 ORDER BY priority_order ASC`,
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
		"UPDATE client SET name = ?, address = ?, zone_id = ? WHERE id = ?",
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
		"UPDATE client SET unregister_date = date('now') WHERE id = ?",
		id,
	)
	if err != nil {
		return fmt.Errorf("unregistering client %d: %w", id, err)
	}
	return nil
}

// GetInactive returns all inactive (unregistered) clients ordered by priority_order.
func (r *ClientRepository) GetInactive(ctx context.Context) ([]model.Client, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, address, zone_id, create_date, unregister_date, priority_order, reference_uuid
		 FROM client
		 WHERE unregister_date IS NOT NULL
		 ORDER BY priority_order ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing inactive clients: %w", err)
	}
	defer rows.Close()

	return scanClients(rows)
}

// Restore reactivates a client by setting unregister_date to NULL.
func (r *ClientRepository) Restore(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE client SET unregister_date = NULL WHERE id = ?",
		id,
	)
	if err != nil {
		return fmt.Errorf("restoring client %d: %w", id, err)
	}
	return nil
}

// Delete removes a client by its ID.
func (r *ClientRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM client WHERE id = ?",
		id,
	)
	if err != nil {
		return fmt.Errorf("deleting client %d: %w", id, err)
	}
	return nil
}

// BulkUpdateZone moves a set of clients to a new zone in a single transaction.
func (r *ClientRepository) BulkUpdateZone(ctx context.Context, clientRefUUIDs []string, targetZoneRefUUID string) error {
	if len(clientRefUUIDs) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning bulk update zone tx: %w", err)
	}
	defer tx.Rollback()

	placeholders := make([]string, len(clientRefUUIDs))
	args := make([]interface{}, len(clientRefUUIDs)+1)
	args[0] = targetZoneRefUUID
	for i, refUUID := range clientRefUUIDs {
		placeholders[i] = "?"
		args[i+1] = refUUID
	}
	query := fmt.Sprintf(
		"UPDATE client SET zone_id = ? WHERE reference_uuid IN (%s)",
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
	clients := make([]model.Client, 0)
	for rows.Next() {
		var c model.Client
		var addr sql.NullString
		var unregDate sql.NullString

		if err := rows.Scan(&c.ID, &c.Name, &addr, &c.ZoneID, &c.CreateDate, &unregDate, &c.PriorityOrder, &c.ReferenceUUID); err != nil {
			return nil, fmt.Errorf("scanning client: %w", err)
		}

		if addr.Valid {
			c.Address = addr.String
		}
		if unregDate.Valid {
			c.UnregisterDate = &unregDate.String
		}

		clients = append(clients, c)
	}
	return clients, rows.Err()
}
