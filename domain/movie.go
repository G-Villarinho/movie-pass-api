package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	ErrGetAllIndicativeRating    = errors.New("failed to obtain all indicative ratings")
	ErrIndicativeRatingsNotFound = errors.New("indicative ratings not found")
)

type Movie struct {
	ID                 uuid.UUID        `gorm:"column:id;type:char(36);primaryKey"`
	Title              string           `gorm:"column:title;type:varchar(255);not null;index"`
	Duration           int              `gorm:"column:duration;type:int;not null"`
	IndicativeRatingID uuid.UUID        `gorm:"column:id;type:char(36)"`
	IndicativeRating   IndicativeRating `gorm:"foreignKey:IndicativeRatingID"`
	UserID             uuid.UUID        `gorm:"column:userId;type:char(36);not null"`
	User               User             `gorm:"foreignKey:UserID"`
	CreatedAt          time.Time        `gorm:"column:createdAt;not null"`
	UpdatedAt          time.Time        `gorm:"column:updatedAt;default:NULL"`
}

func (Movie) TableName() string {
	return "Movie"
}

type IndicativeRating struct {
	ID          uuid.UUID `gorm:"column:id;type:char(36);primaryKey"`
	Description string    `gorm:"column:description;type:char(4);not null;uniqueIndex"`
	ImageURL    string    `gorm:"column:imageUrl;type:varchar(255);not null"`
	CreatedAt   time.Time `gorm:"column:createdAt;not null"`
	UpdatedAt   time.Time `gorm:"column:updatedAt;default:NULL"`
}

func (IndicativeRating) TableName() string {
	return "IndicativeRating"
}

type IndicativeRatingResponse struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	ImageURL    string    `json:"imageUrl"`
}

type MovieHandler interface {
	GetAllIndicativeRating(ctx echo.Context) error
}

type MovieService interface {
	GetAllIndicativeRating(ctx context.Context) ([]*IndicativeRatingResponse, error)
}
type MovieRepository interface {
	GetAllIndicativeRating(ctx context.Context) ([]*IndicativeRating, error)
}

func (i *IndicativeRating) ToIndicativeRatingResponse() *IndicativeRatingResponse {
	return &IndicativeRatingResponse{
		ID:          i.ID,
		Description: i.Description,
		ImageURL:    i.ImageURL,
	}
}
