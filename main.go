package main

import (
	"fmt"
	"net/http"

	handlerCategory "github.com/vadhe/api-category/internal/category/handler"
	handlerProduct "github.com/vadhe/api-category/internal/product/handler"
	handlerTransaction "github.com/vadhe/api-category/internal/transaction/handler"
	repositoryTransaction "github.com/vadhe/api-category/internal/transaction/repository"
	serviceTransaction "github.com/vadhe/api-category/internal/transaction/service"

	repositoryCategory "github.com/vadhe/api-category/internal/category/repository"
	serviceCategory "github.com/vadhe/api-category/internal/category/service"
	"github.com/vadhe/api-category/internal/database"
	repositoryProduct "github.com/vadhe/api-category/internal/product/repository"
	serviceProduct "github.com/vadhe/api-category/internal/product/service"
)

func main() {
	db, err := database.OpenPostgres()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()
	repoCategory := repositoryCategory.NewCategoryRepository(db)
	svcCategory := serviceCategory.NewCategoryService(repoCategory)
	hCategory := handlerCategory.NewCategoryHandler(svcCategory)

	repoProduct := repositoryProduct.NewProductRepository(db)
	svcProduct := serviceProduct.NewProductService(repoProduct)
	hProduct := handlerProduct.NewProductHandler(svcProduct)
	repoTransaction := repositoryTransaction.NewTransactionRepository(db)
	svcTransaction := serviceTransaction.NewTransactionService(repoTransaction, repoProduct, db)
	hTransaction := handlerTransaction.NewTransactionHandler(svcTransaction)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		handlerProduct.HandlerProduct(w, r, hProduct)
	})
	http.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		handlerProduct.HandlerProductById(w, r, hProduct)
	})
	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		handlerCategory.HandlerCategory(w, r, hCategory)
	})
	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		handlerTransaction.HandlerTransaction(w, r, hTransaction)
	})
	http.HandleFunc("/transactions-checkout", func(w http.ResponseWriter, r *http.Request) {
		handlerTransaction.HandlerTransactionCheckout(w, r, hTransaction)
	})
	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		handlerCategory.HandlerCategoryById(w, r, hCategory)
	})
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
