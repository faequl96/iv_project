package handlers

import (
	gallery_dto "iv_project/dto/gallery"
	"iv_project/models"
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

func ConvertToGalleryResponse(gallery *models.Gallery) gallery_dto.GalleryResponse {
	return gallery_dto.GalleryResponse{
		ID:         gallery.ID,
		ImageURL1:  gallery.ImageURL1,
		ImageURL2:  gallery.ImageURL2,
		ImageURL3:  gallery.ImageURL3,
		ImageURL4:  gallery.ImageURL4,
		ImageURL5:  gallery.ImageURL5,
		ImageURL6:  gallery.ImageURL6,
		ImageURL7:  gallery.ImageURL7,
		ImageURL8:  gallery.ImageURL8,
		ImageURL9:  gallery.ImageURL9,
		ImageURL10: gallery.ImageURL10,
		ImageURL11: gallery.ImageURL11,
		ImageURL12: gallery.ImageURL12,
	}
}

func (h *galleryHandlers) GetGalleryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid gallery ID format. Please provide a numeric ID.")
		return
	}

	gallery, err := h.GalleryRepositories.GetGalleryByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No gallery found with the provided ID.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Gallery retrieved successfully", ConvertToGalleryResponse(gallery))
}

func (h *galleryHandlers) DeleteGallery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid gallery ID format. Please provide a numeric ID.")
		return
	}

	if _, err = h.GalleryRepositories.GetGalleryByID(uint(id)); err != nil {
		ErrorResponse(w, http.StatusNotFound, "No gallery found with the provided ID.")
		return
	}

	if err := h.GalleryRepositories.DeleteGallery(uint(id)); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while deleting the gallery.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Gallery deleted successfully", nil)
}
