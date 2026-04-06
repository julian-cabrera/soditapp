package service

import (
	"context"
	"database/sql"
	"fmt"

	"soditapp/backend/dto"
	"soditapp/backend/model"
	"soditapp/backend/repository"
)

// BalanceAdjustment represents a single balance change in a delivery confirmation.
type BalanceAdjustment struct {
	ClientID     int `json:"client_id"`
	ProductID    int `json:"product_id"`
	Delta        int `json:"delta"`          // negative value = credit units (increases debt)
	ChargedUnits int `json:"charged_units"`  // units paid on the spot (no balance change)
}

// ClientProductService exposes client-product balance operations as Wails-bound methods.
type ClientProductService struct {
	ctx      context.Context
	db       *sql.DB
	repo     *repository.ClientProductRepository
	zoneRepo *repository.ZoneRepository
	txRepo   *repository.DeliveryTransactionRepository
}

// Init sets the runtime dependencies after Wails startup.
func (s *ClientProductService) Init(ctx context.Context, db *sql.DB, repo *repository.ClientProductRepository, zoneRepo *repository.ZoneRepository, txRepo *repository.DeliveryTransactionRepository) {
	s.ctx = ctx
	s.db = db
	s.repo = repo
	s.zoneRepo = zoneRepo
	s.txRepo = txRepo
}

func (s *ClientProductService) EnsureClientProduct(clientID, productID int) (model.ClientProduct, error) {
	return s.repo.Ensure(s.ctx, clientID, productID)
}

func (s *ClientProductService) GetClientBalances(clientID int) ([]model.ClientProduct, error) {
	return s.repo.GetByClient(s.ctx, clientID)
}

func (s *ClientProductService) AdjustBalance(clientID, productID, delta int) error {
	return s.repo.AdjustBalance(s.ctx, clientID, productID, delta)
}

// ConfirmDelivery applies multiple balance adjustments in a single transaction.
func (s *ClientProductService) ConfirmDelivery(adjustments []BalanceAdjustment) error {
	if len(adjustments) == 0 {
		return nil
	}

	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning delivery transaction: %w", err)
	}
	defer tx.Rollback()

	for _, adj := range adjustments {
		if adj.Delta == 0 && adj.ChargedUnits == 0 {
			continue
		}
		if adj.Delta != 0 {
			if err := s.repo.AdjustBalanceTx(s.ctx, tx, adj.ClientID, adj.ProductID, adj.Delta); err != nil {
				return fmt.Errorf("adjusting balance for client %d, product %d: %w", adj.ClientID, adj.ProductID, err)
			}
		}
		creditUnits := -adj.Delta // delta is stored as negative; invert to get positive credit count
		if creditUnits < 0 {
			creditUnits = 0
		}
		if err := s.txRepo.InsertTx(s.ctx, tx, adj.ClientID, adj.ProductID, adj.ChargedUnits, creditUnits); err != nil {
			return fmt.Errorf("recording transaction for client %d, product %d: %w", adj.ClientID, adj.ProductID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing delivery transaction: %w", err)
	}
	return nil
}

// GetDeliverySheets builds paginated delivery sheets for a zone.
func (s *ClientProductService) GetDeliverySheets(zoneID int) ([]dto.DeliverySheetDTO, error) {
	zone, err := s.zoneRepo.GetByID(s.ctx, zoneID)
	if err != nil {
		return nil, fmt.Errorf("fetching zone %d: %w", zoneID, err)
	}

	rows, err := s.repo.GetNegativeBalancesByZone(s.ctx, zoneID)
	if err != nil {
		return nil, fmt.Errorf("fetching negative balances for zone %d: %w", zoneID, err)
	}

	return dto.BuildDeliverySheets(zone.ID, zone.Name, rows, dto.DefaultEntriesPerPage), nil
}
