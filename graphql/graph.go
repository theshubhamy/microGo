package graphql

import (
	"log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/theshubhamy/microGo/services/account"
	"github.com/theshubhamy/microGo/services/catalog"
	"github.com/theshubhamy/microGo/services/order"
)

type Server struct {
	accountClient *account.Client
	catalogClient *catalog.Client
	orderClient   *order.Client
}

func NewGraphQLServer(accountUrl, catalogURL, orderURL string) (*Server, error) {
	// Connect to account service
	accountClient, err := account.NewClient(accountUrl)
	if err != nil {
		log.Println("accountClienterror", err)
		return nil, err
	}

	// Connect to product service
	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return nil, err
	}

	// Connect to order service
	orderClient, err := order.NewClient(orderURL)
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return nil, err
	}

	return &Server{
		accountClient,
		catalogClient,
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

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}
