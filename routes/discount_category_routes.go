package routes

import (
	"iv_project/handlers"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func DiscountCategoryRoutes(r *mux.Router) {
	discountCategoryRepository := repositories.DiscountCategoryRepository(mysql.DB)
	h := handlers.DiscountCategoryHandler(discountCategoryRepository)

	r.HandleFunc("/category", h.CreateDiscountCategory).Methods("POST")
	r.HandleFunc("/category/id/{id}", h.GetDiscountCategoryByID).Methods("GET")
	r.HandleFunc("/categories", h.GetDiscountCategories).Methods("GET")
	r.HandleFunc("/category/id/{id}", h.UpdateDiscountCategory).Methods("PATCH")
	r.HandleFunc("/category/id/{id}", h.DeleteDiscountCategory).Methods("DELETE")
}
