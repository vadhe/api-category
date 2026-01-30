package repository

import (
	"database/sql"

	"github.com/vadhe/api-category/internal/category/domain"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindAll() ([]domain.Category, error) {
	rows, err := r.db.Query(`
		SELECT id, name, description
		FROM categories
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (r *CategoryRepository) FindByID(id int) (*domain.Category, error) {
	row := r.db.QueryRow(`
		SELECT id, name, description FROM categories WHERE id = $1;
	`, id)

	var categorie domain.Category
	err := row.Scan(&categorie.ID, &categorie.Name, &categorie.Description)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &categorie, nil
}

func (r *CategoryRepository) Create(name, description string) (*domain.Category, error) {
	row := r.db.QueryRow(`
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description;
	`, name, description)

	var categorie domain.Category
	err := row.Scan(&categorie.ID, &categorie.Name, &categorie.Description)
	if err != nil {
		return nil, err
	}
	return &categorie, nil
}

func (r *CategoryRepository) Update(id int, name, description string) (*domain.Category, error) {
	row := r.db.QueryRow(`
		UPDATE categories SET name = $2, description = $3 WHERE id = $1
		RETURNING id, name, description;
	`, id, name, description)

	var categorie domain.Category
	err := row.Scan(&categorie.ID, &categorie.Name, &categorie.Description)
	if err != nil {
		return nil, err
	}
	return &categorie, nil
}

func (r *CategoryRepository) Delete(id int) error {
	err := r.db.QueryRow(`DELETE FROM categories WHERE id = $1 RETURNING id`, id).Scan(&id)
	if err == sql.ErrNoRows {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}
