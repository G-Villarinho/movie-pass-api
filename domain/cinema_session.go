package domain

import (
	"time"

	"github.com/google/uuid"
)

type CinemaSession struct {
	ID           uuid.UUID  `gorm:"column:id;type:char(36);primaryKey"`
	CinemaRoomID uuid.UUID  `gorm:"column:cinemaRoomId;type:char(36);not null"`
	CinemaRoom   CinemaRoom `gorm:"foreignKey:CinemaRoomID"`
	MovieID      uuid.UUID  `gorm:"column:MovieId;type:char(36);not null"`
	Movie        CinemaRoom `gorm:"foreignKey:MovieID"`
	StartTime    time.Time  `gorm:"column:startTime;type:time;not null"`
	EndTime      time.Time  `gorm:"column:endTime;type:time;not null"`
	CreatedAt    time.Time  `gorm:"column:createdAt;not null"`
	UpdatedAt    time.Time  `gorm:"column:updatedAt;default:NULL"`
}

func (CinemaSession) TableName() string {
	return "CinemaSession"
}
