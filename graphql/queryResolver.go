package graphql

import "context"

type queryResolver struct {
	server *Server
}

// Accounts
// Products

func (r *mutationResolver) GetAccounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
}
func (r *mutationResolver) GetProducts(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {

}
