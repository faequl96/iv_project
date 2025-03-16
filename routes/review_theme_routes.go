package routes

import (
	"iv_project/handlers"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func ReviewRoutes(r *mux.Router) {
	reviewRepository := repositories.ReviewRepository(mysql.DB)
	h := handlers.ReviewHandlers(reviewRepository)

	r.HandleFunc("/reviews", h.CreateReview).Methods("POST")
	r.HandleFunc("/reviews/{id}", h.GetReviewByID).Methods("GET")
	r.HandleFunc("/reviews", h.GetReviews).Methods("GET")
	r.HandleFunc("/reviews/{invitationThemeId}", h.GetReviewsByInvitationThemeID).Methods("GET")
	r.HandleFunc("/reviews/{id}", h.UpdateReview).Methods("PATCH")
	r.HandleFunc("/reviews/{id}", h.DeleteReview).Methods("DELETE")
}
