package repository

import (
	"database/sql"
	"fmt"

	"github.com/vadhe/api-category/internal/transaction/domain"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) FindAll() ([]domain.Transaction, error) {
	var query = `SELECT transactions.id, transactions.total_amount, transactions.created_at FROM transactions`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var transaction domain.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.TotalAmount, &transaction.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *TransactionRepository) CreateTransaction(tx *sql.Tx, items *domain.Transaction) (*domain.Transaction, error) {
	var transactionID int

	err := tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", items.TotalAmount).Scan(&transactionID)
	if err != nil {
		fmt.Println(err)

		return nil, err
	}

	for i := range items.Details {
		items.Details[i].TransactionID = transactionID
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
			transactionID, items.Details[i].ProductID, items.Details[i].Quantity, items.Details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &domain.Transaction{
		ID:          transactionID,
		TotalAmount: items.TotalAmount,
		Details:     items.Details,
	}, nil
}
