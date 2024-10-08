package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"

	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/go-redis/redis/v8"
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

func (m *MovieRepository) GetAllIndicativeRating(ctx context.Context) ([]*domain.IndicativeRating, error) {
	var indicativeRating []*domain.IndicativeRating
	if err := m.db.WithContext(ctx).Find(&indicativeRating).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return indicativeRating, nil
}

func (m *MovieRepository) GetIndicativeRatingByID(ctx context.Context, id uuid.UUID) (*domain.IndicativeRating, error) {
	var indicativeRating domain.IndicativeRating

	if err := m.db.WithContext(ctx).First(&indicativeRating, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &indicativeRating, nil
}

func (m *MovieRepository) Create(ctx context.Context, movie domain.Movie) error {
	if err := m.db.WithContext(ctx).Create(&movie).Error; err != nil {
		return err
	}

	return nil
}

func (m *MovieRepository) CreateMovieImage(ctx context.Context, movieImage domain.MovieImage) error {
	if err := m.db.WithContext(ctx).Create(&movieImage).Error; err != nil {
		return err
	}

	return nil
}

func (m *MovieRepository) AddUploadTaskToQueue(ctx context.Context, task domain.MovieImageUploadTask) error {
	data, err := jsoniter.Marshal(task)
	if err != nil {
		return fmt.Errorf("error to serialize task for Redis. Error: %w", err)
	}

	if err := m.redisClient.RPush(ctx, m.getImageUploadKey(), data).Err(); err != nil {
		return fmt.Errorf("error to add task to Redis queue. Error: %w", err)
	}

	return nil
}

func (m *MovieRepository) GetNextUploadTask(ctx context.Context) (*domain.MovieImageUploadTask, error) {
	data, err := m.redisClient.LPop(ctx, m.getImageUploadKey()).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	var task domain.MovieImageUploadTask
	if err := jsoniter.Unmarshal([]byte(data), &task); err != nil {
		return nil, fmt.Errorf("error to deserialize task from Redis. error:%w", err)
	}

	return &task, nil
}

func (m *MovieRepository) GetALlByUserID(ctx context.Context, userID uuid.UUID, pagination *domain.Pagination) (*domain.Pagination, error) {
	var movies []*domain.Movie
	if err := m.db.WithContext(ctx).
		Where("userId = ?", userID.String()).
		Scopes(paginate(&movies, pagination, m.db)).
		Preload("Images").
		Preload("IndicativeRating").
		Find(&movies).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	for _, movie := range movies {
		if len(movie.Images) > 0 {
			movie.Images = movie.Images[:1]
		}
	}

	pagination.Rows = movies

	return pagination, nil
}

func (m *MovieRepository) Update(ctx context.Context, movie domain.Movie) error {
	if err := m.db.Save(movie).Error; err != nil {
		return err
	}

	return nil
}

func (m *MovieRepository) GetByID(ctx context.Context, ID uuid.UUID, withPreload bool) (*domain.Movie, error) {
	var movie domain.Movie
	db := m.db.WithContext(ctx)

	if withPreload {
		db = db.Preload("Images").Preload("IndicativeRating")
	}

	if err := db.First(&movie, "id = ?", ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &movie, nil
}

func (m *MovieRepository) DeleteMovieImage(ctx context.Context, cloudFlareID uuid.UUID) error {
	if err := m.db.WithContext(ctx).Where("cloudFlareId = ?", cloudFlareID).Delete(&domain.MovieImage{}).Error; err != nil {
		return err
	}

	return nil
}

func (m *MovieRepository) AddDeleteTaskToQueue(ctx context.Context, task domain.MovieImageDeleteTask) error {
	data, err := jsoniter.Marshal(task)
	if err != nil {
		return fmt.Errorf("error to serialize task for Redis. Error: %w", err)
	}

	if err := m.redisClient.RPush(ctx, m.getImageDeleteKey(), data).Err(); err != nil {
		return fmt.Errorf("error to add task to Redis queue. Error: %w", err)
	}

	return nil
}

func (m *MovieRepository) GetNextDeleteTask(ctx context.Context) (*domain.MovieImageDeleteTask, error) {
	data, err := m.redisClient.LPop(ctx, m.getImageDeleteKey()).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	var task domain.MovieImageDeleteTask
	if err := jsoniter.Unmarshal([]byte(data), &task); err != nil {
		return nil, fmt.Errorf("error to deserialize task from Redis. error:%w", err)
	}

	return &task, nil
}

func (m *MovieRepository) getImageUploadKey() string {
	return "image_upload_queue"
}

func (m *MovieRepository) getImageDeleteKey() string {
	return "image_delete_queue"
}
