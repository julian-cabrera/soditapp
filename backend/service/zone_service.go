package service

import (
	"context"
	"fmt"

	"soditapp/backend/model"
	"soditapp/backend/repository"
)

// ZoneService exposes zone operations as Wails-bound methods.
type ZoneService struct {
	ctx  context.Context
	repo *repository.ZoneRepository
}

// Init sets the runtime dependencies after Wails startup.
func (s *ZoneService) Init(ctx context.Context, repo *repository.ZoneRepository) {
	s.ctx = ctx
	s.repo = repo
}

func (s *ZoneService) CreateZone(name string) (model.Zone, error) {
	if name == "" {
		return model.Zone{}, fmt.Errorf("zone name cannot be empty")
	}
	return s.repo.Create(s.ctx, name)
}

func (s *ZoneService) ListZones() ([]model.Zone, error) {
	return s.repo.GetAll(s.ctx)
}

func (s *ZoneService) UpdateZone(id int, name string) error {
	if name == "" {
		return fmt.Errorf("zone name cannot be empty")
	}
	return s.repo.Update(s.ctx, id, name)
}

func (s *ZoneService) DeleteZone(id int) error {
	return s.repo.Delete(s.ctx, id)
}
