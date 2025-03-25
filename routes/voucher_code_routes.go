package routes

import (
	"iv_project/handlers"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func VoucherCodeRoutes(r *mux.Router, jwtServices jwtToken.JWTServices) {
	voucherCodeRepository := repositories.VoucherCodeRepository(mysql.DB)
	h := handlers.VoucherCodeHandler(voucherCodeRepository)

	r.HandleFunc("/voucher-code", middleware.Auth(jwtServices, h.CreateVoucherCode)).Methods("POST")
	r.HandleFunc("/voucher-code/id/{id}", middleware.Auth(jwtServices, h.GetVoucherCodeByID)).Methods("GET")
	r.HandleFunc("/voucher-code/id/{id}", middleware.Auth(jwtServices, h.UpdateVoucherCodeByID)).Methods("PATCH")
	r.HandleFunc("/voucher-code/id/{id}", middleware.Auth(jwtServices, h.DeleteVoucherCodeByID)).Methods("DELETE")
}
