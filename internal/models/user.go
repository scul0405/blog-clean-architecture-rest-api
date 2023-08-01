package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// User full model
type User struct {
	UserID      uuid.UUID  `json:"user_id" db:"user_id" validate:"omitempty"`
	FirstName   string     `json:"first_name" db:"first_name" validate:"required,lte=30"`
	LastName    string     `json:"last_name" db:"last_name" validate:"required,lte=30"`
	Email       string     `json:"email,omitempty" db:"email" validate:"omitempty,lte=60,email"`
	Password    string     `json:"password,omitempty" db:"password" validate:"omitempty,required,gte=6"`
	Role        *string    `json:"role,omitempty" db:"role" validate:"omitempty,lte=10"`
	About       *string    `json:"about,omitempty" db:"about" validate:"omitempty,lte=1024"`
	Avatar      *string    `json:"avatar,omitempty" db:"avatar" validate:"omitempty,lte=512,url"`
	PhoneNumber *string    `json:"phone_number,omitempty" db:"phone_number" validate:"omitempty,lte=20"`
	Address     *string    `json:"address,omitempty" db:"address" validate:"omitempty,lte=250"`
	City        *string    `json:"city,omitempty" db:"city" validate:"omitempty,lte=24"`
	Country     *string    `json:"country,omitempty" db:"country" validate:"omitempty,lte=24"`
	Gender      *string    `json:"gender,omitempty" db:"gender" validate:"omitempty,lte=10"`
	Postcode    *int       `json:"postcode,omitempty" db:"postcode" validate:"omitempty"`
	Birthday    *time.Time `json:"birthday,omitempty" db:"birthday" validate:"omitempty,lte=10"`
	CreatedAt   time.Time  `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at,omitempty" db:"updated_at"`
	LoginDate   time.Time  `json:"login_date" db:"login_date"`
}

// PrepareCreate prepare for register
func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	if u.PhoneNumber != nil {
		*u.PhoneNumber = strings.TrimSpace(*u.PhoneNumber)
	}

	if u.Role != nil {
		*u.Role = strings.ToLower(strings.TrimSpace(*u.Role))
	}
	return nil
}
