package scripts

import (
	"MicroShopik/configs"
	"MicroShopik/internal/database"
	"MicroShopik/internal/repositories"
	"log"
)

func RunSeedData(cfg *configs.Config) error {
	return runSeedData(cfg)
}

func RunSeedRoles(cfg *configs.Config) error {
	return runSeedRoles(cfg)
}

func RunInitialSeed(cfg *configs.Config) error {
	userRepo := repositories.NewUserRepository(database.GetDB())
	roleRepo := repositories.NewRoleRepository(database.GetDB())

	roles, err := roleRepo.GetAll()
	if err != nil {
		log.Printf("Warning: Could not check existing roles: %v", err)
		return err
	}

	users, err := userRepo.GetAll()
	if err != nil {
		log.Printf("Warning: Could not check existing users: %v", err)
		return err
	}

	if len(roles) > 0 && len(users) > 0 {
		log.Println("Database already contains data, skipping seed...")
		return nil
	}

	log.Println("Database is empty, running initial seed...")

	if err := RunSeedRoles(cfg); err != nil {
		log.Printf("Warning: Failed to seed roles: %v", err)
		return err
	}

	if err := RunSeedData(cfg); err != nil {
		log.Printf("Warning: Failed to seed data: %v", err)
		return err
	}

	log.Println("Initial seed completed successfully!")
	return nil
}
