package domain

import (
	"time"

	"github.com/google/uuid"
)

type CinemaRoom struct {
	ID        uuid.UUID `gorm:"column:id;type:char(36);primaryKey"`
	Name      string    `gorm:"column:name;type:varchar(255);not null"`
	SeatCount int       `gorm:"column:seatCount;type:int;not null"`
	CinemaID  uuid.UUID `gorm:"column:cinemaId;type:char(36);not null"`
	Cinema    Cinema    `gorm:"foreignKey:CinemaID"`
	Rows      int       `gorm:"column:rows;type:int;not null"`
	Collumns  int       `gorm:"column:collumns;type:int;not null"`
	CreatedAt time.Time `gorm:"column:createdAt;not null"`
	UpdatedAt time.Time `gorm:"column:updatedAt;default:NULL"`
}

func (CinemaRoom) TableName() string {
	return "CinemaRoom"
}
