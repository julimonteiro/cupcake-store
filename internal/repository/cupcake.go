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

func (r *CupcakeRepository) FindByID(id uint) (*models.Cupcake, error) {
	var cupcake models.Cupcake
	err := r.db.First(&cupcake, id).Error
	if err != nil {
		return nil, err
	}
	return &cupcake, nil
}

func (r *CupcakeRepository) FindAll() ([]models.Cupcake, error) {
	var cupcakes []models.Cupcake
	err := r.db.Find(&cupcakes).Error
	return cupcakes, err
}

func (r *CupcakeRepository) Update(cupcake *models.Cupcake) error {
	return r.db.Save(cupcake).Error
}

func (r *CupcakeRepository) Delete(id uint) error {
	return r.db.Delete(&models.Cupcake{}, id).Error
}