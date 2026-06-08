package user

import (
	"Project2-v7/pkg/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := uh.service.Create(req); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
}

func (uh *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uh.service.GetAll()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, users)
}

func (uh *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	user, err := uh.service.GetById(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := uh.service.Update(id, req); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}

func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	if err := uh.service.Delete(id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req ChangeRoleRequest
	if err := utils.DecodeAndValidate(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.ChangeRole(r.Context(), id, req); err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.WriteJSON(w, http.StatusNotFound, "user not found")
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "role updated")
}
