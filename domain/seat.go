package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Seat struct {
	ID             uuid.UUID  `gorm:"column:id;type:char(36);primaryKey"`
	CinemaRoomID   uuid.UUID  `gorm:"column:cinemaRoomId;type:char(36);not null"`
	CinemaRoom     CinemaRoom `gorm:"foreignKey:CinemaRoomID"`
	SeatIdentifier string     `gorm:"column:seatIdentifier;type:char(5);not null"`
	CreatedAt      time.Time  `gorm:"column:createdAt;not null"`
	UpdatedAt      time.Time  `gorm:"column:updatedAt;default:NULL"`
}

type SeatPayload struct {
	SeatIdentifier string
}

func (Seat) TableName() string {
	return "Seat"
}

type SeatRepository interface {
	Create(ctx context.Context, seat Seat) error
	GetByID(ctx context.Context, seatID uuid.UUID) (*Seat, error)
	GetAll(ctx context.Context) ([]Seat, error)
	Delete(ctx context.Context, seatID uuid.UUID) error
}
