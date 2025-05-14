package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/cors"
)

type config struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatlogURL  string `envconfig:"CATLOG_SERVICE_URL"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
}

func main() {
	var config config

	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}
	server, err := NewGraphQLServer(config.AccountURL, config.CatlogURL, config.OrderURL)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/graphql", handler.New(server.ToExecutableSchema()))
	http.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))

	// Add CORS middleware to allow cross-origin requests
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow all origins
	}).Handler(http.DefaultServeMux)

	log.Println("server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
