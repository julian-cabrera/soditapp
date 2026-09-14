package service

import (
	"context"
	"fmt"

	"soditapp/backend/dto"
	"soditapp/backend/repository"
)

// ReportService provides delivery history and zone report queries.
type ReportService struct {
	ctx            context.Context
	deliveryTxRepo *repository.DeliveryTransactionRepository
}

// Init wires runtime dependencies into the service.
func (s *ReportService) Init(ctx context.Context, deliveryTxRepo *repository.DeliveryTransactionRepository) {
	s.ctx = ctx
	s.deliveryTxRepo = deliveryTxRepo
}

// Reinit updates runtime dependencies (called after database import/reopen).
func (s *ReportService) Reinit(deliveryTxRepo *repository.DeliveryTransactionRepository) {
	s.deliveryTxRepo = deliveryTxRepo
}

// GetClientHistory returns the delivery transactions for a client within the given date range.
// startDate and endDate must be in "YYYY-MM-DD" format.
func (s *ReportService) GetClientHistory(clientRefUUID string, startDate, endDate string) ([]dto.DeliveryHistoryEntryDTO, error) {
	rows, err := s.deliveryTxRepo.GetByClient(s.ctx, clientRefUUID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("GetClientHistory: %w", err)
	}

	result := make([]dto.DeliveryHistoryEntryDTO, len(rows))
	for i, row := range rows {
		result[i] = dto.DeliveryHistoryEntryDTO{
			ClientRefUUID: row.ClientRefUUID,
			ClientName:    row.ClientName,
			ProductRefUUID: row.ProductRefUUID,
			ProductName:   row.ProductName,
			ChargedUnits:  row.ChargedUnits,
			CreditUnits:   row.CreditUnits,
			DeliveryDate:  row.DeliveryDate,
		}
	}
	return result, nil
}

// GetZoneReport returns the delivery transactions for all clients in a zone within the given date range.
// startDate and endDate must be in "YYYY-MM-DD" format.
func (s *ReportService) GetZoneReport(zoneRefUUID string, startDate, endDate string) ([]dto.ZoneReportEntryDTO, error) {
	rows, err := s.deliveryTxRepo.GetByZone(s.ctx, zoneRefUUID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("GetZoneReport: %w", err)
	}

	result := make([]dto.ZoneReportEntryDTO, len(rows))
	for i, row := range rows {
		result[i] = dto.ZoneReportEntryDTO{
			ClientRefUUID: row.ClientRefUUID,
			ClientName:    row.ClientName,
			ProductRefUUID: row.ProductRefUUID,
			ProductName:   row.ProductName,
			ChargedUnits:  row.ChargedUnits,
			CreditUnits:   row.CreditUnits,
			DeliveryDate:  row.DeliveryDate,
		}
	}
	return result, nil
}
