package account

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PutAccount(ctx context.Context, name string) (*Account, error)
	GetAccountbyId(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type accountService struct {
	repository Repository
}

func newService(repository Repository) Service {
	return &accountService{repository}
}

func (as *accountService) PostAccount(ctx context.Context, name string) (*Account, error) {
	account := &Account{
		Name: name,
		ID:   ksuid.New().String(),
	}

	err := as.repository.PutAccount(ctx, *account)
	if err != nil {
		return nil, err
	}
	return account, nil

}

func (as *accountService) GetAccount(ctx context.Context, id string) (*Account, error) {
	account, err := as.repository.GetAccountbyId(ctx, id)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (as *accountService) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	accounts, err := as.repository.ListAccounts(ctx, skip, take)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
