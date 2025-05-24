package main

import (
	"log"
	"time"

	"github.com/theshubhamy/microGo/services/account"
	"github.com/tinrab/retry"
)

func main() {
	account.LoadConfig()
	redisClient := account.InitRedis(account.AppConfig.REDIS_URL)
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("Error closing Redis client: %v", err)
		}
	}()

	var accRepo account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		accRepo, err = account.NewPostgresRepository(account.AppConfig.DATABASE_URL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer accRepo.Close()
	log.Println("Server running at 8080 ...")
	s := account.NewService(accRepo, redisClient)
	log.Fatal(account.ListenGrpcServer(s, 8080))
}
