package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/cors"
	"github.com/theshubhamy/microGo/graphql"
)

type config struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
}

func main() {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	// Your custom server that provides resolvers
	customServer, err := graphql.NewGraphQLServer(cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL)
	if err != nil {
		log.Fatal("GraphQLServer error:", err)
	}

	// Build the gqlgen executable schema
	schema := customServer.ToExecutableSchema()

	// ✅ Create a real gqlgen handler.Server
	srv := handler.New(schema)
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	http.Handle("/graphql", srv)
	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	}).Handler(http.DefaultServeMux)

	log.Println("server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
