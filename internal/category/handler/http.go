package handler

import (
	"encoding/json"
	"net/http"

	domain "github.com/vadhe/api-category/internal/category/domain"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	res := domain.Category{
		ID:          1,
		Name:        "vadhe",
		Description: "description",
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	// Implementation of CreateCategory function
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	// Implementation of UpdateCategory function
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// Implementation of DeleteCategory function
}
