package handlers

import (
	"encoding/json"
	"iv_project/dto"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type galleryHandlers struct {
	GalleryRepositories repositories.GalleryRepositories
}

func GalleryHandler(GalleryRepositories repositories.GalleryRepositories) *galleryHandlers {
	return &galleryHandlers{GalleryRepositories}
}

func (h *galleryHandlers) GetGalleryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	gallery, err := h.GalleryRepositories.GetGalleryByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: gallery}
	json.NewEncoder(w).Encode(response)
}

func (h *galleryHandlers) DeleteGallery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	gallery, err := h.GalleryRepositories.GetGalleryByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	err = h.GalleryRepositories.DeleteGallery(gallery)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: "Gallery deleted successfully"}
	json.NewEncoder(w).Encode(response)
}
