package service

import (
	"context"
	"fmt"

	"soditapp/backend/model"
	"soditapp/backend/repository"
)

// ProductService exposes product operations as Wails-bound methods.
type ProductService struct {
	ctx  context.Context
	repo *repository.ProductRepository
}

// Init sets the runtime dependencies after Wails startup.
func (s *ProductService) Init(ctx context.Context, repo *repository.ProductRepository) {
	s.ctx = ctx
	s.repo = repo
}

// Reinit updates runtime dependencies (called after database import/reopen).
func (s *ProductService) Reinit(repo *repository.ProductRepository) {
	s.repo = repo
}

func (s *ProductService) CreateProduct(name, typ string) (model.Product, error) {
	if name == "" {
		return model.Product{}, fmt.Errorf("product name cannot be empty")
	}
	return s.repo.Create(s.ctx, name, typ)
}

func (s *ProductService) ListProducts() ([]model.Product, error) {
	return s.repo.GetAll(s.ctx)
}

func (s *ProductService) ListActiveProducts() ([]model.Product, error) {
	return s.repo.GetActive(s.ctx)
}

func (s *ProductService) UpdateProduct(id int, name string) error {
	if name == "" {
		return fmt.Errorf("product name cannot be empty")
	}
	return s.repo.Update(s.ctx, id, name)
}

func (s *ProductService) ToggleProductActive(id int) error {
	return s.repo.ToggleActive(s.ctx, id)
}
