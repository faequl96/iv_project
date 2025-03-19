package routes

import (
	"iv_project/handlers"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func TransactionRoutes(r *mux.Router) {
	transactionRepository := repositories.TransactionRepository(mysql.DB)
	invitationRepository := repositories.InvitationRepository(mysql.DB)
	ivCoinPackageRepository := repositories.IVCoinPackageRepository(mysql.DB)
	ivCoinRepository := repositories.IVCoinRepository(mysql.DB)
	h := handlers.TransactionHandler(transactionRepository, invitationRepository, ivCoinPackageRepository, ivCoinRepository)

	r.HandleFunc("/transaction", h.CreateTransaction).Methods("POST")
	r.HandleFunc("/transaction/id/{id}", h.GetTransactionByID).Methods("GET")
	r.HandleFunc("/transactions", h.GetTransactions).Methods("GET")
	r.HandleFunc("/transactions/user-id/{userId}", h.GetTransactionsByUserID).Methods("GET")
	r.HandleFunc("/transaction/id/{id}", h.UpdateTransaction).Methods("PATCH")
	r.HandleFunc("/transaction/id/{id}", h.CreateTransaction).Methods("DELETE")
}
