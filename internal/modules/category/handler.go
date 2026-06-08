package category

import (
	"Project2-v7/pkg/utils"
	"encoding/json"
	"net/http"
)

type CategoryHandler struct {
	service *CategoryService
}

func NewCategoryHandler(service *CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (ch *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if _, err := ch.service.Create(req); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Category created successfully"})
}

func (ch *CategoryHandler) GetAllCategory(w http.ResponseWriter, r *http.Request) {
	categories, err := ch.service.GetAll()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, categories)
}

func (ch *CategoryHandler) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}
	category, err := ch.service.GetById(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "category not found")
		return
	}
	utils.WriteJSON(w, http.StatusOK, category)
}

func (ch *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}
	var req CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := ch.service.Update(id, req); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Category updated successfully"})
}

func (ch *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}
	if err := ch.service.Delete(id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
