package routes

import (
	"iv_project/handlers"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func InvitedGuestRoutes(r *mux.Router) {
	invitedGuestRepository := repositories.InvitedGuestRepository(mysql.DB)
	h := handlers.InvitedGuestHandlers(invitedGuestRepository)

	r.Use(middleware.Language)

	r.HandleFunc("/invited-guest", h.CreateInvitedGuest).Methods("POST")
	r.HandleFunc("/invited-guests/invitation-id/{invitationId}", h.GetInvitedGuestsByInvitationID).Methods("GET")
	r.HandleFunc("/invited-guest/id/{id}", h.UpdateInvitedGuestByID).Methods("PATCH")
}
