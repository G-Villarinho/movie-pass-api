package repository

import (
	"context"
	"errors"

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
	if err := u.db.WithContext(ctx).Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user *domain.User
	if err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}
