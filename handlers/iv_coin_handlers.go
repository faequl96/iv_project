package handlers

import (
	"encoding/json"
	"iv_project/dto"
	user_dto "iv_project/dto/user"
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

func (h *ivCoinHandlers) GetIVCoinByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	ivCoin, err := h.IVCoinRepositories.GetIVCoinByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{Code: http.StatusNotFound, Message: "IVCoin not found"}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: ivCoin}
	json.NewEncoder(w).Encode(response)
}

func (h *ivCoinHandlers) UpdateIVCoin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(user_dto.IVCoinRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	ivCoin, err := h.IVCoinRepositories.GetIVCoinByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{Code: http.StatusNotFound, Message: "IVCoin not found"}
		json.NewEncoder(w).Encode(response)
		return
	}

	ivCoin.Balance = request.Balance

	err = h.IVCoinRepositories.UpdateIVCoin(ivCoin)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: "IVCoin updated successfully"}
	json.NewEncoder(w).Encode(response)
}
