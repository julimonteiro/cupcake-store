package repository

import "github.com/julimonteiro/cupcake-store/internal/models"


type CupcakeRepositoryInterface interface {
	Create(cupcake *models.Cupcake) error
	FindByID(id uint) (*models.Cupcake, error)
	FindAll() ([]models.Cupcake, error)
}

