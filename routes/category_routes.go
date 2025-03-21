package routes

import (
	"iv_project/handlers"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func CategoryRoutes(r *mux.Router) {
	iategoryRepository := repositories.CategoryRepository(mysql.DB)
	h := handlers.CategoryHandler(iategoryRepository)

	r.HandleFunc("/category", h.CreateCategory).Methods("POST")
	r.HandleFunc("/category/id/{id}", h.GetCategoryByID).Methods("GET")
	r.HandleFunc("/categories", h.GetCategories).Methods("GET")
	r.HandleFunc("/category/id/{id}", h.UpdateCategory).Methods("PATCH")
	r.HandleFunc("/category/id/{id}", h.DeleteCategory).Methods("DELETE")
}
