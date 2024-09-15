package service

import (
	"context"
	"fmt"

	"github.com/GSVillas/movie-pass-api/config"
	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/golang-jwt/jwt"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/do"
)

type sessionService struct {
	i                 *do.Injector
	sessionRepository domain.SessionRepository
}

func NewSessionService(i *do.Injector) (domain.SessionService, error) {
	sessionRepository, err := do.Invoke[domain.SessionRepository](i)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SessionRepository: %w", err)
	}

	return &sessionService{
		i:                 i,
		sessionRepository: sessionRepository,
	}, nil
}

func (s *sessionService) Create(ctx context.Context, user domain.User) (string, error) {
	token, err := s.createToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to create token for user ID %s: %w", user.ID, err)
	}

	session := &domain.Session{
		Token:     token,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UserID:    user.ID,
		Email:     user.Email,
	}

	if err := s.sessionRepository.Create(ctx, *session); err != nil {
		return "", fmt.Errorf("failed to create session for user ID %s: %w", user.ID, err)
	}

	return token, nil
}

func (s *sessionService) GetSession(ctx context.Context, token string) (*domain.Session, error) {
	sessionToken, err := s.extractSessionFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to extract session from token: %w", err)
	}

	session, err := s.sessionRepository.GetSession(ctx, sessionToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session for user ID %s: %w", sessionToken.UserID, err)
	}

	if session == nil {
		return nil, domain.ErrSessionNotFound
	}

	if token != session.Token {
		return nil, domain.ErrSessionMismatch
	}

	return session, nil
}

func (s *sessionService) createToken(user domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"moviePassId": user.ID,
		"firstName":   user.FirstName,
		"lastName":    user.LastName,
		"email":       user.Email,
	})

	tokenString, err := token.SignedString(config.Env.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token for user ID %s: %w", user.ID, err)
	}

	return tokenString, nil
}

func (s *sessionService) extractSessionFromToken(tokenString string) (*domain.Session, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, domain.ErrorUnexpectedMethod
		}
		return config.Env.PublicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, domain.ErrTokenInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domain.ErrTokenInvalid
	}

	sessionJSON, err := jsoniter.Marshal(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal claims into JSON: %w", err)
	}

	var session domain.Session
	if err := jsoniter.Unmarshal(sessionJSON, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session from JSON: %w", err)
	}

	return &session, nil
}
