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

// Reinit updates runtime dependencies (called after database import/reopen).
func (s *ClientService) Reinit(repo *repository.ClientRepository) {
	s.repo = repo
}

func (s *ClientService) CreateClient(name, address, zoneRefUUID string) (model.Client, error) {
	if name == "" {
		return model.Client{}, fmt.Errorf("client name cannot be empty")
	}
	return s.repo.Create(s.ctx, name, address, zoneRefUUID)
}

func (s *ClientService) ListClients() ([]model.Client, error) {
	return s.repo.GetAll(s.ctx)
}

func (s *ClientService) ListClientsByZone(zoneRefUUID string) ([]model.Client, error) {
	return s.repo.GetByZone(s.ctx, zoneRefUUID)
}

func (s *ClientService) UpdateClient(id int, name, address, zoneRefUUID string) error {
	if name == "" {
		return fmt.Errorf("client name cannot be empty")
	}
	return s.repo.Update(s.ctx, model.Client{
		ID:      id,
		Name:    name,
		Address: address,
		ZoneID:  zoneRefUUID,
	})
}

func (s *ClientService) UnregisterClient(id int) error {
	return s.repo.Unregister(s.ctx, id)
}

func (s *ClientService) ListInactiveClients() ([]model.Client, error) {
	return s.repo.GetInactive(s.ctx)
}

func (s *ClientService) RestoreClient(id int) error {
	return s.repo.Restore(s.ctx, id)
}

func (s *ClientService) MoveClientsToZone(clientRefUUIDs []string, targetZoneRefUUID string) error {
	if len(clientRefUUIDs) == 0 {
		return fmt.Errorf("no clients selected")
	}
	if targetZoneRefUUID == "" {
		return fmt.Errorf("invalid target zone")
	}
	return s.repo.BulkUpdateZone(s.ctx, clientRefUUIDs, targetZoneRefUUID)
}

// ReorderClient moves a client between two neighbors within the same zone.
// prevID or nextID can be 0 to indicate start or end.
// It uses gap-based integer ordering and will trigger a reindex when gaps are exhausted.
func (s *ClientService) ReorderClient(clientID int, prevID int, nextID int) error {
	// Validate client exists and determine zone
	client, err := s.repo.GetByID(s.ctx, clientID)
	if err != nil {
		return fmt.Errorf("getting client: %w", err)
	}
	zoneRefUUID := client.ZoneID

	// Helper to get priority of a neighbor (0 if not present)
	getPriority := func(id int) (int, error) {
		if id == 0 {
			return 0, nil
		}
		c, err := s.repo.GetByID(s.ctx, id)
		if err != nil {
			return 0, fmt.Errorf("getting neighbor %d: %w", id, err)
		}
		return c.PriorityOrder, nil
	}

	prevPrio, err := getPriority(prevID)
	if err != nil {
		return err
	}
	nextPrio, err := getPriority(nextID)
	if err != nil {
		return err
	}

	// Compute new priority
	var newPriority int
	if prevID == 0 && nextID == 0 {
		// only client in zone
		newPriority = repository.DefaultPriorityGap
	} else if prevID == 0 {
		// place at beginning
		newPriority = nextPrio - repository.DefaultPriorityGap
		if newPriority <= 0 {
			// need reindex to make room
			if err := s.repo.ReindexZone(s.ctx, zoneRefUUID, repository.DefaultPriorityGap); err != nil {
				return fmt.Errorf("reindex failed: %w", err)
			}
			// recompute
			nextPrio, _ = getPriority(nextID)
			newPriority = nextPrio - repository.DefaultPriorityGap
		}
	} else if nextID == 0 {
		// place at end
		newPriority = prevPrio + repository.DefaultPriorityGap
	} else {
		// between two items
		if nextPrio-prevPrio > 1 {
			newPriority = (prevPrio + nextPrio) / 2
		} else {
			// collision, reindex and recompute
			if err := s.repo.ReindexZone(s.ctx, zoneRefUUID, repository.DefaultPriorityGap); err != nil {
				return fmt.Errorf("reindex failed: %w", err)
			}
			// re-fetch priorities
			prevPrio, _ = getPriority(prevID)
			nextPrio, _ = getPriority(nextID)
			if nextPrio-prevPrio <= 1 {
				// as a last resort, assign prev + 1
				newPriority = prevPrio + 1
			} else {
				newPriority = (prevPrio + nextPrio) / 2
			}
		}
	}

	// Persist change
	if err := s.repo.UpdatePriority(s.ctx, clientID, newPriority); err != nil {
		return fmt.Errorf("updating priority: %w", err)
	}
	return nil
}
