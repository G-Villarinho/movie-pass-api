package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

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
	log := slog.With(
		slog.String("repository", "movie"),
		slog.String("func", "GetAllIndicativeRating"),
	)

	log.Info("Initializing get all indicative rating process")

	var indicativeRating []*domain.IndicativeRating
	if err := m.db.WithContext(ctx).Find(&indicativeRating).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("No indicative rating records found")
			return nil, nil
		}

		log.Error("Failed to get all indicative rating", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("Get all indicative rating process executed successfully")
	return indicativeRating, nil
}

func (m *MovieRepository) Create(ctx context.Context, movie domain.Movie) error {
	log := slog.With(
		slog.String("repository", "movie"),
		slog.String("func", "create"),
	)

	log.Info("Initializing create movie process")

	if err := m.db.WithContext(ctx).Create(&movie).Error; err != nil {
		log.Error("Failed to create movie", slog.String("error", err.Error()))
		return err
	}

	log.Info("movie creation process excuted succefully")
	return nil
}

func (m *MovieRepository) CreateMovieImage(ctx context.Context, movieImage domain.MovieImage) error {
	log := slog.With(
		slog.String("repository", "movie"),
		slog.String("func", "create"),
	)

	log.Info("Initializing create movie image process")

	if err := m.db.WithContext(ctx).Create(&movieImage).Error; err != nil {
		log.Error("Failed to create movie image", slog.String("error", err.Error()))
		return err
	}

	log.Info("movie image creation process excuted succefully")
	return nil
}

func (m *MovieRepository) AddUploadImageTaskToQueue(ctx context.Context, task domain.MovieImageUploadTask) error {
	log := slog.With(
		slog.String("repository", "movie"),
		slog.String("func", "AddUploadImageTaskToQueue"),
		slog.String("movieID", task.MovieID.String()),
		slog.String("userID", task.UserID.String()),
	)

	log.Info("Adding image upload task to Redis queue")

	data, err := jsoniter.Marshal(task)
	if err != nil {
		log.Error("Failed to serialize task for Redis", slog.String("error", err.Error()))
		return err
	}

	if err := m.redisClient.RPush(ctx, m.getImageUploadKey(), data).Err(); err != nil {
		log.Error("Failed to add task to Redis queue", slog.String("error", err.Error()))
		return err
	}

	log.Info("Image upload task added to Redis queue successfully")
	return nil
}

func (m *MovieRepository) GetNextUploadImageTaskFromQueue(ctx context.Context) (*domain.MovieImageUploadTask, error) {
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

func (m *MovieRepository) GetALlByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Movie, error) {
	log := slog.With(
		slog.String("repository", "movie"),
		slog.String("func", "GetALlByUserID"),
	)

	log.Info("Initializing create movie image process")

	var movies []*domain.Movie
	if err := m.db.WithContext(ctx).Where("userId = ?", userID.String()).Preload("Images").Preload("IndicativeRating").Find(&movies).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("No movies records found")
			return nil, nil
		}

		log.Error("Failed to get all movies by user id", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("Get all movies by user id process executed successfully")
	return movies, nil
}

func (m *MovieRepository) Update(ctx context.Context, movieID uuid.UUID, updates map[string]any) error {
	log := slog.With(
		slog.String("repository", "movie"),
		slog.String("func", "Update"),
		slog.String("movieID", movieID.String()),
	)

	log.Info("Initializing update movie process")

	if err := m.db.Model(&domain.Movie{}).Where("id = ?", movieID).Updates(updates).Error; err != nil {
		log.Error("Failed to update movie", slog.String("error", err.Error()))
		return err
	}

	log.Info("Movie update process executed successfully")
	return nil
}

func (m *MovieRepository) GetByID(ctx context.Context, movieID uuid.UUID, withPreload bool) (*domain.Movie, error) {
	log := slog.With(
		slog.String("repository", "movie"),
		slog.String("func", "GetByID"),
	)

	log.Info("Initializing get movie by ID process")

	var movie domain.Movie
	db := m.db.WithContext(ctx)

	if withPreload {
		db = db.Preload("Images").Preload("IndicativeRating")
	}

	if err := db.First(&movie, "id = ?", movieID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("No movie record found", slog.String("movieID", movieID.String()))
			return nil, nil
		}
		log.Error("Failed to get movie by ID", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("Get movie by ID process executed successfully")
	return &movie, nil
}

func (m *MovieRepository) getImageUploadKey() string {
	return "image_upload_queue"
}
