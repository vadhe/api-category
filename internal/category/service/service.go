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
