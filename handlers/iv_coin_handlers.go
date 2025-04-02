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
	ivCoinResponse := iv_coin_dto.IVCoinResponse{
		ID:      ivCoin.ID,
		Balance: ivCoin.Balance,
	}

	return ivCoinResponse
}

func (h *ivCoinHandlers) GetIVCoin(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIdKey).(string)
	iVCoin, err := h.IVCoinRepositories.GetIVCoinByUserID(userID)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No iv coin found with the provided user.",
			"id": "IV coin tidak ditemukan dengan pengguna yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "IV coin retrieved successfully", ConvertToIVCoinResponse(iVCoin))
}

func (h *ivCoinHandlers) GetIVCoinByID(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() && role != models.UserRoleAdmin.String() {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "You do not have permission to access this resource.",
			"id": "Anda tidak memiliki izin untuk mengakses sumber daya ini.",
		}
		ErrorResponse(w, http.StatusForbidden, messages, lang)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid iv coin ID format. Please provide a numeric ID.",
			"id": "Format ID iv coin tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	ivCoin, err := h.IVCoinRepositories.GetIVCoinByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No iv coin found with the provided ID.",
			"id": "IV coin tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "IV coin retrieved successfully", ConvertToIVCoinResponse(ivCoin))
}

func (h *ivCoinHandlers) UpdateIVCoinByID(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() && role != models.UserRoleAdmin.String() {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "You do not have permission to access this resource.",
			"id": "Anda tidak memiliki izin untuk mengakses sumber daya ini.",
		}
		ErrorResponse(w, http.StatusForbidden, messages, lang)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid iv coin ID format. Please provide a numeric ID.",
			"id": "Format ID iv coin tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	ivCoin, err := h.IVCoinRepositories.GetIVCoinByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No iv coin found with the provided ID.",
			"id": "IV coin tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	var request iv_coin_dto.IVCoinRequest
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

	ivCoin.Balance = request.Balance

	if err = h.IVCoinRepositories.UpdateIVCoin(ivCoin); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Failed to update iv coin.",
			"id": "Gagal mengupdate iv coin.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "IV coin updated successfully", ConvertToIVCoinResponse(ivCoin))
}
