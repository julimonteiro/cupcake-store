package repository

import (
	"github.com/julimonteiro/cupcake-store/internal/models"
	"gorm.io/gorm"
)

type CupcakeRepository struct {
	db *gorm.DB
}

var _ CupcakeRepositoryInterface = (*CupcakeRepository)(nil)

func NewCupcakeRepository(db *gorm.DB) *CupcakeRepository {
	return &CupcakeRepository{db: db}
}

func (r *CupcakeRepository) Create(cupcake *models.Cupcake) error {
	return r.db.Create(cupcake).Error
}