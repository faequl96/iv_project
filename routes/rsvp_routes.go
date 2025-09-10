package routes

import (
	"iv_project/handlers"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func RSVPRoutes(r *mux.Router) {
	rsvpRepository := repositories.RSVPRepository(mysql.DB)
	h := handlers.RSVPHandlers(rsvpRepository)

	r.Use(middleware.Language)

	r.HandleFunc("/rsvp", h.CreateRSVP).Methods("POST")
	r.HandleFunc("/rsvps/invitation-id/{invitationId}", h.GetRSVPsByInvitationID).Methods("GET")
}
