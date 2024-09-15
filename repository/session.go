package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/GSVillas/movie-pass-api/config"
	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type sessionRepository struct {
	i           *do.Injector
	db          *gorm.DB
	redisClient *redis.Client
}

func NewSessionRepository(i *do.Injector) (domain.SessionRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, fmt.Errorf("error to initialize DB connection: %w", err)
	}

	redisClient, err := do.Invoke[*redis.Client](i)
	if err != nil {
		return nil, fmt.Errorf("error to initialize Redis client: %w", err)
	}

	return &sessionRepository{
		i:           i,
		db:          db,
		redisClient: redisClient,
	}, nil
}

func (s *sessionRepository) Create(ctx context.Context, session domain.Session) error {
	sessionJSON, err := jsoniter.Marshal(session)
	if err != nil {
		return err
	}

	if err := s.redisClient.Set(ctx, s.getSessionKey(session.UserID.String()), sessionJSON, time.Duration(config.Env.SessionExp)*time.Hour).Err(); err != nil {
		return err
	}

	return nil
}

func (s *sessionRepository) GetSession(ctx context.Context, userID uuid.UUID) (*domain.Session, error) {
	sessionJSON, err := s.redisClient.Get(ctx, s.getSessionKey(userID.String())).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var session domain.Session
	if err := jsoniter.UnmarshalFromString(sessionJSON, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *sessionRepository) getSessionKey(userID string) string {
	return fmt.Sprintf("session_%s", userID)
}
