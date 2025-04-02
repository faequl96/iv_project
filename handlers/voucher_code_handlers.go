package handlers

import (
	"encoding/json"
	voucher_code_dto "iv_project/dto/voucher_code"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type voucherCodeHandlers struct {
	VoucherCodeRepositories repositories.VoucherCodeRepositories
}

func VoucherCodeHandler(VoucherCodeRepositories repositories.VoucherCodeRepositories) *voucherCodeHandlers {
	return &voucherCodeHandlers{VoucherCodeRepositories}
}

func ConvertToVoucherCodeResponse(VoucherCode *models.VoucherCode) voucher_code_dto.VoucherCodeResponse {
	voucherCodeResponse := voucher_code_dto.VoucherCodeResponse{
		ID:                 VoucherCode.ID,
		Name:               VoucherCode.Name,
		DiscountPercentage: VoucherCode.DiscountPercentage,
	}

	return voucherCodeResponse
}

func (h *voucherCodeHandlers) CreateVoucherCode(w http.ResponseWriter, r *http.Request) {
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

	var request voucher_code_dto.VoucherCodeRequest
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

	voucherCode := &models.VoucherCode{
		Name:               request.Name,
		DiscountPercentage: request.DiscountPercentage,
	}

	if err := h.VoucherCodeRepositories.CreateVoucherCode(voucherCode); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Error occurred while creating voucher code. Please try again later.",
			"id": "Terjadi kesalahan saat membuat kode voucher. Coba lagi nanti.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusCreated, "Voucher code created successfully", ConvertToVoucherCodeResponse(voucherCode))
}

func (h *voucherCodeHandlers) GetVoucherCodeByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid voucher code ID format. Please provide a numeric ID.",
			"id": "Format ID kode voucher tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	voucherCode, err := h.VoucherCodeRepositories.GetVoucherCodeByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No voucher code found with the provided ID.",
			"id": "Kode voucher tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Voucher code retrieved successfully", ConvertToVoucherCodeResponse(voucherCode))
}

func (h *voucherCodeHandlers) UpdateVoucherCodeByID(w http.ResponseWriter, r *http.Request) {
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

	var request voucher_code_dto.VoucherCodeRequest
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

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid voucher code ID format. Please provide a numeric ID.",
			"id": "Format ID kode voucher tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	voucherCode, err := h.VoucherCodeRepositories.GetVoucherCodeByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No voucher code found with the provided ID.",
			"id": "Kode voucher tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if request.Name != "" && request.DiscountPercentage != 0 {
		voucherCode.Name = request.Name
		voucherCode.DiscountPercentage = request.DiscountPercentage
	}

	if err := h.VoucherCodeRepositories.UpdateVoucherCode(voucherCode); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while updating the voucher code.",
			"id": "Terjadi kesalahan saat mengupdate kode voucher.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Voucher code updated successfully", ConvertToVoucherCodeResponse(voucherCode))
}

func (h *voucherCodeHandlers) DeleteVoucherCodeByID(w http.ResponseWriter, r *http.Request) {
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
			"en": "Invalid voucher code ID format. Please provide a numeric ID.",
			"id": "Format ID kode voucher tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	if _, err = h.VoucherCodeRepositories.GetVoucherCodeByID(uint(id)); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No voucher code found with the provided ID.",
			"id": "Kode voucher tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if err := h.VoucherCodeRepositories.DeleteVoucherCode(uint(id)); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while deleting the voucher code.",
			"id": "Terjadi kesalahan saat menghapus kode voucher.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Voucher code deleted successfully", nil)
}
