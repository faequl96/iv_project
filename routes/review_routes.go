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

	r.HandleFunc("/review", h.CreateReview).Methods("POST")
	r.HandleFunc("/review/id/{id}", h.GetReviewByID).Methods("GET")
	r.HandleFunc("/reviews/invitation-theme-id/{invitationThemeId}", h.GetReviewsByInvitationThemeID).Methods("GET")
	r.HandleFunc("/review/id/{id}", h.UpdateReview).Methods("PATCH")
	r.HandleFunc("/review/id/{id}", h.DeleteReview).Methods("DELETE")
}
