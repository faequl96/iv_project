package routes

import (
	"iv_project/handlers"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func TransactionConfirmationRoutes(r *mux.Router, jwtServices jwtToken.JWTServices) {
	transactionRepository := repositories.TransactionRepository(mysql.DB)
	invitationRepository := repositories.InvitationRepository(mysql.DB)
	ivCoinPackageRepository := repositories.IVCoinPackageRepository(mysql.DB)
	ivCoinRepository := repositories.IVCoinRepository(mysql.DB)
	h := handlers.TransactionConfirmationHandler(
		transactionRepository,
		invitationRepository,
		ivCoinPackageRepository,
		ivCoinRepository,
	)

	r.HandleFunc("/transaction-confirmation-auto", middleware.Auth(jwtServices, h.AutoByMidtrans)).Methods("PATCH")
	r.HandleFunc("/transaction-confirmation-manual/id/{id}", middleware.Auth(jwtServices, h.ManualByAdminByID)).Methods("PATCH")
}
