package domain

import (
	"time"

	"github.com/google/uuid"
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
