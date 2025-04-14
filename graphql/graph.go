package graphql

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/theshubhamy/microGo/services/account"
)

type Server struct {
	accountClient *account.Client
	catlogClient  *catlog.Client
	orderClient   *order.Client
}

func NewGraphqlServer(accountURL, catalogURL, orderURL string) (*Server, error) {
	accountClient, err := account.NewClient(accountURL)

	if err != nil {
		accountClient.Close()
		return nil, err
	}
	catlogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		catlogClient.Close()
		return nil, err
	}
	orderClient, err := order.NewClient(orderURL)

	if err != nil {
		orderClient.Close()
		return nil, err
	}

	return &Server{
		accountClient,
		catlogClient,
		orderClient,
	}, nil
}

func (s *Server) Mutation() MutationResolver {
	return &mutationResolver{
		server: s,
	}
}

func (s *Server) Query() QueryResolver {
	return &queryResolver{
		server: s,
	}
}

func (s *Server) Account() AccountResolver {
	return &accountResolver{
		server: s,
	}
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{Resolvers: s})
}
