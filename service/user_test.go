package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/GSVillas/movie-pass-api/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Create_WhenUserAlreadyExistsByEmail_ShouldReturnErrEmailAlreadyRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mock.NewMockUserRepository(ctrl)
	userService := &userService{
		userRepository: userRepositoryMock,
	}

	payload := &domain.UserPayload{
		FirstName:       "Test",
		LastName:        "Doe",
		Email:           "test@example.com",
		ConfirmEmail:    "test@example.com",
		Password:        "Str0ngP@ssw0rd!",
		ConfirmPassword: "Str0ngP@ssw0rd!",
		BirthDate:       time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	existingUser := &domain.User{}
	userRepositoryMock.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(existingUser, nil)
	err := userService.Create(context.Background(), *payload)

	assert.ErrorIs(t, err, domain.ErrEmailAlreadyRegister)
}

func TestUserService_Create_WhenGetUserByEmailFails_ShouldReturnErrGetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mock.NewMockUserRepository(ctrl)
	userService := &userService{
		userRepository: userRepositoryMock,
	}

	payload := domain.UserPayload{
		FirstName:       "Test",
		LastName:        "Doe",
		Email:           "test@example.com",
		ConfirmEmail:    "test@example.com",
		Password:        "Str0ngP@ssw0rd!",
		ConfirmPassword: "Str0ngP@ssw0rd!",
		BirthDate:       time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	repError := errors.New("Table 'User' not found")
	userRepositoryMock.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, repError)

	err := userService.Create(context.Background(), payload)

	assert.ErrorIs(t, err, domain.ErrGetUserByEmail)
}

func TestUserService_Create_WhenSuccess_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mock.NewMockUserRepository(ctrl)

	userService := &userService{
		userRepository: userRepositoryMock,
	}

	payload := &domain.UserPayload{
		FirstName:       "Test",
		LastName:        "Doe",
		Email:           "test@example.com",
		ConfirmEmail:    "test@example.com",
		Password:        "Str0ngP@ssw0rd!",
		ConfirmPassword: "Str0ngP@ssw0rd!",
		BirthDate:       time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	userRepositoryMock.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, nil)

	userRepositoryMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	err := userService.Create(context.Background(), *payload)

	assert.NoError(t, err)
}

func TestUserService_Create_WhenHashingPasswordFails_ShouldReturnErrHashingPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mock.NewMockUserRepository(ctrl)

	userService := &userService{
		userRepository: userRepositoryMock,
	}

	payload := &domain.UserPayload{
		FirstName:       "Test",
		LastName:        "Doe",
		Email:           "test@example.com",
		ConfirmEmail:    "test@example.com",
		Password:        "pneumoultramicroscopiosilicovulcanocolioticopneumoultramicroscopiosilicovulcanocolioticopneumoultramicroscopiosilicovulcanocolioticopneumoultramicroscopiosilicovulcanocolioticopneumoultramicroscopiosilicovulcanocolioticopneumoultramicroscopiosilicovulcanocolioticopneumoultramicroscopiosilicovulcanocolioticopneumoultramicroscopiosilicovulcanocoliotico",
		ConfirmPassword: "pneumoultramicroscopiosilicovulcanocolioticopneumoultramicroscopiosilicovulcanocolioticopneumoultramicroscopiosilicovulcanocolioticopneumoultramicroscopiosilicovulcanocolioticopneumoultramicroscopiosilicovulcanocolioticopneumoultramicroscopiosilicovulcanocolioticopneumoultramicroscopiosilicovulcanocolioticopneumoultramicroscopiosilicovulcanocoliotico",
	}

	userRepositoryMock.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, nil)

	err := userService.Create(context.Background(), *payload)

	assert.ErrorIs(t, err, domain.ErrHashingPassword)
}

func TestUserService_Create_WhenCreateUserFails_ShouldReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mock.NewMockUserRepository(ctrl)

	userService := &userService{
		userRepository: userRepositoryMock,
	}

	payload := &domain.UserPayload{
		FirstName:       "Test",
		LastName:        "Doe",
		Email:           "test@example.com",
		ConfirmEmail:    "test@example.com",
		Password:        "Str0ngP@ssw0rd!",
		ConfirmPassword: "Str0ngP@ssw0rd!",
		BirthDate:       time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	userRepositoryMock.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, nil)

	repError := errors.New("Table 'User' not found")
	userRepositoryMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(repError)

	err := userService.Create(context.Background(), *payload)

	assert.ErrorIs(t, err, domain.ErrCreateUser)
}

