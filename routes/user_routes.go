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
	r.HandleFunc("/user/id/{id}", middleware.Auth(jwtServices, h.GetUserByID)).Methods("GET")
	r.HandleFunc("/users", middleware.Auth(jwtServices, h.GetUsers)).Methods("GET")
	r.HandleFunc("/user/id/{id}", middleware.Auth(jwtServices, h.UpdateUserByID)).Methods("PATCH")
	r.HandleFunc("/user", middleware.Auth(jwtServices, h.DeleteUser)).Methods("DELETE")
	r.HandleFunc("/user/id/{id}", middleware.Auth(jwtServices, h.DeleteUserByID)).Methods("DELETE")
}
