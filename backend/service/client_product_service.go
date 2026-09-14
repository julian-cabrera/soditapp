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
	ClientRefUUID string `json:"client_ref_uuid"`
	ProductRefUUID string `json:"product_ref_uuid"`
	Delta         int    `json:"delta"`          // negative value = credit units (increases debt)
	ChargedUnits  int    `json:"charged_units"`  // units paid on the spot (no balance change)
	DeliveryDate  string `json:"delivery_date"`
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

// Reinit updates runtime dependencies (called after database import/reopen).
func (s *ClientProductService) Reinit(db *sql.DB, repo *repository.ClientProductRepository, zoneRepo *repository.ZoneRepository, txRepo *repository.DeliveryTransactionRepository) {
	s.db = db
	s.repo = repo
	s.zoneRepo = zoneRepo
	s.txRepo = txRepo
}

func (s *ClientProductService) EnsureClientProduct(clientRefUUID, productRefUUID string) (model.ClientProduct, error) {
	return s.repo.Ensure(s.ctx, clientRefUUID, productRefUUID)
}

func (s *ClientProductService) GetClientBalances(clientRefUUID string) ([]model.ClientProduct, error) {
	return s.repo.GetByClient(s.ctx, clientRefUUID)
}

func (s *ClientProductService) AdjustBalance(clientRefUUID, productRefUUID string, delta int) error {
	return s.repo.AdjustBalance(s.ctx, clientRefUUID, productRefUUID, delta)
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
			if err := s.repo.AdjustBalanceTx(s.ctx, tx, adj.ClientRefUUID, adj.ProductRefUUID, adj.Delta); err != nil {
				return fmt.Errorf("adjusting balance for client %s, product %s: %w", adj.ClientRefUUID, adj.ProductRefUUID, err)
			}
		}

		creditUnits := -adj.Delta // delta is stored as negative; invert to get positive credit count
		if creditUnits < 0 {
			creditUnits = 0
		}
		if err := s.txRepo.InsertTx(s.ctx, tx, adj.ClientRefUUID, adj.ProductRefUUID, adj.ChargedUnits, creditUnits, adj.DeliveryDate); err != nil {
			return fmt.Errorf("recording transaction for client %s, product %s: %w", adj.ClientRefUUID, adj.ProductRefUUID, err)
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

	rows, err := s.repo.GetAllByZone(s.ctx, zone.ReferenceUUID)
	if err != nil {
		return nil, fmt.Errorf("fetching clients for zone %d: %w", zoneID, err)
	}

	return dto.BuildDeliverySheets(zone.ID, zone.Name, rows, dto.DefaultEntriesPerPage), nil
}
