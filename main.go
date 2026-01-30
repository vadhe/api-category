package main

import (
	"fmt"
	"net/http"

	"github.com/vadhe/api-category/internal/category/handler"
	"github.com/vadhe/api-category/internal/category/repository"
	"github.com/vadhe/api-category/internal/category/service"
	"github.com/vadhe/api-category/internal/database"
)

func main() {
	db, err := database.OpenPostgres()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()
	repo := repository.NewCategoryRepository(db)
	svc := service.NewCategoryService(repo)
	h := handler.NewCategoryHandler(svc)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})
	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		handler.HandlerCategory(w, r, h)
	})
	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		handler.HandlerCategoryById(w, r, h)
	})
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
