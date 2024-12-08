package catalog

import (
	"context"
	"github.com/segmentio/ksuid"
)

type Service interface {
	Close()
	CreateProduct(ctx context.Context, name, description string, price, quantity uint64) (*Product, error)
	GetProductByID(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip, take uint64) ([]*Product, error)
	GetProductByIDs(ctx context.Context, ids []string) ([]*Product, error)
	SearchProducts(ctx context.Context, query string, skip, take uint64) ([]*Product, error)
}

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint64 `json:"price"`
	Quantity    uint64 `json:"quantity"`
}

type catalogService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &catalogService{repo: repo}
}

func (s *catalogService) Close() {
	s.repo.Close()
}

func (s *catalogService) CreateProduct(ctx context.Context, name, description string, price, quantity uint64) (*Product, error) {
	p := &Product{ID: ksuid.New().String(), Name: name, Description: description, Price: price, Quantity: quantity}

	if err := s.repo.UpdateProduct(ctx, *p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *catalogService) GetProductByID(ctx context.Context, id string) (*Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *catalogService) GetProducts(ctx context.Context, skip, take uint64) ([]*Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	return s.repo.GetProducts(ctx, skip, take)
}

func (s *catalogService) GetProductByIDs(ctx context.Context, ids []string) ([]*Product, error) {
	return s.repo.GetProductByIDs(ctx, ids)
}

func (s *catalogService) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]*Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	return s.repo.SearchProducts(ctx, query, skip, take)
}
