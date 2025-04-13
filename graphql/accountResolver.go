package graphql

import "context"

type accountResolver struct {
	server *Server
}

// Orders
func (r *mutationResolver) GetOrders(ctx context.Context, obj *Account) ([]*Order, error) {
	return nil, nil
}
