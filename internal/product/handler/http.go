package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	domain "github.com/vadhe/api-category/internal/product/domain"
	service "github.com/vadhe/api-category/internal/product/service"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}
func HandlerProduct(w http.ResponseWriter, r *http.Request, h *ProductHandler) {
	switch r.Method {
	case http.MethodGet:
		h.GetProduct(w, r)
	case http.MethodPost:
		h.CreateProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func HandlerProductById(w http.ResponseWriter, r *http.Request, h *ProductHandler) {
	switch r.Method {
	case http.MethodGet:
		h.GetProductByID(w, r)
	case http.MethodPut:
		h.UpdateProduct(w, r)
	case http.MethodDelete:
		h.DeleteProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	product, err := h.service.GetProducts()
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

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	product, err := h.service.GetProductByID(id)
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

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var res = domain.Product{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	product, err := h.service.CreateProduct(res.Name, res.Price, res.Stock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var res = domain.Product{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err = decoder.Decode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	product, err := h.service.UpdateProduct(id, res.Name, res.Price, res.Stock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())

}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteProduct(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
