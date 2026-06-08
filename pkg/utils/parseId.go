package utils

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ParseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return id, true
}

func ParseChiParam(w http.ResponseWriter, r *http.Request, param string) (int, bool) {
	str := chi.URLParam(r, param)
	val, err := strconv.Atoi(str)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid "+param)
		return 0, false
	}
	return val, true
}
