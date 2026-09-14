package dto

import (
	"time"

	"soditapp/backend/repository"
)

// DefaultEntriesPerPage is the number of client entries that fit on one A4 page.
const DefaultEntriesPerPage = 25

// DeliveryEntryDTO represents a single client row on a delivery sheet.
type DeliveryEntryDTO struct {
	ClientRefUUID string `json:"client_ref_uuid"`
	ClientName    string `json:"client_name"`
	ClientAddress string `json:"client_address"`
	ClientBalance int    `json:"client_balance"`
	ProductRefUUID string `json:"product_ref_uuid"`
	ProductName   string `json:"product_name"`
	ChargedUnits  int    `json:"charged_units"`
	CreditUnits   int    `json:"credit_units"`
}

// DeliverySheetDTO represents one printed page of a delivery sheet for a zone.
type DeliverySheetDTO struct {
	ZoneID     int                `json:"zone_id"`
	ZoneName   string             `json:"zone_name"`
	CreatedAt  time.Time          `json:"created_at"`
	PageNumber int                `json:"page_number"`
	Entries    []DeliveryEntryDTO `json:"entries"`
}

// BuildDeliverySheets converts query rows into paginated delivery sheets.
// Each sheet holds at most entriesPerPage entries. ChargedUnits and CreditUnits are
// initialized to 0 — they are filled in by the user during delivery.
func BuildDeliverySheets(zoneID int, zoneName string, rows []repository.NegativeBalanceRow, entriesPerPage int) []DeliverySheetDTO {
	if entriesPerPage <= 0 {
		entriesPerPage = DefaultEntriesPerPage
	}

	entries := make([]DeliveryEntryDTO, len(rows))
	for i, row := range rows {
		entries[i] = DeliveryEntryDTO{
			ClientRefUUID: row.ClientRefUUID,
			ClientName:    row.ClientName,
			ClientAddress: row.ClientAddress,
			ClientBalance: row.Balance,
			ProductRefUUID: row.ProductRefUUID,
			ProductName:   row.ProductName,
		}
	}

	now := time.Now()
	var sheets []DeliverySheetDTO

	for i := 0; i < len(entries); i += entriesPerPage {
		end := i + entriesPerPage
		if end > len(entries) {
			end = len(entries)
		}

		sheets = append(sheets, DeliverySheetDTO{
			ZoneID:     zoneID,
			ZoneName:   zoneName,
			CreatedAt:  now,
			PageNumber: len(sheets) + 1,
			Entries:    entries[i:end],
		})
	}

	if len(sheets) == 0 {
		sheets = append(sheets, DeliverySheetDTO{
			ZoneID:     zoneID,
			ZoneName:   zoneName,
			CreatedAt:  now,
			PageNumber: 1,
			Entries:    []DeliveryEntryDTO{},
		})
	}

	return sheets
}
