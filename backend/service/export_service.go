package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"soditapp/backend/repository"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ExportService exposes CSV export operations as Wails-bound methods.
type ExportService struct {
	ctx               context.Context
	zoneRepo          *repository.ZoneRepository
	clientRepo        *repository.ClientRepository
	clientProductRepo *repository.ClientProductRepository
}

// Init sets the runtime dependencies after Wails startup.
func (s *ExportService) Init(
	ctx context.Context,
	zoneRepo *repository.ZoneRepository,
	clientRepo *repository.ClientRepository,
	clientProductRepo *repository.ClientProductRepository,
) {
	s.ctx = ctx
	s.zoneRepo = zoneRepo
	s.clientRepo = clientRepo
	s.clientProductRepo = clientProductRepo
}

// ExportZones opens a save dialog and writes zones to CSV.
func (s *ExportService) ExportZones() error {
	path, err := runtime.SaveFileDialog(s.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "zonas.csv",
		Filters: []runtime.FileFilter{
			{DisplayName: "CSV", Pattern: "*.csv"},
		},
	})
	if err != nil {
		return fmt.Errorf("opening save dialog: %w", err)
	}
	if path == "" {
		return nil // cancelled
	}

	zones, err := s.zoneRepo.GetAll(s.ctx)
	if err != nil {
		return fmt.Errorf("fetching zones: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.Write([]string{"id", "nombre"}); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}
	for _, z := range zones {
		if err := w.Write([]string{strconv.Itoa(z.ID), z.Name}); err != nil {
			return fmt.Errorf("writing row: %w", err)
		}
	}
	w.Flush()
	return w.Error()
}

// ExportClients opens a save dialog and writes clients to CSV.
func (s *ExportService) ExportClients() error {
	path, err := runtime.SaveFileDialog(s.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "clientes.csv",
		Filters: []runtime.FileFilter{
			{DisplayName: "CSV", Pattern: "*.csv"},
		},
	})
	if err != nil {
		return fmt.Errorf("opening save dialog: %w", err)
	}
	if path == "" {
		return nil
	}

	zones, err := s.zoneRepo.GetAll(s.ctx)
	if err != nil {
		return fmt.Errorf("fetching zones: %w", err)
	}
	zoneNames := make(map[int]string, len(zones))
	for _, z := range zones {
		zoneNames[z.ID] = z.Name
	}

	clients, err := s.clientRepo.GetAll(s.ctx)
	if err != nil {
		return fmt.Errorf("fetching clients: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.Write([]string{"id", "nombre", "direccion", "zona"}); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}
	for _, c := range clients {
		if err := w.Write([]string{
			strconv.Itoa(c.ID),
			c.Name,
			c.Address,
			zoneNames[c.ZoneID],
		}); err != nil {
			return fmt.Errorf("writing row: %w", err)
		}
	}
	w.Flush()
	return w.Error()
}

// ExportBalances opens a save dialog and writes client-product balances to CSV.
func (s *ExportService) ExportBalances() error {
	path, err := runtime.SaveFileDialog(s.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "saldos.csv",
		Filters: []runtime.FileFilter{
			{DisplayName: "CSV", Pattern: "*.csv"},
		},
	})
	if err != nil {
		return fmt.Errorf("opening save dialog: %w", err)
	}
	if path == "" {
		return nil
	}

	balances, err := s.clientProductRepo.GetAllBalances(s.ctx)
	if err != nil {
		return fmt.Errorf("fetching balances: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.Write([]string{"cliente", "producto", "saldo"}); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}
	for _, b := range balances {
		if err := w.Write([]string{b.ClientName, b.ProductName, strconv.Itoa(b.Balance)}); err != nil {
			return fmt.Errorf("writing row: %w", err)
		}
	}
	w.Flush()
	return w.Error()
}
