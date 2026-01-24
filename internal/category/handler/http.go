package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	domain "github.com/vadhe/api-category/internal/category/domain"
	service "github.com/vadhe/api-category/internal/category/service"
)

func GetCategories(w http.ResponseWriter, r *http.Request, data []domain.Category) {

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetCategoryByID(w http.ResponseWriter, r *http.Request, data []domain.Category) {

	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	var res = domain.Category{}
	found := false
	for _, v := range data {
		if v.ID == id {
			res = v
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
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

func UpdateCategory(w http.ResponseWriter, r *http.Request, data []domain.Category) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	found := false
	foundIndex := -1
	for i, v := range data {
		if v.ID == id {
			found = true
			foundIndex = i
			break
		}
	}

	if !found {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}
	var res = domain.Category{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = service.CreateCategory(res.Name, res.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res.ID = id
	data[foundIndex].Name = res.Name
	data[foundIndex].Description = res.Description
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

func DeleteCategory(w http.ResponseWriter, r *http.Request, data *[]domain.Category) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	found := false
	for _, v := range *data {
		if v.ID == id {
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}
	var res = []domain.Category{}
	for _, val := range *data {
		if val.ID != id {
			res = append(res, val)
		}
	}
	*data = res
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}
