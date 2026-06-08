package tag

import (
	"Project2-v7/pkg/utils"
	"encoding/json"
	"net/http"
)

type TagHandler struct {
	service *TagService
}

func NewTagHandler(service *TagService) *TagHandler {
	return &TagHandler{service: service}
}

func (th *TagHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	var req TagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if _, err := th.service.Create(req); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Tag created successfully"})
}

func (th *TagHandler) GetAllTags(w http.ResponseWriter, r *http.Request) {
	tags, err := th.service.GetAll()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, tags)
}

func (th *TagHandler) GetTagById(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}
	tag, err := th.service.GetById(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "tag not found")
		return
	}
	utils.WriteJSON(w, http.StatusOK, tag)
}

func (th *TagHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}
	var req TagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := th.service.Update(id, req); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Tag updated successfully"})
}

func (th *TagHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}
	if err := th.service.Delete(id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
