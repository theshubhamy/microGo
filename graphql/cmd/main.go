package main

import (
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-redis/redis/v8"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/cors"
	"github.com/theshubhamy/microGo/graphql"
	"github.com/theshubhamy/microGo/services/account"
)

type config struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
	RedisURL   string `envconfig:"REDIS_URL"`
}

func main() {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}
	account.LoadConfig()
	redisClient := redis.NewClient(opt)
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
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

	// Wrap with Auth middleware

	http.Handle("/graphql", graphql.AuthMiddleware(redisClient)(graphql.InjectRequestMeta(srv)))

	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	}).Handler(http.DefaultServeMux)

	log.Println("server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
