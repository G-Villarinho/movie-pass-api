package main

import (
	"context"
	"log"
	"time"

	"github.com/GSVillas/movie-pass-api/config"
	"github.com/GSVillas/movie-pass-api/config/database"
	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/GSVillas/movie-pass-api/secure"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func main() {
	config.LoadEnvironments()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.NewMysqlConnection(ctx)
	if err != nil {
		log.Fatal("Fail to connect to mysql: ", err)
	}

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Cinema{},
		&domain.CinemaSession{},
		&domain.CinemaRoom{},
		&domain.IndicativeRating{},
		&domain.Movie{},
		&domain.MovieImage{},
		&domain.SeatReservation{},
		&domain.Seat{},
		&domain.Role{},
	); err != nil {
		log.Fatal("Fail to migrate: ", err)
	}

	populateIndicativeRatings(db)
	pupulateRoles(db)
	populateSuperAdmin(db)

	log.Println("Migration executed successfully")
}

func populateIndicativeRatings(db *gorm.DB) {
	indicativeRatings := []domain.IndicativeRating{
		{
			ID:          uuid.MustParse("dffab792-689b-11ef-b065-0242ac110002"),
			Description: "AL",
			ImageURL:    "https://imagedelivery.net/Zphe8Y_ziiz_0wgSXjC_Qg/eaa56fbb-6f29-4be8-1bcb-a717e5f1e900/public",
		},
		{
			ID:          uuid.MustParse("dffb1391-689b-11ef-b065-0242ac110002"),
			Description: "A10",
			ImageURL:    "https://imagedelivery.net/Zphe8Y_ziiz_0wgSXjC_Qg/bc21525d-e8e0-4b1e-3cfb-baaae79b6800/public",
		},
		{
			ID:          uuid.MustParse("dffb1acb-689b-11ef-b065-0242ac110002"),
			Description: "A12",
			ImageURL:    "https://imagedelivery.net/Zphe8Y_ziiz_0wgSXjC_Qg/e8f234fa-b114-456f-f51b-0a5ce3c15700/public",
		},
		{
			ID:          uuid.MustParse("dffb1b9a-689b-11ef-b065-0242ac110002"),
			Description: "A14",
			ImageURL:    "https://imagedelivery.net/Zphe8Y_ziiz_0wgSXjC_Qg/76c45af1-dcd9-40e7-9188-b87cfe71f600/public",
		},
		{
			ID:          uuid.MustParse("dffb1d12-689b-11ef-b065-0242ac110002"),
			Description: "A16",
			ImageURL:    "https://imagedelivery.net/Zphe8Y_ziiz_0wgSXjC_Qg/d123bf65-a7ad-4c86-bab7-358679d6d000/public",
		},
		{
			ID:          uuid.MustParse("dffb1d82-689b-11ef-b065-0242ac110002"),
			Description: "A18",
			ImageURL:    "https://imagedelivery.net/Zphe8Y_ziiz_0wgSXjC_Qg/36cccb63-b5cd-4ed4-f8cf-e01dad4bff00/public",
		},
	}

	var existingRatings []domain.IndicativeRating
	if err := db.Find(&existingRatings).Error; err != nil {
		log.Printf("Error retrieving existing ratings: %v", err)
		return
	}

	existingRatingsMap := make(map[uuid.UUID]bool)
	for _, rating := range existingRatings {
		existingRatingsMap[rating.ID] = true
	}

	for _, rating := range indicativeRatings {
		if !existingRatingsMap[rating.ID] {
			if err := db.Create(&rating).Error; err != nil {
				log.Printf("Error inserting rating %s: %v", rating.Description, err)
			} else {
				log.Printf("Inserted rating %s", rating.Description)
			}
		} else {
			log.Printf("Rating %s already exists, skipping", rating.Description)
		}
	}
}

func pupulateRoles(db *gorm.DB) {
	var roles = []domain.Role{
		{ID: uuid.MustParse("aab5c388-1559-4a09-9c64-88aaa94fe9c3"), Name: string(domain.AdminRoleLevel1), Description: "Administrator role with basic privileges"},
		{ID: uuid.MustParse("62f31a37-f586-4d51-a593-ee292b1e5090"), Name: string(domain.AdminRoleLevel2), Description: "Administrator role with advanced privileges"},
		{ID: uuid.MustParse("1f7d5ea5-3994-4823-bbe5-a5aad8c7322c"), Name: string(domain.AdminRoleLevel3), Description: "Super administrator role"},
		{ID: uuid.MustParse("7efe4510-169f-4511-8e12-de7ebaea31e5"), Name: string(domain.UserRole), Description: "Authenticated user role"},
	}

	for _, role := range roles {
		if err := db.Where("name = ?", role.Name).FirstOrCreate(&role).Error; err != nil {
			log.Printf("Error inserting role %s: %v", role.Name, err)
		} else {
			log.Printf("Inserted role %s", role.Name)
		}
	}
}

func populateSuperAdmin(db *gorm.DB) {
	var superAdminRole domain.Role
	if err := db.Where("name = ?", string(domain.AdminRoleLevel3)).First(&superAdminRole).Error; err != nil {
		log.Printf("Super admin role not found: %v", err)
		return
	}

	superAdmin := domain.User{
		ID:           uuid.New(),
		FirstName:    "Alexandre",
		LastName:     "Falcão",
		Email:        config.Env.SuperAdminEmail,
		PasswordHash: config.Env.SuperAdminPassword,
		RoleID:       superAdminRole.ID,
		BirthDate:    time.Now().UTC(),
		CreatedAt:    time.Now().UTC(),
	}

	var existingUser domain.User
	if err := db.Where("email = ?", superAdmin.Email).First(&existingUser).Error; err == nil {
		log.Printf("Super admin user already exists, skipping")
		return
	}

	passwordHash, err := secure.HashPassword(config.Env.SuperAdminPassword)
	if err != nil {
		panic(err)
	}

	superAdmin.PasswordHash = string(passwordHash)

	if err := db.Create(&superAdmin).Error; err != nil {
		log.Printf("Error creating super admin: %v", err)
	} else {
		log.Println("Super admin created successfully")
	}
}
