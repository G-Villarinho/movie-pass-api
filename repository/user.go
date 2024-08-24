package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/go-redis/redis/v8"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type userRepository struct {
	i           *do.Injector
	db          *gorm.DB
	redisClient *redis.Client
}

func NewUserRepository(i *do.Injector) (domain.UserRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}

	redisClient, err := do.Invoke[*redis.Client](i)
	if err != nil {
		return nil, err
	}

	return &userRepository{
		i:           i,
		db:          db,
		redisClient: redisClient,
	}, nil
}

func (u *userRepository) Create(ctx context.Context, user domain.User) error {
	log := slog.With(
		slog.String("repository", "user"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing user creation process")
	if err := u.db.WithContext(ctx).Create(&user).Error; err != nil {
		log.Error("Failed to create user", slog.String("error", err.Error()))
		return err
	}

	log.Info("Create user process executed successfully")
	return nil
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	log := slog.With(
		slog.String("repository", "user"),
		slog.String("func", "GetByEmail"),
	)

	log.Info("Initializing process of obtaining user by email")

	var user *domain.User
	if err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("User not found")
			return nil, nil
		}

		log.Error("Failed to get user by email", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("Process of obtaining user by email executed successfully")
	return user, nil
}
