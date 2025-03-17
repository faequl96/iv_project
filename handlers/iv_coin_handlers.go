package handlers

import (
	"encoding/json"
	iv_coin_dto "iv_project/dto/iv_coin"
	"iv_project/models"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ivCoinHandlers struct {
	IVCoinRepositories repositories.IVCoinRepositories
}

func IVCoinHandlers(IVCoinRepositories repositories.IVCoinRepositories) *ivCoinHandlers {
	return &ivCoinHandlers{IVCoinRepositories}
}

func convertToIVCoinResponse(ivCoin *models.IVCoin) iv_coin_dto.IVCoinResponse {
	return iv_coin_dto.IVCoinResponse{
		ID:      ivCoin.ID,
		Balance: ivCoin.Balance,
	}
}

func (h *ivCoinHandlers) GetIVCoinByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	ivCoin, err := h.IVCoinRepositories.GetIVCoinByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "IVCoin not found")
		return
	}

	SuccessResponse(w, http.StatusOK, "IVCoin retrieved successfully", convertToIVCoinResponse(ivCoin))
}

func (h *ivCoinHandlers) UpdateIVCoin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	ivCoin, err := h.IVCoinRepositories.GetIVCoinByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "IVCoin not found")
		return
	}

	request := new(iv_coin_dto.IVCoinRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	ivCoin.Balance = request.Balance

	if err = h.IVCoinRepositories.UpdateIVCoin(ivCoin); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to update IVCoin: "+err.Error())
		return
	}

	SuccessResponse(w, http.StatusOK, "IVCoin updated successfully", convertToIVCoinResponse(ivCoin))
}
