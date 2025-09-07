package service

import (
	"errors"
	"strings"

	"github.com/julimonteiro/cupcake-store/internal/models"
	"github.com/julimonteiro/cupcake-store/internal/repository"
)

type CupcakeService struct {
	repo repository.CupcakeRepositoryInterface
}

func NewCupcakeService(repo repository.CupcakeRepositoryInterface) *CupcakeService {
	return &CupcakeService{repo: repo}
}

func (s *CupcakeService) CreateCupcake(req *models.CreateCupcakeRequest) (*models.Cupcake, error) {
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	cupcake := &models.Cupcake{
		Name:        strings.TrimSpace(req.Name),
		Flavor:      strings.TrimSpace(req.Flavor),
		PriceCents:  req.PriceCents,
		IsAvailable: true,
	}

	if err := s.repo.Create(cupcake); err != nil {
		return nil, err
	}

	return cupcake, nil
}

func (s *CupcakeService) GetCupcake(id uint) (*models.Cupcake, error) {
	cupcake, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return cupcake, nil
}

func (s *CupcakeService) GetAllCupcakes() ([]models.Cupcake, error) {
	return s.repo.FindAll()
}

func (s *CupcakeService) UpdateCupcake(id uint, req *models.UpdateCupcakeRequest) (*models.Cupcake, error) {
	cupcake, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if len(name) < 2 {
			return nil, errors.New("name must have at least 2 characters")
		}
		cupcake.Name = name
	}

	if req.Flavor != nil {
		cupcake.Flavor = strings.TrimSpace(*req.Flavor)
	}

	if req.PriceCents != nil {
		if *req.PriceCents <= 0 {
			return nil, errors.New("price must be greater than zero")
		}
		cupcake.PriceCents = *req.PriceCents
	}

	if req.IsAvailable != nil {
		cupcake.IsAvailable = *req.IsAvailable
	}

	if err := s.repo.Update(cupcake); err != nil {
		return nil, err
	}

	return cupcake, nil
}

func (s *CupcakeService) DeleteCupcake(id uint) error {
	return s.repo.Delete(id)
}

func (s *CupcakeService) validateCreateRequest(req *models.CreateCupcakeRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}

	if len(strings.TrimSpace(req.Name)) < 2 {
		return errors.New("name must have at least 2 characters")
	}

	if strings.TrimSpace(req.Flavor) == "" {
		return errors.New("flavor is required")
	}

	if req.PriceCents <= 0 {
		return errors.New("price must be greater than zero")
	}

	return nil
}
