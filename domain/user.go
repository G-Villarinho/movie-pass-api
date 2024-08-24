package domain

//go:generate mockgen -source=user.go -destination=../mock/user_mock.go -package=mock

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailAlreadyRegister  = errors.New("email already exists")
	ErrHashingPassword       = errors.New("failed to hash password")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrUserNotFoundInContext = errors.New("user not found in context")
	ErrCreateUser            = errors.New("create user fail")
	ErrGetUserByEmail        = errors.New("get user by email fail")
)

type User struct {
	ID           uuid.UUID `gorm:"column:id;type:char(36);primaryKey"`
	FirstName    string    `gorm:"column:firstName;type:varchar(255);not null"`
	LastName     string    `gorm:"column:lastName;type:varchar(255);not null"`
	Email        string    `gorm:"column:email;type:varchar(255);uniqueIndex;not null"`
	BirthDate    time.Time `gorm:"column:birthDate;type:date;not null"`
	PasswordHash string    `gorm:"column:passwordHash;type:varchar(255);not null"`
	CreatedAt    time.Time `gorm:"column:createdAt;type:date;not null"`
	UpdatedAt    time.Time `gorm:"column:updatedAt;type:date;default:NULL"`
}

func (User) TableName() string {
	return "User"
}

type UserPayload struct {
	FirstName       string    `json:"firstName" validate:"required,min=1,max=255"`
	LastName        string    `json:"lastName" validate:"required,min=1,max=255"`
	Email           string    `json:"email" validate:"required,email,max=255"`
	ConfirmEmail    string    `json:"confirmEmail" validate:"required,eqfield=Email"`
	Password        string    `json:"password,omitempty" validate:"required,max=255,strongpassword"`
	ConfirmPassword string    `json:"confirmPassword" validate:"required,eqfield=Password"`
	BirthDate       time.Time `json:"birthDate" validate:"required,nottooold,notfuturedate"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birthDate"`
}

type SignInPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

type UserHandler interface {
	Create(ctx echo.Context) error
	SignIn(ctx echo.Context) error
}

type UserService interface {
	Create(ctx context.Context, payload UserPayload) error
	SignIn(ctx context.Context, payload SignInPayload) (*SignInResponse, error)
}

type UserRepository interface {
	Create(ctx context.Context, user User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}

func (u *UserPayload) trim() {
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	u.ConfirmEmail = strings.TrimSpace(strings.ToLower(u.ConfirmEmail))
}

func (s *SignInPayload) trim() {
	s.Email = strings.TrimSpace(strings.ToLower(s.Email))
}

func (s *SignInPayload) Validate() ValidationErrors {
	s.trim()
	return ValidateStruct(s)
}

func (u *UserPayload) Validate() ValidationErrors {
	u.trim()
	return ValidateStruct(u)
}

func (u *UserPayload) ToUser(passwordHash string) *User {
	return &User{
		ID:           uuid.New(),
		FirstName:    u.FirstName,
		Email:        u.Email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
	}
}
