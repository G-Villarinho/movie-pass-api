package main

import (
	"context"
	"log"
	"time"

	"github.com/GSVillas/movie-pass-api/config"
	"github.com/GSVillas/movie-pass-api/config/database"
	"github.com/GSVillas/movie-pass-api/domain"
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

	log.Println("Migration executed successfully")
}