func TestUserService_SignIn_WhenGetUserByEmailFails_ShouldReturnNiilAndErrGetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mock.NewMockUserRepository(ctrl)
	sessionServiceMock := mock.NewMockSessionService(ctrl)
	userService := &userService{
		userRepository: userRepositoryMock,
		sessionService: sessionServiceMock,
	}

	payload := &domain.SignInPayload{
		Email:    "test@example.com",
		Password: "Str0ngP@ssw0rd!",
	}

	repError := errors.New("Table 'User' not found")
	userRepositoryMock.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, repError)

	response, err := userService.SignIn(context.Background(), *payload)

	assert.ErrorIs(t, err, domain.ErrGetUserByEmail)
	assert.Nil(t, response)
}

func TestUserService_SignIn_WhenUserNotFound_ShouldReturnNilAndErrUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mock.NewMockUserRepository(ctrl)
	sessionService := mock.NewMockSessionService(ctrl)

	userService := &userService{
		userRepository: userRepositoryMock,
		sessionService: sessionService,
	}

	payload := &domain.SignInPayload{
		Email:    "test@example.com",
		Password: "Str0ngP@ssw0rd!",
	}

	userRepositoryMock.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(nil, nil)

	response, err := userService.SignIn(context.Background(), *payload)

	assert.ErrorIs(t, err, domain.ErrUserNotFound)
	assert.Nil(t, response)
}

func TestUserService_SignIn_WhenSuccess_ShouldReturnSignInResponseAndNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mock.NewMockUserRepository(ctrl)
	sessionServiceMock := mock.NewMockSessionService(ctrl)

	userService := &userService{
		userRepository: userRepositoryMock,
		sessionService: sessionServiceMock,
	}

	payload := &domain.SignInPayload{
		Email:    "test@example.com",
		Password: "Str0ngP@ssw0rd!",
	}

	user := &domain.User{
		Email:        payload.Email,
		PasswordHash: "$2a$10$gUWyf9ESKUNGxIAByzPCdOP9UMLLhC039R5jGNivSPhQJFNl4P0OC",
	}

	userRepositoryMock.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(user, nil)
	sessionServiceMock.EXPECT().Create(gomock.Any(), user).Return("validtoken", nil)

	response, err := userService.SignIn(context.Background(), *payload)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "validtoken", response.Token)
}

func TestUserService_SignIn_WhenCreateSessionFails_ShouldReturnErrCreateSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mock.NewMockUserRepository(ctrl)
	sessionServiceMock := mock.NewMockSessionService(ctrl)

	userService := &userService{
		userRepository: userRepositoryMock,
		sessionService: sessionServiceMock,
	}

	payload := &domain.SignInPayload{
		Email:    "test@example.com",
		Password: "Str0ngP@ssw0rd!",
	}

	user := &domain.User{
		Email:        payload.Email,
		PasswordHash: "$2a$10$gUWyf9ESKUNGxIAByzPCdOP9UMLLhC039R5jGNivSPhQJFNl4P0OC",
	}

	userRepositoryMock.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(user, nil)

	sessionServiceMock.EXPECT().Create(gomock.Any(), user).Return("", errors.New("session error"))

	response, err := userService.SignIn(context.Background(), *payload)

	assert.ErrorIs(t, err, domain.ErrCreateSession)
	assert.Nil(t, response)
}

func TestUserService_SignIn_WhenPasswordIsInvalid_ShouldReturnErrInvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mock.NewMockUserRepository(ctrl)
	sessionServiceMock := mock.NewMockSessionService(ctrl)

	userService := &userService{
		userRepository: userRepositoryMock,
		sessionService: sessionServiceMock,
	}

	payload := &domain.SignInPayload{
		Email:    "test@example.com",
		Password: "Teste@123",
	}

	user := &domain.User{
		Email:        payload.Email,
		PasswordHash: "wrong_password",
	}

	userRepositoryMock.EXPECT().GetByEmail(gomock.Any(), payload.Email).Return(user, nil)

	_, err := userService.SignIn(context.Background(), *payload)

	assert.ErrorIs(t, err, domain.ErrInvalidPassword)
}
