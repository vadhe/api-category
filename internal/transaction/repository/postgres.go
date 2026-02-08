package repository

import (
	"database/sql"

	"github.com/vadhe/api-category/internal/transaction/domain"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) FindAll() ([]domain.Transaction, error) {
	var query = `
 SELECT  transactions.id, transactions.total_amount, transactions.created_at,
  transaction_details.transaction_id,
   transaction_details.id,
 transaction_details.product_id, transaction_details.quantity, transaction_details.subtotal
 FROM transactions
INNER JOIN transaction_details ON transactions.id = transaction_details.transaction_id
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction
	a := make(map[int]*domain.Transaction)
	for rows.Next() {
		var transaction domain.Transaction
		var transactionDetails domain.TransactionDetail
		if err := rows.Scan(&transaction.ID, &transaction.TotalAmount, &transaction.CreatedAt,
			&transactionDetails.TransactionID,
			&transactionDetails.ID,
			&transactionDetails.ProductID, &transactionDetails.Quantity,
			&transactionDetails.Subtotal,
		); err != nil {
			return nil, err
		}
		if a[transaction.ID] == nil {
			a[transaction.ID] = &transaction
		}
		a[transaction.ID].Details = append(a[transaction.ID].Details, transactionDetails)
	}
	for _, v := range a {
		transactions = append(transactions, *v)
	}
	return transactions, nil
}

func (r *TransactionRepository) InsertTransaction(
	tx *sql.Tx,
	totalAmount int,
) (int, error) {

	var id int
	err := tx.QueryRow(
		"INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id",
		totalAmount,
	).Scan(&id)

	return id, err
}

func (r *TransactionRepository) InsertTransactionDetails(
	tx *sql.Tx,
	details []domain.TransactionDetail,
) error {

	for i := range details {
		err := tx.QueryRow(
			`INSERT INTO transaction_details
			 (transaction_id, product_id, quantity, subtotal)
			 VALUES ($1, $2, $3, $4)
			 RETURNING id`,
			details[i].TransactionID,
			details[i].ProductID,
			details[i].Quantity,
			details[i].Subtotal,
		).Scan(&details[i].ID)

		if err != nil {
			return err
		}
	}
	return nil
}

func (r *TransactionRepository) GetReport() (*domain.SalesSummary, error) {
	var res domain.SalesSummary
	var productID int

	if err := r.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0)
		FROM transactions
		WHERE created_at >= CURRENT_DATE
		  AND created_at < CURRENT_DATE + INTERVAL '1 day'
	`).Scan(&res.TotalRevenue); err != nil {
		return nil, err
	}

	if err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM transactions
		WHERE created_at >= CURRENT_DATE
		  AND created_at < CURRENT_DATE + INTERVAL '1 day'
	`).Scan(&res.TotalTransactions); err != nil {
		return nil, err
	}

	if err := r.db.QueryRow(`
		SELECT td.product_id, SUM(td.quantity)
		FROM transaction_details td
		JOIN transactions t ON t.id = td.transaction_id
		WHERE t.created_at >= CURRENT_DATE
		  AND t.created_at < CURRENT_DATE + INTERVAL '1 day'
		GROUP BY td.product_id
		ORDER BY SUM(td.quantity) DESC
		LIMIT 1
	`).Scan(&productID, &res.MostSoldProduct.QtySold); err != nil {
		return nil, err
	}

	if err := r.db.QueryRow(`
		SELECT name FROM products WHERE id = $1
	`, productID).Scan(&res.MostSoldProduct.Name); err != nil {
		return nil, err
	}

	return &res, nil
}
