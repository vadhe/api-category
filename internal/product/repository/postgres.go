package repository

import (
	"database/sql"

	"github.com/vadhe/api-category/internal/product/domain"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindAll() ([]domain.Product, error) {
	rows, err := r.db.Query(`
		SELECT id, name, price, stock
		FROM products
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) FindByID(id int) (*domain.Product, error) {
	row := r.db.QueryRow(`
		SELECT id, name, price, stock FROM products WHERE id = $1;
	`, id)

	var product domain.Product
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Create(name string, price int, stock int) (*domain.Product, error) {
	row := r.db.QueryRow(`
		INSERT INTO products (name, price, stock)
		VALUES ($1, $2, $3)
		RETURNING id, name, price, stock;
	`, name, price, stock)

	var product domain.Product
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Update(id int, name string, price int, stock int) (*domain.Product, error) {
	row := r.db.QueryRow(`
		UPDATE products SET name = $2, price = $3, stock = $4 WHERE id = $1
		RETURNING id, name, price, stock;
	`, id, name, price, stock)

	var product domain.Product
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Delete(id int) error {
	err := r.db.QueryRow(`DELETE FROM products WHERE id = $1 RETURNING id`, id).Scan(&id)
	if err == sql.ErrNoRows {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}
