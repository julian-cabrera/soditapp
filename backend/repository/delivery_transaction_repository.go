package repository

import (
	"context"
	"database/sql"
	"fmt"
)

// DeliveryTransactionRow holds joined delivery transaction data for reports.
type DeliveryTransactionRow struct {
	ID            int
	ClientRefUUID string
	ClientName    string
	ProductRefUUID string
	ProductName   string
	ChargedUnits  int
	CreditUnits   int
	DeliveryDate  string
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
func (r *DeliveryTransactionRepository) InsertTx(ctx context.Context, tx *sql.Tx, clientRefUUID, productRefUUID string, chargedUnits, creditUnits int, deliveryDate string) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO delivery_transaction (client_id, product_id, charged_units, credit_units, delivery_date)
		 VALUES (?, ?, ?, ?, ?)`,
		clientRefUUID, productRefUUID, chargedUnits, creditUnits, deliveryDate,
	)
	if err != nil {
		return fmt.Errorf("inserting delivery transaction for client %s, product %s: %w", clientRefUUID, productRefUUID, err)
	}
	return nil
}

// GetByClient returns delivery transactions for a client within a date range.
func (r *DeliveryTransactionRepository) GetByClient(ctx context.Context, clientRefUUID string, startDate, endDate string) ([]DeliveryTransactionRow, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT dt.id, dt.client_id, c.name, dt.product_id, p.name,
		        dt.charged_units, dt.credit_units, dt.delivery_date
		 FROM delivery_transaction dt
		 JOIN client c ON c.reference_uuid = dt.client_id
		 JOIN product p ON p.reference_uuid = dt.product_id
		 WHERE dt.client_id = ?
		   AND dt.delivery_date BETWEEN ? AND ?
		 ORDER BY dt.delivery_date DESC, p.name`,
		clientRefUUID, startDate, endDate,
	)
	if err != nil {
		return nil, fmt.Errorf("getting transactions for client %s: %w", clientRefUUID, err)
	}
	defer rows.Close()

	return scanTransactionRows(rows)
}

// GetByZone returns delivery transactions for all clients in a zone within a date range.
func (r *DeliveryTransactionRepository) GetByZone(ctx context.Context, zoneRefUUID string, startDate, endDate string) ([]DeliveryTransactionRow, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT dt.id, dt.client_id, c.name, dt.product_id, p.name,
		        dt.charged_units, dt.credit_units, dt.delivery_date
		 FROM delivery_transaction dt
		 JOIN client c ON c.reference_uuid = dt.client_id
		 JOIN product p ON p.reference_uuid = dt.product_id
		 WHERE c.zone_id = ?
		   AND dt.delivery_date BETWEEN ? AND ?
		 ORDER BY dt.delivery_date DESC, c.name, p.name`,
		zoneRefUUID, startDate, endDate,
	)
	if err != nil {
		return nil, fmt.Errorf("getting transactions for zone %s: %w", zoneRefUUID, err)
	}
	defer rows.Close()

	return scanTransactionRows(rows)
}

func scanTransactionRows(rows *sql.Rows) ([]DeliveryTransactionRow, error) {
	var results []DeliveryTransactionRow
	for rows.Next() {
		var row DeliveryTransactionRow
		if err := rows.Scan(
			&row.ID, &row.ClientRefUUID, &row.ClientName,
			&row.ProductRefUUID, &row.ProductName,
			&row.ChargedUnits, &row.CreditUnits, &row.DeliveryDate,
		); err != nil {
			return nil, fmt.Errorf("scanning delivery transaction row: %w", err)
		}
		results = append(results, row)
	}
	return results, rows.Err()
}
