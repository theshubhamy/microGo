package graphql

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/theshubhamy/microGo/services/order"
)

var ErrInvalidParameter = errors.New("invalid parameter")

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, in AccountInput) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	a, err := r.server.accountClient.PostAccount(ctx, in.Name, in.Email, in.Phone, in.Password)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Account{
		ID:    a.ID,
		Name:  a.Name,
		Email: a.Email,
		Phone: a.Phone,
	}, nil
}

func (r *mutationResolver) LoginAccount(ctx context.Context, in LoginInput) (*LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	ip, ok := ctx.Value(ctxKeyIP).(string)
	if !ok {
		return nil, errors.New("missing IP from context")
	}
	userAgent, ok := ctx.Value(ctxKeyUserAgent).(string)
	if !ok {
		return nil, errors.New("missing user agent from context")
	}

	acc, accessToken, refreshToken, err := r.server.accountClient.LoginAccount(ctx, in.Emailorphone, in.Password, ip, userAgent)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &LoginResponse{
		ID:           acc.ID,
		Name:         acc.Name,
		Email:        acc.Email,
		Phone:        acc.Phone,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, in ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	p, err := r.server.catalogClient.PostProduct(ctx, in.Name, in.Description, in.Price)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, in OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var products []order.OrderedProduct
	for _, p := range in.Products {
		if p.Quantity <= 0 {
			return nil, ErrInvalidParameter
		}
		products = append(products, order.OrderedProduct{
			ID:       p.ID,
			Quantity: uint32(p.Quantity),
		})
	}
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok || userID == "" {
		return nil, errors.New("unauthorized: user ID not found")
	}
	o, err := r.server.orderClient.PostOrder(ctx, userID, products)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Order{
		ID:         o.ID,
		CreatedAt:  o.CreatedAt,
		TotalPrice: o.TotalPrice,
	}, nil
}
