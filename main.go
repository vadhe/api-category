package main

import (
	"fmt"
	"net/http"

	category "github.com/vadhe/api-category/internal/category/handler"
)

func main() {
	fmt.Println("Hello, World!")
	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		category.GetCategories(w, r)
	})
	// http.HandleFunc("/categories", category.CreateCategory)
	// http.HandleFunc("/categories", category.UpdateCategory)
	// http.HandleFunc("/categories", category.DeleteCategory)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
