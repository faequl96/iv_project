package routes

import (
	"iv_project/handlers"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func InvitationDataRoutes(r *mux.Router, jwtServices jwtToken.JWTServices) {
	invitationDataRepository := repositories.InvitationDataRepository(mysql.DB)
	h := handlers.InvitationDataHandler(invitationDataRepository)

	r.HandleFunc("/invitation-data/id/{id}", middleware.Auth(jwtServices, h.GetInvitationDataByID)).Methods("GET")
	r.HandleFunc("/invitation-data/id/{id}", middleware.Auth(jwtServices, middleware.InvitationImagesUploader(h.UpdateInvitationDataByID))).Methods("PATCH")
}
