package routes

import (
	"iv_project/handlers"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func TransactionPaymentRoutes(r *mux.Router, jwtServices jwtToken.JWTServices) {
	transactionRepository := repositories.TransactionRepository(mysql.DB)
	invitationRepository := repositories.InvitationRepository(mysql.DB)
	ivCoinPackageRepository := repositories.IVCoinPackageRepository(mysql.DB)
	ivCoinRepository := repositories.IVCoinRepository(mysql.DB)
	userRepository := repositories.UserRepository(mysql.DB)
	h := handlers.TransactionPaymentHandler(
		transactionRepository,
		invitationRepository,
		ivCoinPackageRepository,
		ivCoinRepository,
		userRepository,
	)

	r.HandleFunc("/issue/id/{id}", middleware.Auth(jwtServices, h.IssueByID)).Methods("PATCH")
}
