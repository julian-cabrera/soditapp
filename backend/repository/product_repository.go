package repository

import (
	"context"
	"database/sql"
	"fmt"

	"soditapp/backend/model"

	"github.com/google/uuid"
)

// ProductRepository handles database operations for products.
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new ProductRepository.
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create inserts a new product and returns it.
func (r *ProductRepository) Create(ctx context.Context, name, typ string) (model.Product, error) {
	var p model.Product
	refUUID, err := uuid.NewV7()
	if err != nil {
		return p, fmt.Errorf("generating UUID: %w", err)
	}
	err = r.db.QueryRowContext(ctx,
		"INSERT INTO product (name, type, reference_uuid) VALUES (?, ?, ?) RETURNING id, name, type, active, reference_uuid",
		name, typ, refUUID.String(),
	).Scan(&p.ID, &p.Name, &p.Type, &p.Active, &p.ReferenceUUID)
	if err != nil {
		return p, fmt.Errorf("creating product: %w", err)
	}
	return p, nil
}

// GetByID returns a product by its ID.
func (r *ProductRepository) GetByID(ctx context.Context, id int) (model.Product, error) {
	var p model.Product
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, type, active FROM product WHERE id = ?",
		id,
	).Scan(&p.ID, &p.Name, &p.Type, &p.Active)
	if err != nil {
		return p, fmt.Errorf("getting product %d: %w", id, err)
	}
	return p, nil
}

// GetAll returns all products ordered by name.
func (r *ProductRepository) GetAll(ctx context.Context) ([]model.Product, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, type, active FROM product ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("listing products: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Type, &p.Active); err != nil {
			return nil, fmt.Errorf("scanning product: %w", err)
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// GetActive returns only active products ordered by name.
func (r *ProductRepository) GetActive(ctx context.Context) ([]model.Product, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, name, type, active FROM product WHERE active = 1 ORDER BY name",
	)
	if err != nil {
		return nil, fmt.Errorf("listing active products: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Type, &p.Active); err != nil {
			return nil, fmt.Errorf("scanning product: %w", err)
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// Update modifies a product's name.
func (r *ProductRepository) Update(ctx context.Context, id int, name string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE product SET name = ? WHERE id = ?",
		name, id,
	)
	if err != nil {
		return fmt.Errorf("updating product %d: %w", id, err)
	}
	return nil
}

// ToggleActive flips a product's active status.
func (r *ProductRepository) ToggleActive(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE product SET active = NOT active WHERE id = ?",
		id,
	)
	if err != nil {
		return fmt.Errorf("toggling product %d active status: %w", id, err)
	}
	return nil
}

// Delete removes a product by its ID.
func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM product WHERE id = ?",
		id,
	)
	if err != nil {
		return fmt.Errorf("deleting product %d: %w", id, err)
	}
	return nil
}
