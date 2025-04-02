package routes

import (
	"iv_project/handlers"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func AdMobRoutes(r *mux.Router, jwtServices jwtToken.JWTServices) {
	ivCoinRepository := repositories.IVCoinRepository(mysql.DB)
	h := handlers.AdMobHandlers(ivCoinRepository)

	r.Use(middleware.Language)

	r.HandleFunc("/ad-mob", middleware.Auth(jwtServices, h.AddExtraIVCoins)).Methods("PATCH")
}
