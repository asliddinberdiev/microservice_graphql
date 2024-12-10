package order

import (
	"context"
	"github.com/segmentio/ksuid"
	"time"
)

type Service interface {
	Close()
	Create(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountId string) ([]*Order, error)
}

type Order struct {
	ID         string            `json:"id"`
	AccountID  string            `json:"account_id"`
	TotalPrice uint64           `json:"total_price"`
	CreatedAt  time.Time         `json:"created_at"`
	Products   []OrderedProduct `json:"products"`
}

type OrderedProduct struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       uint64  `json:"price"`
	Quantity    uint64  `json:"quantity"`
}

type orderService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &orderService{repo: repo}
}

func (s *orderService) Close() {
	s.repo.Close()
}

func (s *orderService) Create(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	order := &Order{ID: ksuid.New().String(), AccountID: accountID, TotalPrice: 0, CreatedAt: time.Now(), Products: products}
	order.TotalPrice = 0

	for _, p := range products {
		order.TotalPrice += p.Price * p.Quantity
	}

	if err := s.repo.UpdateOrder(ctx, *order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) GetOrdersForAccount(ctx context.Context, accountId string) ([]*Order, error) {
	return s.repo.GetOrdersForAccount(ctx, accountId)
}
