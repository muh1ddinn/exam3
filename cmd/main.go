package main

import (
	"context"
	"exam3/api"
	"exam3/config"
	"exam3/pkg/logger"
	"exam3/service"
	"exam3/storage/postgres"
	"exam3/storage/redis"

	"fmt"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.ServiceName)
	newRedis := redis.New(cfg)

	store, err := postgres.New(context.Background(), cfg, log)
	if err != nil {
		fmt.Println("error while connecting db, err: ", err)
		return
	}
	defer store.CloseDB()

	service := service.New(store, log, newRedis)

	c := api.New(service, log)

	fmt.Println("programm is running on localhost:9090...")
	c.Run(":8080")
}
