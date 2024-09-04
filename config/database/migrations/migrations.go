package main

import (
	"context"
	"log"
	"time"

	"github.com/GSVillas/movie-pass-api/config"
	"github.com/GSVillas/movie-pass-api/config/database"
	"github.com/GSVillas/movie-pass-api/domain"
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
		&domain.SeatReservation{},
		&domain.Seat{},
	); err != nil {
		log.Fatal("Fail to migrate: ", err)
	}

	populateIndicativeRatings(db)

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
