package routes

import (
	"iv_project/handlers"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func IVCoinPackageRoutes(r *mux.Router, jwtServices jwtToken.JWTServices) {
	ivCoinPackageRepository := repositories.IVCoinPackageRepository(mysql.DB)
	discountCategoryRepository := repositories.DiscountCategoryRepository(mysql.DB)
	h := handlers.IVCoinPackageHandler(ivCoinPackageRepository, discountCategoryRepository)

	r.HandleFunc("/iv-coin-package", middleware.Auth(jwtServices, h.CreateIVCoinPackage)).Methods("POST")
	r.HandleFunc("/iv-coin-package/id/{id}", middleware.Auth(jwtServices, h.GetIVCoinPackageByID)).Methods("GET")
	r.HandleFunc("/iv-coin-packages", middleware.Auth(jwtServices, h.GetIVCoinPackages)).Methods("GET")
	r.HandleFunc("/iv-coin-package/id/{id}", middleware.Auth(jwtServices, h.UpdateIVCoinPackageByID)).Methods("PATCH")
	r.HandleFunc("/iv-coin-package/id/{id}", middleware.Auth(jwtServices, h.DeleteIVCoinPackageByID)).Methods("DELETE")
}
