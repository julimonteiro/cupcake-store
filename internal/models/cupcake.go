package models

import "time"

type Cupcake struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"not null;size:100"`
	Flavor      string    `json:"flavor" gorm:"not null;size:100"`
	PriceCents  int       `json:"price_cents" gorm:"not null"`
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Cupcake) TableName() string {
	return "cupcakes"
}

type CreateCupcakeRequest struct {
	Name       string `json:"name" validate:"required,min=2"`
	Flavor     string `json:"flavor" validate:"required"`
	PriceCents int    `json:"price_cents" validate:"required,gt=0"`
}

type UpdateCupcakeRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2"`
	Flavor      *string `json:"flavor,omitempty" validate:"omitempty"`
	PriceCents  *int    `json:"price_cents,omitempty" validate:"omitempty,gt=0"`
	IsAvailable *bool   `json:"is_available,omitempty"`
}
