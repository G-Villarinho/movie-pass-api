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

type cinemaRepository struct {
	i           *do.Injector
	db          *gorm.DB
	redisClient *redis.Client
}

func NewCinemaRepository(i *do.Injector) (domain.CinemaRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}

	redisClient, err := do.Invoke[*redis.Client](i)
	if err != nil {
		return nil, err
	}

	return &cinemaRepository{
		i:           i,
		db:          db,
		redisClient: redisClient,
	}, nil
}

func (c *cinemaRepository) Create(ctx context.Context, cinema domain.Cinema) error {
	log := slog.With(
		slog.String("repository", "cinema"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing cinema creation process")

	if err := c.db.WithContext(ctx).Create(&cinema).Error; err != nil {
		log.Error("Failed to create cinema", slog.String("error", err.Error()))
		return err
	}

	log.Info("cinema creation process excuted succefully")
	return nil
}

func (c *cinemaRepository) GetByID(ctx context.Context, cinemaID uuid.UUID) (*domain.Cinema, error) {
	log := slog.With(
		slog.String("repository", "cinema"),
		slog.String("func", "GetByID"),
	)

	log.Info("Initializing get cinema by ID process")

	var cinema domain.Cinema
	if err := c.db.WithContext(ctx).Where("id = ?", cinemaID).First(&cinema).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("No cinema records found with the provided ID", slog.String("cinemaId", cinemaID.String()))
			return nil, nil
		}

		log.Error("Failed to get cinema by ID", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("Get cinema by ID process executed successfully")
	return &cinema, nil
}

func (c *cinemaRepository) GetAll(ctx context.Context, userID uuid.UUID) ([]domain.Cinema, error) {
	log := slog.With(
		slog.String("repository", "cinema"),
		slog.String("func", "GetAll"),
	)

	log.Info("Initializing get all cinemas process")

	var cinemas []domain.Cinema
	if err := c.db.Where("userId = ?", userID.String()).WithContext(ctx).Find(&cinemas).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("No cinema records found")
			return nil, nil
		}

		log.Error("Failed to get all cinemas", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("Get all cinemas process executed successfully")
	return cinemas, nil
}

func (c *cinemaRepository) Delete(ctx context.Context, cinemaID uuid.UUID) error {
	log := slog.With(
		slog.String("repository", "cinema"),
		slog.String("func", "Delete"),
	)

	log.Info("Initializing delete cinema process")

	if err := c.db.WithContext(ctx).Where("id = ?", cinemaID).Delete(&domain.Cinema{}).Error; err != nil {
		log.Error("Failed to delete cinema", slog.String("error", err.Error()))
		return err
	}

	log.Info("Delete cinema process executed successfully")
	return nil
}
