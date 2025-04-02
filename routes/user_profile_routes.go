package routes

import (
	"iv_project/handlers"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func UserProfileRoutes(r *mux.Router, jwtServices jwtToken.JWTServices) {
	userProfileRepository := repositories.UserProfileRepository(mysql.DB)
	h := handlers.UserProfileHandlers(userProfileRepository)

	r.Use(middleware.Language)

	r.HandleFunc("/user-profile", middleware.Auth(jwtServices, h.GetUserProfile)).Methods("GET")
	r.HandleFunc("/user-profile/id/{id}", middleware.Auth(jwtServices, h.GetUserProfileByID)).Methods("GET")
	r.HandleFunc("/user-profile", middleware.Auth(jwtServices, h.UpdateUserProfile)).Methods("PATCH")
	r.HandleFunc("/user-profile/id/{id}", middleware.Auth(jwtServices, h.UpdateUserProfileByID)).Methods("PATCH")
}
