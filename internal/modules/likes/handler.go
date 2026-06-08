package likes

import (
	"Project2-v7/internal/shared/middleware"
	"Project2-v7/pkg/utils"
	"net/http"
)

type LikeHandler struct {
	service *LikeService
}

func NewLikeHandler(service *LikeService) *LikeHandler {
	return &LikeHandler{service: service}
}

func (h *LikeHandler) Like(w http.ResponseWriter, r *http.Request) {
	var req LikeRequest
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	claims, ok := middleware.GetClaimsFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	req.UserId = claims.UserId

	if err := h.service.Like(r.Context(), LikeRequest{BlogId: id, UserId: req.UserId}); err != nil {
		utils.WriteError(w, http.StatusConflict, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Blog liked successfully"})
}

func (h *LikeHandler) Unlike(w http.ResponseWriter, r *http.Request) {
	var req LikeRequest
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	claims, ok := middleware.GetClaimsFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	req.UserId = claims.UserId

	if err := h.service.Unlike(r.Context(), LikeRequest{BlogId: id, UserId: req.UserId}); err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *LikeHandler) GetCount(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	resp, err := h.service.GetCount(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}
