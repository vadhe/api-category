package service

import (
	"errors"

	domain "github.com/vadhe/api-category/internal/category/domain"
)

func CreateCategory(name, description string) (*domain.Category, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if description == "" {
		return nil, errors.New("description is required")
	}

	category := &domain.Category{
		Name:        name,
		Description: description,
	}

	return category, nil
}

type Repository interface {
	FindAll() ([]domain.Category, error)
	FindByID(id int) (*domain.Category, error)
	Create(name, description string) (*domain.Category, error)
	Update(id int, name, description string) (*domain.Category, error)
	Delete(id int) error
}

type CategoryService struct {
	repo Repository
}

func NewCategoryService(repo Repository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetCategories() ([]domain.Category, error) {
	return s.repo.FindAll()
}

func (s *CategoryService) GetCategoryByID(id int) (*domain.Category, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	category, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}
	return category, nil
}
func (s *CategoryService) CreateCategory(name, description string) (*domain.Category, error) {
	newCategory, err := CreateCategory(name, description)
	if err != nil {
		return nil, err
	}

	category, err := s.repo.Create(newCategory.Name, newCategory.Description)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) UpdateCategory(id int, name, description string) (*domain.Category, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	newCategory, err := CreateCategory(name, description)
	if err != nil {
		return nil, err
	}

	category, err := s.repo.Update(id, newCategory.Name, newCategory.Description)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) DeleteCategory(id int) error {
	if id <= 0 {
		return errors.New("invalid id")
	}

	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
