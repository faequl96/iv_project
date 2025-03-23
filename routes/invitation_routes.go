package routes

import (
	"iv_project/handlers"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func InvitationRoutes(r *mux.Router, jwtServices jwtToken.JWTServices) {
	invitationRepository := repositories.InvitationRepository(mysql.DB)
	invitationThemeRepository := repositories.InvitationThemeRepository(mysql.DB)
	h := handlers.InvitationHandler(invitationRepository, invitationThemeRepository)

	r.HandleFunc("/invitation", middleware.Auth(jwtServices, middleware.InvitationImagesUploader(h.CreateInvitation))).Methods("POST")
	r.HandleFunc("/invitation/id/{id}", middleware.Auth(jwtServices, h.GetInvitationByID)).Methods("GET")
	r.HandleFunc("/invitations", middleware.Auth(jwtServices, h.GetInvitations)).Methods("GET")
	r.HandleFunc("/invitations/user-id/{userId}", middleware.Auth(jwtServices, h.GetInvitationsByUserID)).Methods("GET")
	r.HandleFunc("/invitation/id/{id}", middleware.Auth(jwtServices, middleware.InvitationImagesUploader(h.UpdateInvitationByID))).Methods("PATCH")
	r.HandleFunc("/invitation/id/{id}", middleware.Auth(jwtServices, h.DeleteInvitationByID)).Methods("DELETE")
}
