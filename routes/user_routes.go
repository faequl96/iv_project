package routes

import (
	"iv_project/handlers"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	userRepository := repositories.UserRepository(mysql.DB)
	h := handlers.UserHandlers(userRepository)

	r.HandleFunc("/user", h.CreateUser).Methods("POST")
	r.HandleFunc("/user/id/{id}", h.GetUserByID).Methods("GET")
	r.HandleFunc("/users", h.GetUsers).Methods("GET")
	r.HandleFunc("/user/id/{id}", h.DeleteUser).Methods("DELETE")
}
