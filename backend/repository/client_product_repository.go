package repository

import (
	"context"
	"database/sql"
	"fmt"

	"soditapp/backend/model"
)

// NegativeBalanceRow holds the joined data needed to build delivery sheets.
type NegativeBalanceRow struct {
	ClientRefUUID string `json:"client_ref_uuid"`
	ClientName    string `json:"client_name"`
	ClientAddress string `json:"client_address"`
	Balance       int    `json:"balance"`
	ProductRefUUID string `json:"product_ref_uuid"`
	ProductName   string `json:"product_name"`
}

// ClientProductRepository handles database operations for client-product balances.
type ClientProductRepository struct {
	db *sql.DB
}

// NewClientProductRepository creates a new ClientProductRepository.
func NewClientProductRepository(db *sql.DB) *ClientProductRepository {
	return &ClientProductRepository{db: db}
}

// Ensure creates a client-product record if it doesn't exist, then returns it.
func (r *ClientProductRepository) Ensure(ctx context.Context, clientRefUUID, productRefUUID string) (model.ClientProduct, error) {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO client_product (client_id, product_id)
		 VALUES (?, ?)
		 ON CONFLICT (client_id, product_id) DO NOTHING`,
		clientRefUUID, productRefUUID,
	)
	if err != nil {
		return model.ClientProduct{}, fmt.Errorf("ensuring client_product (%s, %s): %w", clientRefUUID, productRefUUID, err)
	}

	var cp model.ClientProduct
	err = r.db.QueryRowContext(ctx,
		`SELECT id, client_id, product_id, balance
		 FROM client_product
		 WHERE client_id = ? AND product_id = ?`,
		clientRefUUID, productRefUUID,
	).Scan(&cp.ID, &cp.ClientID, &cp.ProductID, &cp.Balance)
	if err != nil {
		return cp, fmt.Errorf("fetching client_product (%s, %s): %w", clientRefUUID, productRefUUID, err)
	}

	return cp, nil
}

// GetByClient returns all client-product records for a given client.
func (r *ClientProductRepository) GetByClient(ctx context.Context, clientRefUUID string) ([]model.ClientProduct, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, client_id, product_id, balance
		 FROM client_product
		 WHERE client_id = ?
		 ORDER BY product_id`,
		clientRefUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing client_products for client %s: %w", clientRefUUID, err)
	}
	defer rows.Close()

	var results []model.ClientProduct
	for rows.Next() {
		var cp model.ClientProduct
		if err := rows.Scan(&cp.ID, &cp.ClientID, &cp.ProductID, &cp.Balance); err != nil {
			return nil, fmt.Errorf("scanning client_product: %w", err)
		}
		results = append(results, cp)
	}
	return results, rows.Err()
}

// AdjustBalance adds delta to the balance of a client-product pair.
func (r *ClientProductRepository) AdjustBalance(ctx context.Context, clientRefUUID, productRefUUID string, delta int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE client_product
		 SET balance = balance + ?
		 WHERE client_id = ? AND product_id = ?`,
		delta, clientRefUUID, productRefUUID,
	)
	if err != nil {
		return fmt.Errorf("adjusting balance for (%s, %s): %w", clientRefUUID, productRefUUID, err)
	}
	return nil
}

// AdjustBalanceTx adds delta to the balance within a transaction.
func (r *ClientProductRepository) AdjustBalanceTx(ctx context.Context, tx *sql.Tx, clientRefUUID, productRefUUID string, delta int) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE client_product
		 SET balance = balance + ?
		 WHERE client_id = ? AND product_id = ?`,
		delta, clientRefUUID, productRefUUID,
	)
	if err != nil {
		return fmt.Errorf("adjusting balance in tx for (%s, %s): %w", clientRefUUID, productRefUUID, err)
	}
	return nil
}

// BalanceExportRow holds joined data for the balances CSV export.
type BalanceExportRow struct {
	ClientName  string
	ProductName string
	Balance     int
}

// GetAllBalances returns all active client-product balances joined with client and product names.
func (r *ClientProductRepository) GetAllBalances(ctx context.Context) ([]BalanceExportRow, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT c.name, p.name, cp.balance
		 FROM client_product cp
		 JOIN client c ON c.reference_uuid = cp.client_id
		 JOIN product p ON p.reference_uuid = cp.product_id
		 WHERE c.unregister_date IS NULL
		 ORDER BY c.name, p.name`,
	)
	if err != nil {
		return nil, fmt.Errorf("getting all balances: %w", err)
	}
	defer rows.Close()

	var results []BalanceExportRow
	for rows.Next() {
		var r BalanceExportRow
		if err := rows.Scan(&r.ClientName, &r.ProductName, &r.Balance); err != nil {
			return nil, fmt.Errorf("scanning balance export row: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

// GetNegativeBalancesByZone returns clients with negative balances in a zone,
// joined with client and product details for delivery sheet building.
func (r *ClientProductRepository) GetNegativeBalancesByZone(ctx context.Context, zoneRefUUID string) ([]NegativeBalanceRow, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT c.reference_uuid, c.name, COALESCE(c.address, ''), cp.balance, p.reference_uuid, p.name
		 FROM client_product cp
		 JOIN client c ON c.reference_uuid = cp.client_id
		 JOIN product p ON p.reference_uuid = cp.product_id
		 WHERE c.zone_id = ?
		   AND c.unregister_date IS NULL
		   AND cp.balance < 0
		 ORDER BY c.name, p.name`,
		zoneRefUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("getting negative balances for zone %s: %w", zoneRefUUID, err)
	}
	defer rows.Close()

	var results []NegativeBalanceRow
	for rows.Next() {
		var r NegativeBalanceRow
		if err := rows.Scan(&r.ClientRefUUID, &r.ClientName, &r.ClientAddress, &r.Balance, &r.ProductRefUUID, &r.ProductName); err != nil {
			return nil, fmt.Errorf("scanning negative balance row: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

// GetAllByZone returns all active clients in a zone with their product balances,
// joined with client and product details for delivery sheet building.
func (r *ClientProductRepository) GetAllByZone(ctx context.Context, zoneRefUUID string) ([]NegativeBalanceRow, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT c.reference_uuid, c.name, COALESCE(c.address, ''), COALESCE(cp.balance, 0), p.reference_uuid, p.name
		 FROM client c
		 CROSS JOIN product p
		 LEFT JOIN client_product cp ON cp.client_id = c.reference_uuid AND cp.product_id = p.reference_uuid
		 WHERE c.zone_id = ?
		   AND c.unregister_date IS NULL
		   AND p.active = 1
		 ORDER BY c.priority_order ASC, p.name`,
		zoneRefUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("getting clients for zone %s: %w", zoneRefUUID, err)
	}
	defer rows.Close()

	var results []NegativeBalanceRow
	for rows.Next() {
		var r NegativeBalanceRow
		if err := rows.Scan(&r.ClientRefUUID, &r.ClientName, &r.ClientAddress, &r.Balance, &r.ProductRefUUID, &r.ProductName); err != nil {
			return nil, fmt.Errorf("scanning negative balance row: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}
