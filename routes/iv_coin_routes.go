package routes

import (
	"iv_project/handlers"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func IVCoinRoutes(r *mux.Router) {
	ivCoinRepository := repositories.IVCoinRepository(mysql.DB)
	h := handlers.IVCoinHandlers(ivCoinRepository)

	r.HandleFunc("/iv-coin/{id}", h.GetIVCoinByID).Methods("GET")
	r.HandleFunc("/iv-coin/{userId}", h.GetIVCoinByUserID).Methods("GET")
	r.HandleFunc("/iv-coin/{id}", h.UpdateIVCoin).Methods("PATCH")
}
