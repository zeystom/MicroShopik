package scripts

import (
	"MicroShopik/configs"
	"MicroShopik/internal/database"
	"MicroShopik/internal/domain"
	"MicroShopik/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// runSeedData runs the seed data function that can be called from main.go
func runSeedData(_ *configs.Config) error {
	userRepo := repositories.NewUserRepository(database.GetDB())
	categoryRepo := repositories.NewCategoryRepository(database.GetDB())
	productRepo := repositories.NewProductRepository(database.GetDB())

	adminHash, err := bcrypt.GenerateFromPassword([]byte("admin12345"), bcrypt.DefaultCost)
	sellerHash, err := bcrypt.GenerateFromPassword([]byte("seller12345"), bcrypt.DefaultCost)
	buyerHash, err := bcrypt.GenerateFromPassword([]byte("buyer12345"), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	users := []domain.User{
		{
			Username: "admin",
			Email:    "admin@microshopik.com",
			Password: string(adminHash),
		},
		{
			Username: "seller1",
			Email:    "seller1@microshopik.com",
			Password: string(sellerHash),
		},
		{
			Username: "buyer1",
			Email:    "buyer1@microshopik.com",
			Password: string(buyerHash),
		},
	}

	for i := range users {
		err := userRepo.Create(&users[i])
		if err != nil {
			log.Printf("Failed to create user %s: %v", users[i].Username, err)
		} else {
			log.Printf("Created user: %s", users[i].Username)
		}
	}

	roleAssignments := map[string][]string{
		"admin":   {"admin"},
		"seller1": {"seller"},
		"buyer1":  {"buyer"},
	}

	for username, roles := range roleAssignments {
		user, err := userRepo.GetByEmail(username + "@microshopik.com")
		if err != nil {
			log.Printf("Failed to get user %s: %v", username, err)
			continue
		}

		for _, roleName := range roles {
			err := userRepo.AssignRole(user.ID, roleName)
			if err != nil {
				log.Printf("Failed to assign role %s to user %s: %v", roleName, username, err)
			} else {
				log.Printf("Assigned role %s to user %s", roleName, username)
			}
		}
	}

	categories := []domain.Category{
		{
			Name: "Gaming",
		},
		{
			Name: "Streaming",
		},
		{
			Name: "Software",
		},
		{
			Name: "Music",
		},
	}

	for i := range categories {
		err := categoryRepo.Create(&categories[i])
		if err != nil {
			log.Printf("Failed to create category %s: %v", categories[i].Name, err)
		} else {
			log.Printf("Created category: %s", categories[i].Name)
		}
	}

	seller, err := userRepo.GetByEmail("seller1@microshopik.com")
	if err != nil {
		log.Printf("Failed to get seller user: %v", err)
		return err
	}

	products := []domain.Product{
		{
			SellerID:    seller.ID,
			Title:       "Premium Netflix Account",
			Description: "4K Ultra HD streaming, 4 screens, no ads, premium content access",
			Price:       2999,
			CategoryID:  2,
			IsActive:    true,
			MaxSales:    1000,
			SoldCount:   1247,
		},
		{
			SellerID:    seller.ID,
			Title:       "Valorant High Rank Account",
			Description: "Immortal rank account with 50+ skins, completed battle pass, all agents unlocked",
			Price:       8999,
			CategoryID:  1,
			IsActive:    true,
			MaxSales:    500,
			SoldCount:   892,
		},
		{
			SellerID:    seller.ID,
			Title:       "Spotify Premium Family",
			Description: "6 accounts, ad-free, offline downloads, high quality audio",
			Price:       1999,
			CategoryID:  4,
			IsActive:    true,
			MaxSales:    2000,
			SoldCount:   2156,
		},
		{
			SellerID:    seller.ID,
			Title:       "Adobe Creative Suite License",
			Description: "Full access to Photoshop, Illustrator, Premiere Pro, After Effects",
			Price:       4999,
			CategoryID:  3,
			IsActive:    true,
			MaxSales:    800,
			SoldCount:   567,
		},
	}

	for i := range products {
		productID, err := productRepo.Create(&products[i])
		if err != nil {
			log.Printf("Failed to create product %s: %v", products[i].Title, err)
		} else {
			log.Printf("Created product: %s with ID: %d", products[i].Title, productID)

		}
	}

	log.Println("Data seeding completed!")
	return nil
}
