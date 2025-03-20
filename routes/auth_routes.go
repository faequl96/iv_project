package routes

import (
	"iv_project/handlers"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func AuthRoutes(r *mux.Router) {
	userRepository := repositories.UserRepository(mysql.DB)
	h := handlers.AuthHandlers(userRepository)

	r.HandleFunc("/login", h.Login).Methods("POST")
}
