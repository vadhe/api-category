package service

import (
	"database/sql"
	"errors"
	"time"

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
	GetReport(startDate time.Time, endDate time.Time) (*domain.SalesSummary, error)
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
	// dari fe map gabungin qty by id
	totalQty := make(map[int]int)
	for _, item := range items {
		if totalQty[item.ProductID] == 0 {
			totalQty[item.ProductID] = item.Quantity
		} else {
			totalQty[item.ProductID] += item.Quantity
		}
	}
	products := make(map[int]productDomain.Product)
	for productID := range totalQty {
		p, err := s.product.FindByID(productID)
		if err != nil {
			return nil, err
		}
		products[productID] = *p
	}

	for productID, qty := range totalQty {
		p := products[productID]

		if p.Stock < qty {
			return nil, errors.New("stock not enough")
		}

		if _, err := s.product.DecreaseStockTx(tx, productID, qty); err != nil {
			return nil, err
		}
	}

	transaction := &domain.Transaction{}
	for _, item := range items {
		p := products[item.ProductID]

		subtotal := p.Price * item.Quantity
		transaction.TotalAmount += subtotal

		transaction.Details = append(transaction.Details, domain.TransactionDetail{
			ProductID: item.ProductID,
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

func (s *TransactionService) GetReport(startDate time.Time, endDate time.Time) (*domain.SalesSummary, error) {
	return s.repo.GetReport(startDate, endDate)
}
