package routes

import (
	"iv_project/handlers"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func IVCoinPackageRoutes(r *mux.Router) {
	ivCoinPackageRepository := repositories.IVCoinPackageRepository(mysql.DB)
	h := handlers.IVCoinPackageHandler(ivCoinPackageRepository)

	r.HandleFunc("/iv-coin-package", h.CreateIVCoinPackage).Methods("POST")
	r.HandleFunc("/iv-coin-package/id/{id}", h.GetIVCoinPackageByID).Methods("GET")
	r.HandleFunc("/iv-coin-packages", h.GetIVCoinPackages).Methods("GET")
	r.HandleFunc("/iv-coin-package/id/{id}", h.UpdateIVCoinPackage).Methods("PATCH")
	r.HandleFunc("/iv-coin-package/id/{id}", h.DeleteIVCoinPackage).Methods("DELETE")
}
