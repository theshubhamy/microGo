package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/theshubhamy/microGo/services/catalog"
	"github.com/tinrab/retry"
)

type Config struct {
	DATABASE_URL string `envconfig:"DATABASE_URL"`
}

func main() {
	var config Config
	log.Fatal("databaseURL", config.DATABASE_URL)
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	var r catalog.Repository

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = catalog.NewElasticRepository(config.DATABASE_URL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()
	log.Println("Server running at 8080 ...")
	s := catalog.NewService(r)
	log.Fatal(catalog.ListenGrpcServer(s, 8080))
}
