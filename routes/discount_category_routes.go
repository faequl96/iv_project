package routes

import (
	"iv_project/handlers"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func DiscountCategoryRoutes(r *mux.Router, jwtServices jwtToken.JWTServices) {
	discountCategoryRepository := repositories.DiscountCategoryRepository(mysql.DB)
	h := handlers.DiscountCategoryHandler(discountCategoryRepository)

	r.HandleFunc("/discount-category", middleware.Auth(jwtServices, h.CreateDiscountCategory)).Methods("POST")
	r.HandleFunc("/discount-category/id/{id}", middleware.Auth(jwtServices, h.GetDiscountCategoryByID)).Methods("GET")
	r.HandleFunc("/discount-categories", middleware.Auth(jwtServices, h.GetDiscountCategories)).Methods("GET")
	r.HandleFunc("/discount-category/id/{id}", middleware.Auth(jwtServices, h.UpdateDiscountCategory)).Methods("PATCH")
	r.HandleFunc("/discount-category/id/{id}", middleware.Auth(jwtServices, h.DeleteDiscountCategory)).Methods("DELETE")
}
