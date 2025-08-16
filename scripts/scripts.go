package scripts

import (
	"MicroShopik/configs"
	"MicroShopik/internal/database"
	"MicroShopik/internal/repositories"
	"log"
)

// RunSeedData runs the seed data function that can be called from main.go
func RunSeedData(cfg *configs.Config) error {
	return runSeedData(cfg)
}

// RunSeedRoles runs the role seed data function that can be called from main.go
func RunSeedRoles(cfg *configs.Config) error {
	return runSeedRoles(cfg)
}

// RunInitialSeed checks if database is empty and runs seed data only if needed
func RunInitialSeed(cfg *configs.Config) error {
	// Check if database already has data
	userRepo := repositories.NewUserRepository(database.GetDB())
	roleRepo := repositories.NewRoleRepository(database.GetDB())

	// Check if roles exist
	roles, err := roleRepo.GetAll()
	if err != nil {
		log.Printf("Warning: Could not check existing roles: %v", err)
		return err
	}

	// Check if users exist
	users, err := userRepo.GetAll()
	if err != nil {
		log.Printf("Warning: Could not check existing users: %v", err)
		return err
	}

	// If database already has data, skip seeding
	if len(roles) > 0 && len(users) > 0 {
		log.Println("Database already contains data, skipping seed...")
		return nil
	}

	log.Println("Database is empty, running initial seed...")

	// First create roles, then users and other data
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
