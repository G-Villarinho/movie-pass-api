package domain

import (
	"time"

	"github.com/google/uuid"
)

type Cinema struct {
	ID        uuid.UUID `gorm:"column:id;type:char(36);primaryKey"`
	Name      string    `gorm:"column:name;type:varchar(255);not null"`
	Location  string    `gorm:"column:name;type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"column:createdAt;not null"`
	UpdatedAt time.Time `gorm:"column:updatedAt;default:NULL"`
}

func (Cinema) TableName() string {
	return "Cinema"
}
