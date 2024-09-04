package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type MovieRepository struct {
	i           *do.Injector
	db          *gorm.DB
	redisClient *redis.Client
}

func NewMovieRepository(i *do.Injector) (domain.MovieRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}

	redisClient, err := do.Invoke[*redis.Client](i)
	if err != nil {
		return nil, err
	}

	return &MovieRepository{
		i:           i,
		db:          db,
		redisClient: redisClient,
	}, nil
}

func (m *MovieRepository) GetAllIndicativeRating(ctx context.Context, userID uuid.UUID) ([]domain.IndicativeRating, error) {
	log := slog.With(
		slog.String("repository", "movie"),
		slog.String("func", "GetAllIndicativeRating"),
	)

	log.Info("Initializing get all indicative rating process")

	var indicativeRating []domain.IndicativeRating
	if err := m.db.Where("userId = ?", userID.String()).WithContext(ctx).Find(&indicativeRating).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Warn("No indicative rating records found")
				return nil, nil
			}
		}

		log.Error("Failed to get all indicative rating", slog.String("error", err.Error()))
		return nil, err
	}

	if len(indicativeRating) == 0 {
		log.Warn("No indicative rating records found")
		return nil, nil
	}

	log.Info("Get all indicative rating process executed successfully")
	return indicativeRating, nil
}
