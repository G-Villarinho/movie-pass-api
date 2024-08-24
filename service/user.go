package service

import (
	"context"
	"log/slog"

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
		return nil, err
	}

	sessionService, err := do.Invoke[domain.SessionService](i)
	if err != nil {
		return nil, err
	}

	return &userService{
		i:              i,
		userRepository: userRepository,
		sessionService: sessionService,
	}, nil
}

func (u *userService) Create(ctx context.Context, payload domain.UserPayload) error {
	log := slog.With(
		slog.String("service", "user"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing user creation process")

	user, err := u.userRepository.GetByEmail(ctx, payload.Email)
	if err != nil {
		log.Error("Failed to get user by email", slog.String("error", err.Error()))
		return domain.ErrGetUserByEmail
	}

	if user != nil {
		log.Warn("There is already a user with this ", slog.String("email", payload.Email))
		return domain.ErrEmailAlreadyRegister
	}

	passwordHash, err := secure.HashPassword(payload.Password)
	if err != nil {
		log.Error("Failed to hash password", slog.String("error", err.Error()))
		return domain.ErrHashingPassword
	}

	user = payload.ToUser(string(passwordHash))

	if err := u.userRepository.Create(ctx, *user); err != nil {
		log.Error("Failed to create user", slog.String("error", err.Error()))
		return domain.ErrCreateUser
	}

	log.Info("User creation process executed successfully")
	return nil
}

func (u *userService) SignIn(ctx context.Context, payload domain.SignInPayload) (*domain.SignInResponse, error) {
	log := slog.With(
		slog.String("service", "user"),
		slog.String("func", "SignIn"),
	)

	log.Info("Initializing user creation process")

	user, err := u.userRepository.GetByEmail(ctx, payload.Email)
	if err != nil {
		log.Error("Failed to get user by email", slog.String("error", err.Error()))
		return nil, domain.ErrGetUserByEmail
	}

	if user == nil {
		log.Warn("User not found with this email", slog.String("email", payload.Email))
		return nil, domain.ErrUserNotFound
	}

	if err := secure.CheckPassword(user.PasswordHash, payload.Password); err != nil {
		log.Warn("The password entered is invalid for the ", slog.String("email", payload.Email))
		return nil, domain.ErrInvalidPassword
	}

	token, err := u.sessionService.Create(ctx, user)
	if err != nil {
		log.Error("Was not possible create the session for the user", slog.String("error", err.Error()))
		return nil, domain.ErrCreateSession
	}

	log.Info("user sign in process executed successfully")
	return &domain.SignInResponse{
		Token: token,
	}, nil
}
