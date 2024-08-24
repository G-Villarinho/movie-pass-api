package domain

//go:generate mockgen -source=session.go -destination=../mock/session_mock.go -package=mock

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrTokenInvalid           = errors.New("invalid token")
	ErrSessionNotFound        = errors.New("token not found")
	ErrorUnexpectedMethod     = errors.New("unexpected signing method")
	ErrTokenNotFoundInContext = errors.New("token not found in context")
	ErrSessionMismatch        = errors.New("session icompatible for user requested")
	ErrCreateSession          = errors.New("create session fails")
	ErrCreateToken            = errors.New("create session fails")
)

type Session struct {
	Token     string    `json:"token"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	UserID    uuid.UUID `json:"MoviePassId"`
	Email     string    `json:"email"`
}

type SessionService interface {
	Create(ctx context.Context, user User) (string, error)
	GetSession(ctx context.Context, token string) (*Session, error)
}

type SessionRepository interface {
	Create(ctx context.Context, session Session) error
	GetSession(ctx context.Context, userID uuid.UUID) (*Session, error)
}
