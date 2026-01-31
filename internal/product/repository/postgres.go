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
	SELECT products.id, products.name, products.price, products.stock, categories.name, categories.id
    FROM Products
    INNER JOIN categories ON products.category_id  = categories.id;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryName, &product.CategoryId); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) FindByID(id int) (*domain.Product, error) {
	row := r.db.QueryRow(`
		SELECT products.id, products.name, products.price, products.stock, categories.name, categories.id
        FROM Products
        INNER JOIN categories ON products.category_id  = categories.id
        WHERE products.id = $1;
	`, id)

	var product domain.Product
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryName, &product.CategoryId)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Create(name string, price int, stock int, categoryId int) (*domain.Product, error) {
	row := r.db.QueryRow(`
		WITH inserted_product AS (
        INSERT INTO products (name, price, stock, category_id)
        VALUES ($1, $2, $3, $4)
        RETURNING id, name, price, stock, category_id
        )
	    SELECT
	        p.id, p.name, p.price, p.stock, p.category_id, c.name as category_name
	    FROM inserted_product p
	    INNER JOIN categories c ON p.category_id = c.id;
	`, name, price, stock, categoryId)

	var product domain.Product
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryId, &product.CategoryName)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Update(id int, name string, price int, stock int, categoryId int) (*domain.Product, error) {
	row := r.db.QueryRow(`
   		WITH updated_product AS (
        UPDATE products SET name = $2, price = $3, stock = $4, category_id = $5 WHERE id = $1
        RETURNING id, name, price, stock, category_id
           )
           SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name as category_name
           FROM updated_product p
           INNER JOIN categories c ON p.category_id = c.id;
	`, id, name, price, stock, categoryId)

	var product domain.Product
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryId, &product.CategoryName)
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
