package domain

import (
	"time"

	"github.com/google/uuid"
)

type SeatReservation struct {
	ID              uuid.UUID     `gorm:"column:id;type:char(36);primaryKey"`
	CinemaSessionID uuid.UUID     `gorm:"column:cinemaSessionId;type:char(36);not null"`
	CinemaSession   CinemaSession `gorm:"foreignKey:CinemaSessionID"`
	SeatID          uuid.UUID     `gorm:"column:SeatId;type:char(36);not null"`
	Seat            Seat          `gorm:"foreignKey:SeatID"`
	UserID          uuid.UUID     `gorm:"column:UserID;type:char(36);not null"`
	User            User          `gorm:"foreignKey:SeatID"`
	CreatedAt       time.Time     `gorm:"column:createdAt;not null"`
	UpdatedAt       time.Time     `gorm:"column:updatedAt;default:NULL"`
}

func (SeatReservation) TableName() string {
	return "SeatReservation"
}
