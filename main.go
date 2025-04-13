package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
	"github.com/theshubhamy/microGo/graphql"
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
	server, err := graphql.NewGraphqlServer(config.AccountURL, config.CatlogURL, config.OrderURL)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/graphql", handler.New(server.ToExecutableSchema()))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	log.Fatal(http.ListenAndServe(":3300", nil))

}
