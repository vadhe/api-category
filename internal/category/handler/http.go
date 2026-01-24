package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	domain "github.com/vadhe/api-category/internal/category/domain"
	service "github.com/vadhe/api-category/internal/category/service"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {
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
	var res = domain.Category{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = service.CreateCategory(res.Name, res.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	// Implementation of UpdateCategory function
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// Implementation of DeleteCategory function
}
