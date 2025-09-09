package handlers

import (
	"encoding/json"
	ad_mob_dto "iv_project/dto/ad_mob"
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

func (h *adMobHandlers) AddExtraIVCoins(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIdKey).(string)
	ivCoin, err := h.IVCoinRepositories.GetIVCoinByUserID(userID)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No iv coin found with the provided user.",
			"id": "IV coin tidak ditemukan berdasarkan pengguna ini.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	var request ad_mob_dto.AdMobRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid request format",
			"id": "Format request tidak valid",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	if err := validator.New().Struct(request); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Validation failed. Please complete the request field",
			"id": "Validasi gagal. Silahkan lengkapi field request",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	if request.Amount != 0 {
		if ivCoin.AdMobMarker < 5 {
			ivCoin.AdMobMarker += 1
			ivCoin.Balance += request.Amount
		} else {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "Free iv coins have reached the daily limit",
				"id": "IV coin gratis telah mencapai batas harian",
			}
			ErrorResponse(w, http.StatusInternalServerError, messages, lang)
			return
		}
	}

	if err = h.IVCoinRepositories.UpdateIVCoin(ivCoin); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Failed to update iv coin",
			"id": "Gagal mengupdate iv coin",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "IV coin updated successfully", ConvertToIVCoinResponse(ivCoin))
}
