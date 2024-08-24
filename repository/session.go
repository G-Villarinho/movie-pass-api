package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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
		return nil, err
	}

	redisClient, err := do.Invoke[*redis.Client](i)
	if err != nil {
		return nil, err
	}

	return &sessionRepository{
		i:           i,
		db:          db,
		redisClient: redisClient,
	}, nil
}

func (s *sessionRepository) Create(ctx context.Context, session domain.Session) error {
	log := slog.With(
		slog.String("repository", "session"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing session creation process")
	sessionJSON, err := jsoniter.Marshal(session)
	if err != nil {
		log.Error("Failed to marshal session data", slog.String("error", err.Error()))
		return err
	}

	if err := s.redisClient.Set(ctx, s.getSessionKey(session.UserID.String()), sessionJSON, time.Duration(config.Env.SessionExp)*time.Hour).Err(); err != nil {
		log.Error("Failed to save token", slog.String("error", err.Error()))
		return err
	}
	log.Info("Create user session process executed succefully")
	return nil
}

func (s *sessionRepository) GetSession(ctx context.Context, userID uuid.UUID) (*domain.Session, error) {
	log := slog.With(
		slog.String("repository", "session"),
		slog.String("func", "GetSession"),
	)

	log.Info("Initializing get session process")

	sessionJSON, err := s.redisClient.Get(ctx, s.getSessionKey(userID.String())).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Warn("session not found")
			return nil, domain.ErrSessionNotFound
		}

		log.Error("Failed to retrieve session", slog.String("error", err.Error()))
		return nil, err
	}

	var session domain.Session
	if err := jsoniter.UnmarshalFromString(sessionJSON, &session); err != nil {
		log.Error("Failed to unmarshal user data", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("get session process executed succefully")
	return &session, nil
}

func (s *sessionRepository) getSessionKey(userID string) string {
	tokenKey := fmt.Sprintf("session_%s", userID)
	return tokenKey
}
