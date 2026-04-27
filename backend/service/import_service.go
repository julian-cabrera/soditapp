package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"soditapp/backend/repository"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ImportResult summarizes the outcome of a client import operation.
type ImportResult struct {
	Inserted int      `json:"inserted"`
	Skipped  int      `json:"skipped"`
	Errors   []string `json:"errors"`
}

// ImportService exposes CSV import operations as Wails-bound methods.
type ImportService struct {
	ctx        context.Context
	clientRepo *repository.ClientRepository
	zoneRepo   *repository.ZoneRepository
}

// Init sets the runtime dependencies after Wails startup.
func (s *ImportService) Init(ctx context.Context, clientRepo *repository.ClientRepository, zoneRepo *repository.ZoneRepository) {
	s.ctx = ctx
	s.clientRepo = clientRepo
	s.zoneRepo = zoneRepo
}

// ImportClients opens a file dialog, parses the CSV, and upserts clients.
// Expected columns: nombre (required), direccion (optional), zona (must match existing zone name).
// If a client name already exists, it will be updated with the new address and zone.
// If a client name does not exist, a new client will be created.
func (s *ImportService) ImportClients() (ImportResult, error) {
	path, err := runtime.OpenFileDialog(s.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{
			{DisplayName: "CSV", Pattern: "*.csv"},
		},
	})
	if err != nil {
		return ImportResult{}, fmt.Errorf("opening file dialog: %w", err)
	}
	if path == "" {
		return ImportResult{}, nil // cancelled
	}

	f, err := os.Open(path)
	if err != nil {
		return ImportResult{}, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	// Build zone name -> ID map (case-insensitive)
	zones, err := s.zoneRepo.GetAll(s.ctx)
	if err != nil {
		return ImportResult{}, fmt.Errorf("fetching zones: %w", err)
	}
	zoneIDs := make(map[string]int, len(zones))
	for _, z := range zones {
		zoneIDs[strings.ToLower(z.Name)] = z.ID
	}

	reader := csv.NewReader(f)
	header, err := reader.Read()
	if err != nil {
		return ImportResult{}, fmt.Errorf("reading CSV header: %w", err)
	}

	// Find column indices
	colIndex := make(map[string]int)
	for i, col := range header {
		colIndex[strings.ToLower(strings.TrimSpace(col))] = i
	}
	nombreIdx, hasNombre := colIndex["nombre"]
	if !hasNombre {
		return ImportResult{}, fmt.Errorf("CSV must have a 'nombre' column")
	}
	zonaIdx, hasZona := colIndex["zona"]
	if !hasZona {
		return ImportResult{}, fmt.Errorf("CSV must have a 'zona' column")
	}
	direccionIdx, hasDireccion := colIndex["direccion"]

	var result ImportResult
	rowNum := 1
	for {
		row, err := reader.Read()
		if err != nil {
			break // EOF or error
		}
		rowNum++

		if nombreIdx >= len(row) {
			result.Errors = append(result.Errors, fmt.Sprintf("fila %d: columna 'nombre' faltante", rowNum))
			continue
		}
		nombre := strings.TrimSpace(row[nombreIdx])
		if nombre == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("fila %d: nombre vacío", rowNum))
			continue
		}

		if zonaIdx >= len(row) {
			result.Errors = append(result.Errors, fmt.Sprintf("fila %d: columna 'zona' faltante", rowNum))
			continue
		}
		zonaNombre := strings.TrimSpace(row[zonaIdx])
		zoneID, ok := zoneIDs[strings.ToLower(zonaNombre)]
		if !ok {
			result.Errors = append(result.Errors, fmt.Sprintf("fila %d: zona '%s' no existe", rowNum, zonaNombre))
			continue
		}

		direccion := ""
		if hasDireccion && direccionIdx < len(row) {
			direccion = strings.TrimSpace(row[direccionIdx])
		}

		if _, err := s.clientRepo.Upsert(s.ctx, nombre, direccion, zoneID); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("fila %d: %s", rowNum, err.Error()))
			continue
		}

		result.Inserted++
	}

	return result, nil
}
