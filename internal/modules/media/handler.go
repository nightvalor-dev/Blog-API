package media

import (
	"Project2-v7/pkg/utils"
	"net/http"
)

type MediaHandler struct {
	service *MediaService
}

func NewMediaHandler(service *MediaService) *MediaHandler {
	return &MediaHandler{service: service}
}

func (h *MediaHandler) Upload(w http.ResponseWriter, r *http.Request) {
	blogId, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	// 10MB max
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "image field missing")
		return
	}
	defer file.Close()

	resp, err := h.service.Upload(r.Context(), blogId, file)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *MediaHandler) GetByBlogId(w http.ResponseWriter, r *http.Request) {
	blogId, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	media, err := h.service.GetByBlogId(r.Context(), blogId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, media)
}

func (h *MediaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	mediaId, err := utils.ParseChiParam(w, r, "mid")
	if !err {
		return
	}

	if err := h.service.Delete(r.Context(), mediaId); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
