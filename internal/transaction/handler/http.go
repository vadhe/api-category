package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

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
	startDateParams := r.URL.Query().Get("start_date")
	endDateParams := r.URL.Query().Get("end_date")
	layout := "2006-01-02"
	var startDate, endDate time.Time
	if startDateParams != "" {
		v, err := time.Parse(layout, startDateParams)
		startDate = v
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if endDateParams != "" {
		v, err := time.Parse(layout, endDateParams)
		endDate = v
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	report, err := h.service.GetReport(startDate, endDate)
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
