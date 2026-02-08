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
	GetReport() (*domain.SalesSummary, error)
	InsertTransaction(tx *sql.Tx, totalAmount int) (int, error)
	InsertTransactionDetails(tx *sql.Tx, details []domain.TransactionDetail) error
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

func (s *TransactionService) Checkout(
	items []domain.CheckoutItem,
) (*domain.Transaction, error) {

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
		if _, err := s.product.DecreaseStockTx(tx, product.ID, newStock); err != nil {
			return nil, err
		}

		subtotal := product.Price * item.Quantity
		transaction.TotalAmount += subtotal
		transaction.Details = append(transaction.Details, domain.TransactionDetail{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Subtotal:  subtotal,
		})
	}

	transactionID, err := s.repo.InsertTransaction(tx, transaction.TotalAmount)
	if err != nil {
		return nil, err
	}

	for i := range transaction.Details {
		transaction.Details[i].TransactionID = transactionID
	}

	if err := s.repo.InsertTransactionDetails(tx, transaction.Details); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	transaction.ID = transactionID
	
	return transaction, nil
}

func (s *TransactionService) GetReport() (*domain.SalesSummary, error) {
	return s.repo.GetReport()
}
