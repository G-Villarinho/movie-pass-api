package domain

import (
	"context"
	"errors"
	"mime/multipart"
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
	IndicativeRatingID uuid.UUID        `gorm:"column:id;type:char(36)"`
	UserID             uuid.UUID        `gorm:"column:userId;type:char(36);not null"`
	Title              string           `gorm:"column:title;type:varchar(255);not null;index"`
	Duration           int              `gorm:"column:duration;type:int;not null"`
	User               User             `gorm:"foreignKey:UserID"`
	IndicativeRating   IndicativeRating `gorm:"foreignKey:IndicativeRatingID"`
	CreatedAt          time.Time        `gorm:"column:createdAt;not null"`
	UpdatedAt          time.Time        `gorm:"column:updatedAt;default:NULL"`
	Images             []MovieImage     `gorm:"foreignKey:MovieID"`
}

func (Movie) TableName() string {
	return "Movie"
}

type MovieImage struct {
	ID        uuid.UUID `gorm:"column:id;type:char(36);primaryKey"`
	MovieID   uuid.UUID `gorm:"column:movieId;type:char(36);not null"`
	ImageURL  string    `gorm:"column:imageUrl;type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"column:createdAt;not null"`
	UpdatedAt time.Time `gorm:"column:updatedAt;default:NULL"`
}

func (MovieImage) TableName() string {
	return "MovieImage"
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

type MoviePayload struct {
	Images             []*multipart.FileHeader `json:"images" validate:"validateImages"`
	IndicativeRatingID uuid.UUID               `json:"indicativeRatingId" validate:"required,uuid"`
	Title              string                  `json:"title" validate:"required,min=1,max=255"`
	Duration           int                     `json:"duration" validate:"required,gt=0"`
}

type IndicativeRatingResponse struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	ImageURL    string    `json:"imageUrl"`
}

type MovieResponse struct {
	ID               uuid.UUID                `json:"id"`
	ImagesURL        []string                 `json:"imagesUrl,omitempty"`
	Title            string                   `json:"title"`
	Duration         int                      `json:"duration"`
	IndicativeRating IndicativeRatingResponse `json:"indicativeRating,omitempty"`
}

type MovieHandler interface {
	GetAllIndicativeRating(ctx echo.Context) error
	Create(ctx echo.Context) error
}

type MovieService interface {
	GetAllIndicativeRating(ctx context.Context) ([]*IndicativeRatingResponse, error)
	Create(ctx context.Context, payload MoviePayload) (*MovieResponse, error)
}

type MovieRepository interface {
	GetAllIndicativeRating(ctx context.Context) ([]*IndicativeRating, error)
	Create(ctx context.Context, movie Movie) error
	CreateMovieImage(ctx context.Context, movieImage []MovieImage) error
}

func (i *IndicativeRating) ToIndicativeRatingResponse() *IndicativeRatingResponse {
	return &IndicativeRatingResponse{
		ID:          i.ID,
		Description: i.Description,
		ImageURL:    i.ImageURL,
	}
}

func (m *Movie) ToMovieResponse() *MovieResponse {
	imagesURL := make([]string, len(m.Images))
	for i, image := range m.Images {
		imagesURL[i] = image.ImageURL
	}

	return &MovieResponse{
		ID:               m.ID,
		Title:            m.Title,
		Duration:         m.Duration,
		IndicativeRating: *m.IndicativeRating.ToIndicativeRatingResponse(),
		ImagesURL:        imagesURL,
	}
}
