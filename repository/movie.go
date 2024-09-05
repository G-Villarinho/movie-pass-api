package repository

import (
	"context"
	"errors"
	"log/slog"

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
	log := slog.With(
		slog.String("repository", "movie"),
		slog.String("func", "GetNextUploadImageTaskFromQueue"),
	)

	log.Info("Retrieving next image upload task from Redis queue")

	data, err := m.redisClient.LPop(ctx, m.getImageUploadKey()).Result()
	if err != nil {
		if err == redis.Nil {
			log.Warn("No tasks found in Redis queue")
			return nil, nil
		}
		log.Error("Failed to retrieve task from Redis queue", slog.String("error", err.Error()))
		return nil, err
	}

	var task domain.MovieImageUploadTask
	if err := jsoniter.Unmarshal([]byte(data), &task); err != nil {
		log.Error("Failed to deserialize task from Redis", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("Successfully retrieved image upload task from Redis queue")
	return &task, nil
}

func (m *MovieRepository) getImageUploadKey() string {
	return "image_upload_queue"
}
