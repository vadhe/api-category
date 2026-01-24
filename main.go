package main

import (
	"fmt"
	"net/http"

	category "github.com/vadhe/api-category/internal/category/handler"
)

func main() {
	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			category.GetCategories(w, r)
		case http.MethodPost:
			category.CreateCategory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	// http.HandleFunc("/categories", category.UpdateCategory)
	// http.HandleFunc("/categories", category.DeleteCategory)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
