package domain

import (
	"time"

	"github.com/google/uuid"
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
	FirstName       string    `json:"firstName" validate:"required,min=1,max=75"`
	LastName        string    `json:"lastName" validate:"required,min=1,max=75"`
	Email           string    `json:"email" validate:"required,email"`
	ConfirmEmail    string    `json:"confirmEmail" validate:"required,eqfield=Email"`
	Password        string    `json:"password,omitempty" validate:"required,max=255,strongpassword"`
	ConfirmPassword string    `json:"confirmPassword" validate:"required,eqfield=Password"`
	BirthDate       time.Time `json:"birthDate" validate:"required"`
}

type UserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName" validate:"required,min=1,max=75"`
	Email     string `json:"email" validate:"required,email"`
}

type SignInPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required"`
}

type SignInResponse struct {
	Token string `json:"token"`
}
