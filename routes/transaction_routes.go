package routes

import (
	"iv_project/handlers"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func TransactionRoutes(r *mux.Router, jwtServices jwtToken.JWTServices) {
	transactionRepository := repositories.TransactionRepository(mysql.DB)
	invitationRepository := repositories.InvitationRepository(mysql.DB)
	ivCoinPackageRepository := repositories.IVCoinPackageRepository(mysql.DB)
	ivCoinRepository := repositories.IVCoinRepository(mysql.DB)
	userRepository := repositories.UserRepository(mysql.DB)
	h := handlers.TransactionHandler(
		transactionRepository,
		invitationRepository,
		ivCoinPackageRepository,
		ivCoinRepository,
		userRepository,
	)

	r.HandleFunc("/transaction", middleware.Auth(jwtServices, h.CreateTransaction)).Methods("POST")
	r.HandleFunc("/transaction/id/{id}", middleware.Auth(jwtServices, h.GetTransactionByID)).Methods("GET")
	r.HandleFunc("/transactions", middleware.Auth(jwtServices, h.GetTransactions)).Methods("GET")
	r.HandleFunc("/transactions/user-id/{userId}", middleware.Auth(jwtServices, h.GetTransactionsByUserID)).Methods("GET")
	r.HandleFunc("/transaction/id/{id}", middleware.Auth(jwtServices, h.UpdateTransactionPaymentMethodByID)).Methods("PATCH")
	r.HandleFunc("/transaction/id/{id}", middleware.Auth(jwtServices, h.CreateTransaction)).Methods("DELETE")
}
