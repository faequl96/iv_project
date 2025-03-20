package handlers

import (
	"encoding/json"
	ad_mob_dto "iv_project/dto/ad_mob"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/repositories"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type adMobHandlers struct {
	IVCoinRepositories repositories.IVCoinRepositories
}

func AdMobHandlers(IVCoinRepositories repositories.IVCoinRepositories) *adMobHandlers {
	return &adMobHandlers{IVCoinRepositories}
}

func ConvertToAdMobResponse(ivCoin *models.IVCoin) ad_mob_dto.AdMobResponse {
	return ad_mob_dto.AdMobResponse{
		ID:      ivCoin.ID,
		Balance: ivCoin.Balance,
	}
}

func (h *adMobHandlers) AddExtraIVCoins(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value(middleware.UserIdKey).(string)
	ivCoin, err := h.IVCoinRepositories.GetIVCoinByUserID(userID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No iv coin found with the provided user.")
		return
	}

	var request ad_mob_dto.AdMobRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	if request.Amount != 0 {
		ivCoin.Balance += request.Amount
	}

	if err = h.IVCoinRepositories.UpdateIVCoin(ivCoin); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to update iv coin: "+err.Error())
		return
	}

	SuccessResponse(w, http.StatusOK, "IV coin updated successfully", ConvertToIVCoinResponse(ivCoin))
}
