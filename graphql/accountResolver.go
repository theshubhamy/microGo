package graphql

import "context"

type AccountResolver struct {
	server *Server
}

// Orders
func (r *AccountResolver) GetOrders(ctx context.Context, obj *Account) ([]*Order, error) {
	return nil, nil
}
