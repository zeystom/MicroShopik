package main

import (
	"MicroShopik/configs"
	"MicroShopik/internal/database"
	"MicroShopik/internal/domain"
	"MicroShopik/internal/repositories"
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

	roleRepo := repositories.NewRoleRepository(database.GetDB())

	roles := []domain.Role{
		{
			Name:        "admin",
			Description: "Administrator with full system access",
		},
		{
			Name:        "seller",
			Description: "Seller who can create and manage products",
		},
		{
			Name:        "buyer",
			Description: "Regular user who can browse and purchase products",
		},
		{
			Name:        "moderator",
			Description: "Moderator who can review and approve content",
		},
	}

	for _, role := range roles {
		err := roleRepo.Create(&role)
		if err != nil {
			log.Printf("Failed to create role %s: %v", role.Name, err)
		} else {
			log.Printf("Created role: %s", role.Name)
		}
	}

	log.Println("Role seeding completed!")
}
