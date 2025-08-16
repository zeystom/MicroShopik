package main

import (
	"MicroShopik/configs"
	"MicroShopik/internal/database"
	"MicroShopik/scripts"
	"log"
)

func main() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	if err := scripts.RunSeedData(cfg); err != nil {
		log.Fatal("Failed to run seed data:", err)
	}
}
