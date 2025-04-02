package routes

import (
	"iv_project/handlers"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func AuthRoutes(r *mux.Router, JwtServices jwtToken.JWTServices) {
	userRepository := repositories.UserRepository(mysql.DB)
	h := handlers.AuthHandlers(JwtServices, userRepository)

	r.Use(middleware.Language)

	r.HandleFunc("/login", h.Login).Methods("POST")
}
