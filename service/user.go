package service

import (
	"context"
	"fmt"

	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/GSVillas/movie-pass-api/secure"
	"github.com/samber/do"
)

type userService struct {
	i              *do.Injector
	userRepository domain.UserRepository
	sessionService domain.SessionService
}

func NewUserService(i *do.Injector) (domain.UserService, error) {
	userRepository, err := do.Invoke[domain.UserRepository](i)
	if err != nil {
		return nil, fmt.Errorf("error to initialize UserRepository: %w", err)
	}

	sessionService, err := do.Invoke[domain.SessionService](i)
	if err != nil {
		return nil, fmt.Errorf("error to initialize SessionService: %w", err)
	}

	return &userService{
		i:              i,
		userRepository: userRepository,
		sessionService: sessionService,
	}, nil
}

func (u *userService) Create(ctx context.Context, payload domain.UserPayload) error {
	user, err := u.userRepository.GetByEmail(ctx, payload.Email)
	if err != nil {
		return fmt.Errorf("error to retrieve user by email %s: %w", payload.Email, err)
	}

	if user != nil {
		return domain.ErrEmailAlreadyRegister
	}

	passwordHash, err := secure.HashPassword(payload.Password)
	if err != nil {
		return fmt.Errorf("error to hash password for email %s: %w", payload.Email, err)
	}

	user = payload.ToUser(string(passwordHash))

	if err := u.userRepository.Create(ctx, *user); err != nil {
		return fmt.Errorf("error to create user with email %s: %w", payload.Email, err)
	}

	return nil
}

func (u *userService) SignIn(ctx context.Context, payload domain.SignInPayload) (*domain.SignInResponse, error) {
	user, err := u.userRepository.GetByEmail(ctx, payload.Email)
	if err != nil {
		return nil, fmt.Errorf("error to retrieve user by email %s: %w", payload.Email, err)
	}

	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	if err := secure.CheckPassword(user.PasswordHash, payload.Password); err != nil {
		return nil, domain.ErrInvalidPassword
	}

	token, err := u.sessionService.Create(ctx, *user)
	if err != nil {
		return nil, fmt.Errorf("error to create session for user ID %s: %w", user.ID, err)
	}

	return &domain.SignInResponse{
		Token: token,
	}, nil
}
