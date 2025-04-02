package routes

import (
	"iv_project/handlers"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func CategoryRoutes(r *mux.Router, jwtServices jwtToken.JWTServices) {
	categoryRepository := repositories.CategoryRepository(mysql.DB)
	h := handlers.CategoryHandler(categoryRepository)

	r.Use(middleware.Language)

	r.HandleFunc("/category", middleware.Auth(jwtServices, h.CreateCategory)).Methods("POST")
	r.HandleFunc("/category/id/{id}", middleware.Auth(jwtServices, h.GetCategoryByID)).Methods("GET")
	r.HandleFunc("/categories", middleware.Auth(jwtServices, h.GetCategories)).Methods("GET")
	r.HandleFunc("/category/id/{id}", middleware.Auth(jwtServices, h.UpdateCategoryByID)).Methods("PATCH")
	r.HandleFunc("/category/id/{id}", middleware.Auth(jwtServices, h.DeleteCategoryByID)).Methods("DELETE")
}
