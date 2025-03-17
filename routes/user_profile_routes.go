package routes

import (
	"iv_project/handlers"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func UserProfileRoutes(r *mux.Router) {
	userProfileRepository := repositories.UserProfileRepository(mysql.DB)
	h := handlers.UserProfileHandlers(userProfileRepository)

	r.HandleFunc("/user-profile/{id}", h.GetUserProfileByID).Methods("GET")
	r.HandleFunc("/user-profile/{id}", h.UpdateUserProfile).Methods("PATCH")
}
