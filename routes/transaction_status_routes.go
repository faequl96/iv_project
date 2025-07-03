package routes

import (
	"iv_project/handlers"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func TransactionStatusRoutes(r *mux.Router, jwtServices jwtToken.JWTServices) {
	transactionRepository := repositories.TransactionRepository(mysql.DB)
	invitationRepository := repositories.InvitationRepository(mysql.DB)
	ivCoinPackageRepository := repositories.IVCoinPackageRepository(mysql.DB)
	ivCoinRepository := repositories.IVCoinRepository(mysql.DB)
	voucherCodeRepository := repositories.VoucherCodeRepository(mysql.DB)
	userVoucherCodeUsageRepository := repositories.UserVoucherCodeUsageRepository(mysql.DB)
	h := handlers.TransactionStatusHandler(
		transactionRepository,
		invitationRepository,
		ivCoinPackageRepository,
		ivCoinRepository,
		voucherCodeRepository,
		userVoucherCodeUsageRepository,
	)

	r.Use(middleware.Language)

	r.HandleFunc("/transaction-status-check/reference-number/{referenceNumber}", middleware.Auth(jwtServices, h.CheckByReferenceNumber)).Methods("PATCH")
	r.HandleFunc("/transaction-status-reset/id/{id}", middleware.Auth(jwtServices, h.ResetByID)).Methods("PATCH")
}
