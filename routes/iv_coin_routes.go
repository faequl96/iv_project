package routes

import (
	"iv_project/handlers"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func IVCoinRoutes(r *mux.Router, jwtServices jwtToken.JWTServices) {
	ivCoinRepository := repositories.IVCoinRepository(mysql.DB)
	h := handlers.IVCoinHandlers(ivCoinRepository)

	r.Use(middleware.Language)

	r.HandleFunc("/iv-coin", middleware.Auth(jwtServices, h.GetIVCoin)).Methods("GET")
	r.HandleFunc("/iv-coin/id/{id}", middleware.Auth(jwtServices, h.GetIVCoinByID)).Methods("GET")
	r.HandleFunc("/iv-coin/id/{id}", middleware.Auth(jwtServices, h.UpdateIVCoinByID)).Methods("PATCH")
}
