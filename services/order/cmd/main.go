package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/theshubhamy/microGo/services/order"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
	AccountURL  string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL  string `envconfig:"CATALOG_SERVICE_URL"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}
	var r order.Repository

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = order.NewPostgresRepository(config.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()
	log.Println("Server running at 8080 ...")
	s := order.NewService(r)
	log.Fatal(order.ListenGrpcServer(s, config.AccountURL, config.CatalogURL, 8080))
}
