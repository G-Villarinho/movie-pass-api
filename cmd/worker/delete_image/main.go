package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/GSVillas/movie-pass-api/client"
	"github.com/GSVillas/movie-pass-api/config"
	"github.com/GSVillas/movie-pass-api/config/database"
	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/GSVillas/movie-pass-api/repository"
	"github.com/GSVillas/movie-pass-api/service"
	"github.com/go-redis/redis/v8"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func main() {
	config.ConfigureLogger()
	config.LoadEnvironments()

	i := do.New()

	db, err := database.NewMysqlConnection(context.Background())
	if err != nil {
		log.Fatal("Fail to connect to mysql: ", err)
	}

	redisClient, err := database.NewRedisConnection(context.Background())
	if err != nil {
		log.Fatal("Fail to connect to redis: ", err)
	}

	do.Provide(i, func(i *do.Injector) (*gorm.DB, error) {
		return db, nil
	})

	do.Provide(i, func(i *do.Injector) (*redis.Client, error) {
		return redisClient, nil
	})

	do.Provide(i, client.NewCloudFlareService)

	do.Provide(i, service.NewMovieService)

	do.Provide(i, repository.NewMovieRepository)

	movieRepository, err := do.Invoke[domain.MovieRepository](i)
	if err != nil {
		panic(err)
	}

	movieService, err := do.Invoke[domain.MovieService](i)
	if err != nil {
		panic(err)
	}

	for {
		task, err := movieRepository.GetNextDeleteTask(context.Background())
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		if task == nil {
			continue
		}

		slog.Info("start delete image in cloud")
		if err := movieService.ProcessDeleteQueue(context.Background(), *task); err != nil {
			slog.Error(err.Error())
			continue
		}

		slog.Info("Image delete sucesfully")
	}
}
