package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vadhe/api-category/internal/category/domain"
	category "github.com/vadhe/api-category/internal/category/handler"
)

var data []domain.Category = []domain.Category{
	{ID: 1, Name: "Electronics", Description: "Electronic devices"},
	{ID: 2, Name: "Clothing", Description: "Clothing items"},
	{ID: 3, Name: "Books", Description: "Books"},
}

func main() {
	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			path := strings.TrimPrefix(r.URL.Path, "/categories/")
			if path == "" || path == "/" {
				category.GetCategories(w, r, data)
			} else {
				category.GetCategoryByID(w, r, data)
			}
		case http.MethodPost:
			category.CreateCategory(w, r)
		case http.MethodPut:
			category.UpdateCategory(w, r, data)
		case http.MethodDelete:
			category.DeleteCategory(w, r, &data)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
