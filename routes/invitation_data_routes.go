package routes

import (
	"iv_project/handlers"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func InvitationDataRoutes(r *mux.Router) {
	invitationDataRepository := repositories.InvitationDataRepository(mysql.DB)
	h := handlers.InvitationDataHandler(invitationDataRepository)

	r.HandleFunc("/invitation-data", middleware.InvitationDataImagesUploader(h.CreateInvitationData)).Methods("POST")
	r.HandleFunc("/invitation-data/{id}", h.GetInvitationDataByID).Methods("GET")
	r.HandleFunc("/invitation-data/{id}", middleware.InvitationDataImagesUploader(h.UpdateInvitationData)).Methods("PATCH")
	r.HandleFunc("/invitation-data/{id}", h.DeleteInvitationData).Methods("DELETE")
}
