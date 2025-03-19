package routes

import (
	"iv_project/handlers"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func InvitationRoutes(r *mux.Router) {
	invitationRepository := repositories.InvitationRepository(mysql.DB)
	h := handlers.InvitationHandler(invitationRepository)

	r.HandleFunc("/invitation", middleware.InvitationImagesUploader(h.CreateInvitation)).Methods("POST")
	r.HandleFunc("/invitationid//{id}", h.GetInvitationByID).Methods("GET")
	r.HandleFunc("/invitations", h.GetInvitations).Methods("GET")
	r.HandleFunc("/invitations/user-id/{userId}", h.GetInvitationsByUserID).Methods("GET")
	r.HandleFunc("/invitationid//{id}", middleware.InvitationImagesUploader(h.UpdateInvitation)).Methods("PATCH")
	r.HandleFunc("/invitationid//{id}", h.CreateInvitation).Methods("DELETE")
}
