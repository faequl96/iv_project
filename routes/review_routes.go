package routes

import (
	"iv_project/handlers"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func ReviewRoutes(r *mux.Router, jwtServices jwtToken.JWTServices) {
	reviewRepository := repositories.ReviewRepository(mysql.DB)
	userRepository := repositories.UserRepository(mysql.DB)
	h := handlers.ReviewHandlers(reviewRepository, userRepository)

	r.HandleFunc("/review", middleware.Auth(jwtServices, h.CreateReview)).Methods("POST")
	r.HandleFunc("/review/id/{id}", middleware.Auth(jwtServices, h.GetReviewByID)).Methods("GET")
	r.HandleFunc("/reviews", middleware.Auth(jwtServices, h.GetReviews)).Methods("GET")
	r.HandleFunc("/reviews/invitation-theme-id/{invitationThemeId}", middleware.Auth(jwtServices, h.GetReviewsByInvitationThemeID)).Methods("GET")
	r.HandleFunc("/review/id/{id}", middleware.Auth(jwtServices, h.UpdateReviewByID)).Methods("PATCH")
	r.HandleFunc("/review/id/{id}", middleware.Auth(jwtServices, h.DeleteReviewByID)).Methods("DELETE")
}
