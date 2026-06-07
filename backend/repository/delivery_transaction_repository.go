package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// DeliveryTransactionRow holds joined delivery transaction data for reports.
type DeliveryTransactionRow struct {
	ID           int
	ClientID     int
	ClientName   string
	ProductID    int
	ProductName  string
	ChargedUnits int
	CreditUnits  int
	DeliveryDate time.Time
}

// DeliveryTransactionRepository handles database operations for delivery transactions.
type DeliveryTransactionRepository struct {
	db *sql.DB
}

// NewDeliveryTransactionRepository creates a new DeliveryTransactionRepository.
func NewDeliveryTransactionRepository(db *sql.DB) *DeliveryTransactionRepository {
	return &DeliveryTransactionRepository{db: db}
}

// InsertTx inserts a delivery transaction record within an existing DB transaction.
func (r *DeliveryTransactionRepository) InsertTx(ctx context.Context, tx *sql.Tx, clientID, productID, chargedUnits, creditUnits int, deliveryDate time.Time) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO delivery_transaction (client_id, product_id, charged_units, credit_units, delivery_date)
		 VALUES ($1, $2, $3, $4, $5)`,
		clientID, productID, chargedUnits, creditUnits, deliveryDate,
	)
	if err != nil {
		return fmt.Errorf("inserting delivery transaction for client %d, product %d: %w", clientID, productID, err)
	}
	return nil
}

// GetByClient returns delivery transactions for a client within a date range.
func (r *DeliveryTransactionRepository) GetByClient(ctx context.Context, clientID int, startDate, endDate string) ([]DeliveryTransactionRow, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT dt.id, dt.client_id, c.name, dt.product_id, p.name,
		        dt.charged_units, dt.credit_units, dt.delivery_date
		 FROM delivery_transaction dt
		 JOIN client c ON c.id = dt.client_id
		 JOIN product p ON p.id = dt.product_id
		 WHERE dt.client_id = $1
		   AND dt.delivery_date BETWEEN $2 AND $3
		 ORDER BY dt.delivery_date DESC, p.name`,
		clientID, startDate, endDate,
	)
	if err != nil {
		return nil, fmt.Errorf("getting transactions for client %d: %w", clientID, err)
	}
	defer rows.Close()

	return scanTransactionRows(rows)
}

// GetByZone returns delivery transactions for all clients in a zone within a date range.
func (r *DeliveryTransactionRepository) GetByZone(ctx context.Context, zoneID int, startDate, endDate string) ([]DeliveryTransactionRow, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT dt.id, dt.client_id, c.name, dt.product_id, p.name,
		        dt.charged_units, dt.credit_units, dt.delivery_date
		 FROM delivery_transaction dt
		 JOIN client c ON c.id = dt.client_id
		 JOIN product p ON p.id = dt.product_id
		 WHERE c.zone_id = $1
		   AND dt.delivery_date BETWEEN $2 AND $3
		 ORDER BY dt.delivery_date DESC, c.name, p.name`,
		zoneID, startDate, endDate,
	)
	if err != nil {
		return nil, fmt.Errorf("getting transactions for zone %d: %w", zoneID, err)
	}
	defer rows.Close()

	return scanTransactionRows(rows)
}

func scanTransactionRows(rows *sql.Rows) ([]DeliveryTransactionRow, error) {
	var results []DeliveryTransactionRow
	for rows.Next() {
		var row DeliveryTransactionRow
		if err := rows.Scan(
			&row.ID, &row.ClientID, &row.ClientName,
			&row.ProductID, &row.ProductName,
			&row.ChargedUnits, &row.CreditUnits, &row.DeliveryDate,
		); err != nil {
			return nil, fmt.Errorf("scanning delivery transaction row: %w", err)
		}
		results = append(results, row)
	}
	return results, rows.Err()
}
