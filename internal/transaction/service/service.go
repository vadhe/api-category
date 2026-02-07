package service

import (
	"database/sql"
	"errors"

	productDomain "github.com/vadhe/api-category/internal/product/domain"
	"github.com/vadhe/api-category/internal/transaction/domain"
)

func CreateTransaction(totalAmount int) (*domain.Transaction, error) {
	if totalAmount <= 0 {
		return nil, errors.New("Total Amount is required")
	}

	product := &domain.Transaction{
		TotalAmount: totalAmount,
	}

	return product, nil
}

type ProductRepository interface {
	FindByID(id int) (*productDomain.Product, error)
	DecreaseStockTx(tx *sql.Tx, productID int, qty int) (*productDomain.Product, error)
}
type Repository interface {
	FindAll() ([]domain.Transaction, error)
	CreateTransaction(tx *sql.Tx, items *domain.Transaction) (*domain.Transaction, error)
}

type TransactionService struct {
	repo    Repository
	product ProductRepository
	db      *sql.DB
}

func NewTransactionService(repo Repository, product ProductRepository, db *sql.DB) *TransactionService {
	return &TransactionService{repo: repo, product: product, db: db}
}

func (s *TransactionService) GetTransactions() ([]domain.Transaction, error) {
	return s.repo.FindAll()
}
func (s *TransactionService) Checkout(items []domain.CheckoutItem) (*domain.Transaction, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	transaction := &domain.Transaction{}
	for _, item := range items {
		product, err := s.product.FindByID(item.ProductID)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, errors.New("product not found")
		}
		if product.Stock < item.Quantity {
			return nil, errors.New("stock not enough")
		}
		newStock := product.Stock - item.Quantity
		_, err = s.product.DecreaseStockTx(tx, product.ID, newStock)
		if err != nil {
			return nil, err
		}
		data := domain.TransactionDetail{
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			Subtotal:    product.Price * item.Quantity,
		}
		transaction.TotalAmount += product.Price * item.Quantity
		transaction.Details = append(transaction.Details, data)
	}

	return s.repo.CreateTransaction(tx, transaction)
}
