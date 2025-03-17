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

// IVCoinHandlers initializes the handler with the given repository.
func IVCoinHandlers(IVCoinRepositories repositories.IVCoinRepositories) *ivCoinHandlers {
	return &ivCoinHandlers{IVCoinRepositories}
}

// convertToIVCoinResponse converts the IVCoin model to IVCoinResponse DTO
func convertToIVCoinResponse(ivCoin *models.IVCoin) iv_coin_dto.IVCoinResponse {
	return iv_coin_dto.IVCoinResponse{
		ID:      ivCoin.ID,
		Balance: ivCoin.Balance,
	}
}

// GetIVCoinByID retrieves the IVCoin balance by ID
func (h *ivCoinHandlers) GetIVCoinByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the ID from the request URL parameters
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// Fetch IVCoin data from the database
	ivCoin, err := h.IVCoinRepositories.GetIVCoinByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "IVCoin not found")
		return
	}

	// Send a successful response with the IVCoin data
	SuccessResponse(w, http.StatusOK, "IVCoin retrieved successfully", convertToIVCoinResponse(ivCoin))
}

// UpdateIVCoin updates the IVCoin balance by ID
func (h *ivCoinHandlers) UpdateIVCoin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode the request body
	request := new(iv_coin_dto.IVCoinRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	// Extract the ID from the request URL parameters
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// Fetch IVCoin data from the database
	ivCoin, err := h.IVCoinRepositories.GetIVCoinByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "IVCoin not found")
		return
	}

	// Update the IVCoin balance
	ivCoin.Balance = request.Balance

	// Save the updated IVCoin data to the database
	if err = h.IVCoinRepositories.UpdateIVCoin(ivCoin); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to update IVCoin: "+err.Error())
		return
	}

	// Send a successful response with the updated IVCoin data
	SuccessResponse(w, http.StatusOK, "IVCoin updated successfully", convertToIVCoinResponse(ivCoin))
}
