package repository

import (
	"context"
	"errors"

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
	if err := c.db.WithContext(ctx).Create(&cinema).Error; err != nil {
		return err
	}

	return nil
}

func (c *cinemaRepository) GetByID(ctx context.Context, cinemaID uuid.UUID) (*domain.Cinema, error) {
	var cinema domain.Cinema
	if err := c.db.WithContext(ctx).Where("id = ?", cinemaID).First(&cinema).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &cinema, nil
}

func (c *cinemaRepository) GetAll(ctx context.Context, userID uuid.UUID) ([]domain.Cinema, error) {
	var cinemas []domain.Cinema
	if err := c.db.Where("userId = ?", userID.String()).WithContext(ctx).Find(&cinemas).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return cinemas, nil
}

func (c *cinemaRepository) Delete(ctx context.Context, cinemaID uuid.UUID) error {
	if err := c.db.WithContext(ctx).Where("id = ?", cinemaID).Delete(&domain.Cinema{}).Error; err != nil {
		return err
	}

	return nil
}
