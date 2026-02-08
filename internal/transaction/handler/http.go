package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/vadhe/api-category/internal/transaction/domain"
	"github.com/vadhe/api-category/internal/transaction/service"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: s}
}
func HandlerTransaction(w http.ResponseWriter, r *http.Request, h *TransactionHandler) {
	switch r.Method {
	case http.MethodGet:
		h.GetTransaction(w, r)
	// case http.MethodPost:
	// 	h.CreateProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func HandlerTransactionCheckout(w http.ResponseWriter, r *http.Request, h *TransactionHandler) {
	switch r.Method {
	case http.MethodPost:
		h.Checkout(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandlerReport(w http.ResponseWriter, r *http.Request, h *TransactionHandler) {
	switch r.Method {
	case http.MethodGet:
		h.GetReport(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	product, err := h.service.GetTransactions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TransactionHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetReport()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(report)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var res = domain.CheckoutRequest{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	checkout, err := h.service.Checkout(res.Items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(checkout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}

// func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
// 	var req domain.CheckoutRequest
// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	transaction, err := h.service.Checkout(req.Items)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(transaction)
// }
// func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
// 	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid ID format", http.StatusBadRequest)
// 		return
// 	}

// 	var res = domain.Product{}
// 	decoder := json.NewDecoder(r.Body)
// 	decoder.DisallowUnknownFields()

// 	if err = decoder.Decode(&res); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	product, err := h.service.UpdateProduct(id, res.Name, res.Price, res.Stock, res.CategoryId)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	var buf bytes.Buffer
// 	err = json.NewEncoder(&buf).Encode(product)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	w.Write(buf.Bytes())

// }

// func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
// 	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid ID format", http.StatusBadRequest)
// 		return
// 	}
// 	err = h.service.DeleteProduct(id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusNoContent)
// }
