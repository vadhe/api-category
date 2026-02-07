package service

import (
	"errors"

	domain "github.com/vadhe/api-category/internal/category/domain"
)

func CreateCategory(category *domain.Category) (*domain.Category, error) {
	if category.Name == "" {
		return nil, errors.New("name is required")
	}
	if category.Description == "" {
		return nil, errors.New("description is required")
	}

	res := &domain.Category{
		Name:        category.Name,
		Description: category.Description,
	}

	return res, nil
}

type Repository interface {
	FindAll() ([]domain.Category, error)
	FindByID(id int) (*domain.Category, error)
	Create(category *domain.Category) (*domain.Category, error)
	Update(id int, category *domain.Category) (*domain.Category, error)
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
func (s *CategoryService) CreateCategory(category *domain.Category) (*domain.Category, error) {
	newCategory, err := CreateCategory(category)
	if err != nil {
		return nil, err
	}

	res, err := s.repo.Create(newCategory)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *CategoryService) UpdateCategory(id int, category *domain.Category) (*domain.Category, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	newCategory, err := CreateCategory(category)
	if err != nil {
		return nil, err
	}

	res, err := s.repo.Update(id, newCategory)
	if err != nil {
		return nil, err
	}
	return res, nil
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
