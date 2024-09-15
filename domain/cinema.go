package domain

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	ErrCreateCinema   = errors.New("error to create a new cinema")
	ErrGetCinema      = errors.New("error to get cinema")
	ErrCinemaNotFound = errors.New("cinema not found")
	ErrDeleteCinema   = errors.New("error to delete cinema")
)

type Cinema struct {
	ID        uuid.UUID `gorm:"column:id;type:char(36);primaryKey"`
	Name      string    `gorm:"column:name;type:varchar(255);not null"`
	Location  string    `gorm:"column:location;type:varchar(255);not null"`
	UserID    uuid.UUID `gorm:"column:userId;type:char(36);not null"`
	User      User      `gorm:"foreignKey:UserID"`
	CreatedAt time.Time `gorm:"column:createdAt;not null"`
	UpdatedAt time.Time `gorm:"column:updatedAt;default:NULL"`
}

func (Cinema) TableName() string {
	return "Cinema"
}

type CinemaPayload struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Location string `json:"location" validate:"required,min=1,max=255"`
}

type CinemaResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"createdAt"`
}

type CinemaHandler interface {
	Create(ctx echo.Context) error
	GetByID(ctx echo.Context) error
	GetAll(ctx echo.Context) error
	Delete(ctx echo.Context) error
}

type CinemaService interface {
	Create(ctx context.Context, payload CinemaPayload) (*CinemaResponse, error)
	GetByID(ctx context.Context, cinemaID uuid.UUID) (*CinemaResponse, error)
	GetAll(ctx context.Context) ([]CinemaResponse, error)
	Delete(ctx context.Context, cinemaID uuid.UUID) error
}

type CinemaRepository interface {
	Create(ctx context.Context, cinema Cinema) error
	GetByID(ctx context.Context, cinemaID uuid.UUID) (*Cinema, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]Cinema, error)
	Delete(ctx context.Context, cinemaID uuid.UUID) error
}

func (c *CinemaPayload) trim() {
	c.Name = strings.TrimSpace(c.Name)
	c.Location = strings.TrimSpace(c.Location)
}

func (c *CinemaPayload) Validate() ValidationErrors {
	c.trim()
	return ValidateStruct(c)
}

func (c *CinemaPayload) ToCinema(userID uuid.UUID) *Cinema {
	return &Cinema{
		ID:        uuid.New(),
		Name:      c.Name,
		Location:  c.Location,
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
	}
}

func (c *Cinema) ToCinemaResponse() *CinemaResponse {
	return &CinemaResponse{
		ID:        c.ID,
		Name:      c.Name,
		Location:  c.Location,
		CreatedAt: c.CreatedAt,
	}
}
