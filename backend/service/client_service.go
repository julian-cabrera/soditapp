package service

import (
	"context"
	"fmt"

	"soditapp/backend/model"
	"soditapp/backend/repository"
)

// ClientService exposes client operations as Wails-bound methods.
type ClientService struct {
	ctx  context.Context
	repo *repository.ClientRepository
}

// Init sets the runtime dependencies after Wails startup.
func (s *ClientService) Init(ctx context.Context, repo *repository.ClientRepository) {
	s.ctx = ctx
	s.repo = repo
}

func (s *ClientService) CreateClient(name, address string, zoneID int) (model.Client, error) {
	if name == "" {
		return model.Client{}, fmt.Errorf("client name cannot be empty")
	}
	return s.repo.Create(s.ctx, name, address, zoneID)
}

func (s *ClientService) ListClients() ([]model.Client, error) {
	return s.repo.GetAll(s.ctx)
}

func (s *ClientService) ListClientsByZone(zoneID int) ([]model.Client, error) {
	return s.repo.GetByZone(s.ctx, zoneID)
}

func (s *ClientService) UpdateClient(id int, name, address string, zoneID int) error {
	if name == "" {
		return fmt.Errorf("client name cannot be empty")
	}
	return s.repo.Update(s.ctx, model.Client{
		ID:      id,
		Name:    name,
		Address: address,
		ZoneID:  zoneID,
	})
}

func (s *ClientService) UnregisterClient(id int) error {
	return s.repo.Unregister(s.ctx, id)
}
