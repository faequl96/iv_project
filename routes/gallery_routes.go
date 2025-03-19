package routes

import (
	"iv_project/handlers"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func GalleryRoutes(r *mux.Router) {
	galleryRepository := repositories.GalleryRepository(mysql.DB)
	h := handlers.GalleryHandler(galleryRepository)

	r.HandleFunc("/gallery/id/{id}", h.GetGalleryByID).Methods("GET")
	r.HandleFunc("/gallery/id/{id}", h.DeleteGallery).Methods("DELETE")
}
