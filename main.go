package main

import (
	"fmt"
	"net/http"
	"strings"

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
	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			path := strings.TrimPrefix(r.URL.Path, "/categories/")
			if path == "" || path == "/" {
				h.GetCategories(w, r)
			} else {
				h.GetCategoryByID(w, r)
			}
		case http.MethodPost:
			h.CreateCategory(w, r)
		case http.MethodPut:
			h.UpdateCategory(w, r)
		case http.MethodDelete:
			h.DeleteCategory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
