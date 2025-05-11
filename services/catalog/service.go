package catalog

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostProduct(ctx context.Context, name string, description string, price string) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	GetProductsbyIds(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type Product struct {
	ID          string `json:"Id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

type catalogService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &catalogService{r}
}

func (cs *catalogService) PostProduct(ctx context.Context, name string, description string, price string) (*Product, error) {
	product := &Product{
		ID:          ksuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
	}

	err := cs.repository.PutProduct(ctx, *product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (cs *catalogService) GetProduct(ctx context.Context, id string) (*Product, error) {
	product, err := cs.repository.GetProductbyId(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (cs *catalogService) GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	products, err := cs.repository.ListProducts(ctx, skip, take)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (cs *catalogService) GetProductsbyIds(ctx context.Context, ids []string) ([]Product, error) {
	products, err := cs.repository.ListProductsWithIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (cs *catalogService) SearchProducts(ctx context.Context,
	query string, skip uint64, take uint64,
) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	products, err := cs.repository.SearchProduct(ctx, query, skip, take)
	if err != nil {
		return nil, err
	}
	return products, nil
}
