package graphql

import (
	"context"
	"errors"
	"log"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Get single
	if id != nil {
		r, err := r.server.catalogClient.GetProduct(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*Product{{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Price:       r.Price,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	q := ""
	if query != nil {
		q = *query
	}
	productList, err := r.server.catalogClient.GetProducts(ctx, q, nil, skip, take)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var products []*Product
	for _, a := range *productList {
		products = append(products,
			&Product{
				ID:          a.ID,
				Name:        a.Name,
				Description: a.Description,
				Price:       a.Price,
			},
		)
	}

	return products, nil
}

func (r *queryResolver) Orders(ctx context.Context) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok || userID == "" {
		return nil, errors.New("unauthorized: user ID not found")
	}
	orderList, err := r.server.orderClient.GetOrdersForAccount(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var orders []*Order
	for _, o := range orderList {
		var products []*OrderedProduct
		for _, p := range o.Products {
			products = append(products, &OrderedProduct{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Quantity:    int(p.Quantity),
			})
		}
		orders = append(orders, &Order{
			ID:         o.ID,
			CreatedAt:  o.CreatedAt,
			TotalPrice: o.TotalPrice,
			Products:   products,
		})
	}

	return orders, nil
}

func (p PaginationInput) bounds() (uint64, uint64) {
	skipValue := uint64(0)
	takeValue := uint64(100)
	if p.Skip != nil {
		skipValue = uint64(*p.Skip)
	}
	if p.Take != nil {
		takeValue = uint64(*p.Take)
	}
	return skipValue, takeValue
}
