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

	r.HandleFunc("/create-invitation", middleware.InvitationImagesUploader(h.CreateInvitation)).Methods("POST")
	r.HandleFunc("/invitation/{id}", h.GetInvitationByID).Methods("GET")
	r.HandleFunc("/invitations/{id}", h.GetInvitationsByUserID).Methods("GET")
}
