package graphql

import "context"

type QueryResolver struct {
	server *Server
}

// Accounts
// Products

func (r *QueryResolver) GetAccounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	return nil, nil
}
func (r *QueryResolver) GetProducts(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	return nil, nil
}
