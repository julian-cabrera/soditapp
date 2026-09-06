package repository

import (
	"context"
	"database/sql"
	"fmt"

	"soditapp/backend/model"
)

// NegativeBalanceRow holds the joined data needed to build delivery sheets.
type NegativeBalanceRow struct {
	ClientID      int    `json:"client_id"`
	ClientName    string `json:"client_name"`
	ClientAddress string `json:"client_address"`
	Balance       int    `json:"balance"`
	ProductID     int    `json:"product_id"`
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
func (r *ClientProductRepository) Ensure(ctx context.Context, clientID, productID int) (model.ClientProduct, error) {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO client_product (client_id, product_id)
		 VALUES ($1, $2)
		 ON CONFLICT (client_id, product_id) DO NOTHING`,
		clientID, productID,
	)
	if err != nil {
		return model.ClientProduct{}, fmt.Errorf("ensuring client_product (%d, %d): %w", clientID, productID, err)
	}

	var cp model.ClientProduct
	err = r.db.QueryRowContext(ctx,
		`SELECT id, client_id, product_id, balance
		 FROM client_product
		 WHERE client_id = $1 AND product_id = $2`,
		clientID, productID,
	).Scan(&cp.ID, &cp.ClientID, &cp.ProductID, &cp.Balance)
	if err != nil {
		return cp, fmt.Errorf("fetching client_product (%d, %d): %w", clientID, productID, err)
	}

	return cp, nil
}

// GetByClient returns all client-product records for a given client.
func (r *ClientProductRepository) GetByClient(ctx context.Context, clientID int) ([]model.ClientProduct, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, client_id, product_id, balance
		 FROM client_product
		 WHERE client_id = $1
		 ORDER BY product_id`,
		clientID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing client_products for client %d: %w", clientID, err)
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
func (r *ClientProductRepository) AdjustBalance(ctx context.Context, clientID, productID, delta int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE client_product
		 SET balance = balance + $1
		 WHERE client_id = $2 AND product_id = $3`,
		delta, clientID, productID,
	)
	if err != nil {
		return fmt.Errorf("adjusting balance for (%d, %d): %w", clientID, productID, err)
	}
	return nil
}

// AdjustBalanceTx adds delta to the balance within a transaction.
func (r *ClientProductRepository) AdjustBalanceTx(ctx context.Context, tx *sql.Tx, clientID, productID, delta int) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE client_product
		 SET balance = balance + $1
		 WHERE client_id = $2 AND product_id = $3`,
		delta, clientID, productID,
	)
	if err != nil {
		return fmt.Errorf("adjusting balance in tx for (%d, %d): %w", clientID, productID, err)
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
		 JOIN client c ON c.id = cp.client_id
		 JOIN product p ON p.id = cp.product_id
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
func (r *ClientProductRepository) GetNegativeBalancesByZone(ctx context.Context, zoneID int) ([]NegativeBalanceRow, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT c.id, c.name, COALESCE(c.address, ''), cp.balance, p.id, p.name
		 FROM client_product cp
		 JOIN client c ON c.id = cp.client_id
		 JOIN product p ON p.id = cp.product_id
		 WHERE c.zone_id = $1
		   AND c.unregister_date IS NULL
		   AND cp.balance < 0
		 ORDER BY c.name, p.name`,
		zoneID,
	)
	if err != nil {
		return nil, fmt.Errorf("getting negative balances for zone %d: %w", zoneID, err)
	}
	defer rows.Close()

	var results []NegativeBalanceRow
	for rows.Next() {
		var r NegativeBalanceRow
		if err := rows.Scan(&r.ClientID, &r.ClientName, &r.ClientAddress, &r.Balance, &r.ProductID, &r.ProductName); err != nil {
			return nil, fmt.Errorf("scanning negative balance row: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

// GetAllByZone returns all active clients in a zone with their product balances,
// joined with client and product details for delivery sheet building.
func (r *ClientProductRepository) GetAllByZone(ctx context.Context, zoneID int) ([]NegativeBalanceRow, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT c.id, c.name, COALESCE(c.address, ''), COALESCE(cp.balance, 0), p.id, p.name
		 FROM client c
		 CROSS JOIN product p
		 LEFT JOIN client_product cp ON cp.client_id = c.id AND cp.product_id = p.id
		 WHERE c.zone_id = $1
		   AND c.unregister_date IS NULL
		   AND p.active = TRUE
		 ORDER BY c.priority_order ASC, p.name`,
		zoneID,
	)
	if err != nil {
		return nil, fmt.Errorf("getting clients for zone %d: %w", zoneID, err)
	}
	defer rows.Close()

	var results []NegativeBalanceRow
	for rows.Next() {
		var r NegativeBalanceRow
		if err := rows.Scan(&r.ClientID, &r.ClientName, &r.ClientAddress, &r.Balance, &r.ProductID, &r.ProductName); err != nil {
			return nil, fmt.Errorf("scanning negative balance row: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}
