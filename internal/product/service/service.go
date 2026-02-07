package service

import (
	"database/sql"
	"errors"

	domain "github.com/vadhe/api-category/internal/product/domain"
)

func CreateProduct(product *domain.Product) (*domain.Product, error) {
	if product.Name == "" {
		return nil, errors.New("name is required")
	}
	if product.Price <= 0 {
		return nil, errors.New("price is required")
	}
	if product.Stock <= 0 {
		return nil, errors.New("stock is required")
	}
	if product.CategoryId <= 0 {
		return nil, errors.New("categoryId is required")
	}
	res := &domain.Product{
		Name:       product.Name,
		Price:      product.Price,
		Stock:      product.Stock,
		CategoryId: product.CategoryId,
	}

	return res, nil
}

type Repository interface {
	FindAll(name string) ([]domain.Product, error)
	FindByID(id int) (*domain.Product, error)
	Create(product *domain.Product) (*domain.Product, error)
	Update(id int, product *domain.Product) (*domain.Product, error)
	Delete(id int) error
	DecreaseStockTx(tx *sql.Tx, productID int, qty int) (*domain.Product, error)
}

type ProductService struct {
	repo Repository
}

func NewProductService(repo Repository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetProducts(name string) ([]domain.Product, error) {
	return s.repo.FindAll(name)
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
func (s *ProductService) CreateProduct(product *domain.Product) (*domain.Product, error) {
	newProduct, err := CreateProduct(product)
	if err != nil {
		return nil, err
	}

	res, err := s.repo.Create(newProduct)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ProductService) UpdateProduct(id int, product *domain.Product) (*domain.Product, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	newProduct, err := CreateProduct(product)
	if err != nil {
		return nil, err
	}

	res, err := s.repo.Update(id, newProduct)
	if err != nil {
		return nil, err
	}
	return res, nil
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
