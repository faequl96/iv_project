package handlers

import (
	"encoding/json"
	iv_coin_dto "iv_project/dto/iv_coin"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type ivCoinHandlers struct {
	IVCoinRepositories repositories.IVCoinRepositories
}

func IVCoinHandlers(IVCoinRepositories repositories.IVCoinRepositories) *ivCoinHandlers {
	return &ivCoinHandlers{IVCoinRepositories}
}

func ConvertToIVCoinResponse(ivCoin *models.IVCoin) iv_coin_dto.IVCoinResponse {
	return iv_coin_dto.IVCoinResponse{
		ID:      ivCoin.ID,
		Balance: ivCoin.Balance,
	}
}

func (h *ivCoinHandlers) GetIVCoin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value(middleware.UserIdKey).(string)
	iVCoin, err := h.IVCoinRepositories.GetIVCoinByUserID(userID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No iv coin found with the provided user.")
		return
	}

	SuccessResponse(w, http.StatusOK, "IV coin retrieved successfully", ConvertToIVCoinResponse(iVCoin))
}

func (h *ivCoinHandlers) GetIVCoinByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() && role != models.UserRoleAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	ivCoin, err := h.IVCoinRepositories.GetIVCoinByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "IV coin not found")
		return
	}

	SuccessResponse(w, http.StatusOK, "IV coin retrieved successfully", ConvertToIVCoinResponse(ivCoin))
}

func (h *ivCoinHandlers) UpdateIVCoinByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() && role != models.UserRoleAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	ivCoin, err := h.IVCoinRepositories.GetIVCoinByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "IV coin not found")
		return
	}

	var request iv_coin_dto.IVCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	ivCoin.Balance = request.Balance

	if err = h.IVCoinRepositories.UpdateIVCoin(ivCoin); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to update iv coin: "+err.Error())
		return
	}

	SuccessResponse(w, http.StatusOK, "IV coin updated successfully", ConvertToIVCoinResponse(ivCoin))
}
