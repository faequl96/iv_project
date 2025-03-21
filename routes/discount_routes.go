package routes

import (
	"iv_project/handlers"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func DiscountRoutes(r *mux.Router, jwtServices jwtToken.JWTServices) {
	invitationThemeRepository := repositories.InvitationThemeRepository(mysql.DB)
	ivCoinPackageRepository := repositories.IVCoinPackageRepository(mysql.DB)
	h := handlers.DiscountHandlers(invitationThemeRepository, ivCoinPackageRepository)

	r.HandleFunc("/ad-mob", middleware.Auth(jwtServices, h.SetProductPrices)).Methods("PATCH")
}
