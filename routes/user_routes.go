package routes

import (
	"iv_project/handlers"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router, jwtServices jwtToken.JWTServices) {
	userRepository := repositories.UserRepository(mysql.DB)
	h := handlers.UserHandlers(userRepository)

	r.HandleFunc("/user", middleware.Auth(jwtServices, h.GetUser)).Methods("GET")
	r.HandleFunc("/user/id/{id}", h.GetUserByID).Methods("GET")
	r.HandleFunc("/users", h.GetUsers).Methods("GET")
	r.HandleFunc("/user/id/{id}", h.UpdateUser).Methods("PATCH")
	r.HandleFunc("/user/id/{id}", h.DeleteUser).Methods("DELETE")
}
