package order

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostOrder(ctx context.Context, accountId string, products []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error)
}
type Order struct {
	ID         string
	CreatedAt  time.Time
	TotalPrice float64
	AccountId  string
	Products   []OrderedProduct
}

type OrderedProduct struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Quantity    uint32
}

type orderService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &orderService{r}
}

func (os orderService) PostOrder(ctx context.Context, accountId string, products []OrderedProduct) (*Order, error) {
	order := &Order{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		AccountId: accountId,
		Products:  products,
	}

	order.TotalPrice = 0.0

	for _, p := range products {
		order.TotalPrice += p.Price * float64(p.Quantity)
	}
	err := os.repository.PutOrder(ctx, *order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (os orderService) GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error) {
	return os.repository.GetOrderforAccount(ctx, accountId)
}
