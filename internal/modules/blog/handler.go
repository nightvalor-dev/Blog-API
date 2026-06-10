package blog

import (
	"Project2-v7/internal/shared/middleware"
	"Project2-v7/pkg/utils"
	"encoding/json"
	"net/http"
)

type BlogHandler struct {
	service *BlogService
}

func NewBlogHandler(service *BlogService) *BlogHandler {
	return &BlogHandler{service: service}
}

func (bh *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	var req CreateBlogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	claims, ok := middleware.GetClaimsFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	req.UserId = claims.UserId
	err := bh.service.Create(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{"message": "Blog created successfully"})
}

func (bh *BlogHandler) GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	filter := BlogFilter{
		Status:       r.URL.Query().Get("status"),
		SortBy:       r.URL.Query().Get("sort_by"),
		Order:        r.URL.Query().Get("order"),
		CategoryName: r.URL.Query().Get("category"),
		TagName:      r.URL.Query().Get("tag"),
	}

	uid := utils.ParseIntQuery(r, "user_id")
	if uid > 0 {
		filter.UserId = uid
	}

	pagination := utils.ParsePagination(r)

	if isFilteredRequest(filter) {
		result, err := bh.service.GetAllFiltered(r.Context(), filter, pagination)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.WriteJSON(w, http.StatusOK, result)
		return
	}

	result, err := bh.service.GetAllCached(r.Context(), pagination)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, result)
}

// isFilteredRequest returns true if the request carries any filter or sort param.
func isFilteredRequest(f BlogFilter) bool {
	return f.Status != "" || f.UserId > 0 || f.SortBy != "" || f.Order != "" ||
		f.CategoryName != "" || f.TagName != ""
}

func (bh *BlogHandler) GetBlogById(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	blog, err := bh.service.GetById(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "blog not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, blog)
}

func (bh *BlogHandler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	var req UpdateBlogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := bh.service.Update(r.Context(), id, req); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Blog updated successfully"})
}

func (bh *BlogHandler) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	if err := bh.service.Delete(r.Context(), id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
