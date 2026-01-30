package service

import (
	"errors"

	domain "github.com/vadhe/api-category/internal/product/domain"
)

func CreateProduct(name string, price int, stock int) (*domain.Product, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if price <= 0 {
		return nil, errors.New("price is required")
	}
	if stock <= 0 {
		return nil, errors.New("stock is required")
	}

	product := &domain.Product{
		Name:  name,
		Price: price,
		Stock: stock,
	}

	return product, nil
}

type Repository interface {
	FindAll() ([]domain.Product, error)
	FindByID(id int) (*domain.Product, error)
	Create(name string, price int, stock int) (*domain.Product, error)
	Update(id int, name string, price int, stock int) (*domain.Product, error)
	Delete(id int) error
}

type ProductService struct {
	repo Repository
}

func NewProductService(repo Repository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetProducts() ([]domain.Product, error) {
	return s.repo.FindAll()
}

func (s *ProductService) GetProductByID(id int) (*domain.Product, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}
func (s *ProductService) CreateProduct(name string, price int, stock int) (*domain.Product, error) {
	newProduct, err := CreateProduct(name, price, stock)
	if err != nil {
		return nil, err
	}

	product, err := s.repo.Create(newProduct.Name, newProduct.Price, newProduct.Stock)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) UpdateProduct(id int, name string, price int, stock int) (*domain.Product, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	newProduct, err := CreateProduct(name, price, stock)
	if err != nil {
		return nil, err
	}

	product, err := s.repo.Update(id, newProduct.Name, newProduct.Price, newProduct.Stock)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) DeleteProduct(id int) error {
	if id <= 0 {
		return errors.New("invalid id")
	}

	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
