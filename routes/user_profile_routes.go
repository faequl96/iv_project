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

	r.HandleFunc("/user-profile/id/{id}", h.GetUserProfileByID).Methods("GET")
	r.HandleFunc("/user-profile/user-id/{userId}", h.GetUserProfileByUserID).Methods("GET")
	r.HandleFunc("/user-profile/id/{id}", h.UpdateUserProfile).Methods("PATCH")
}
