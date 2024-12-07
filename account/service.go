package account

import (
	"context"
)

type Service interface {
	Close()
	Ping() error
	CreateAccount(ctx context.Context, name string) (*Account, error)
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, skip, take uint64) ([]*Account, error)
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type accountService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &accountService{repo: repo}
}

func (s *accountService) Close() {
	s.repo.Close()
}

func (s *accountService) Ping() error {
	return s.repo.Ping()
}

func (s *accountService) CreateAccount(ctx context.Context, name string) (*Account, error) {
	return s.repo.CreateAccount(ctx, name)
}

func (s *accountService) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	return s.repo.GetAccountByID(ctx, id)
}

func (s *accountService) GetAccounts(ctx context.Context, skip, take uint64) ([]*Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	return s.repo.GetAccounts(ctx, skip, take)
}
