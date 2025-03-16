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

	r.HandleFunc("/user/{id}", h.GetIVCoinByID).Methods("GET")
	r.HandleFunc("/user/{id}", h.UpdateIVCoin).Methods("PATCH")
}
