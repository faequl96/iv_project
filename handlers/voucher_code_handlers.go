package handlers

import (
	"encoding/json"
	user_dto "iv_project/dto/user"
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
	UserRepositories        repositories.UserRepositories
}

func VoucherCodeHandler(
	VoucherCodeRepositories repositories.VoucherCodeRepositories,
	UserRepositories repositories.UserRepositories,
) *voucherCodeHandlers {
	return &voucherCodeHandlers{VoucherCodeRepositories, UserRepositories}
}

func ConvertToVoucherCodeResponse(voucherCode *models.VoucherCode) voucher_code_dto.VoucherCodeResponse {
	voucherCodeResponse := voucher_code_dto.VoucherCodeResponse{
		ID:                 voucherCode.ID,
		Name:               voucherCode.Name,
		DiscountPercentage: voucherCode.DiscountPercentage,
		UsageLimitPerUser:  voucherCode.UsageLimitPerUser,
		IsGlobal:           voucherCode.IsGlobal,
	}

	if len(voucherCode.Users) != 0 {
		var userResponses []user_dto.UserResponse
		for _, user := range voucherCode.Users {
			userCopy := ConvertToUserResponse(&user)
			userResponses = append(userResponses, userCopy)
		}
		voucherCodeResponse.Users = userResponses
	} else {
		voucherCodeResponse.Users = []user_dto.UserResponse{}
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

	users, err := h.UserRepositories.GetUsersByIDs(request.UserIDs)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while fetching users by IDs.",
			"id": "Terjadi kesalahan saat mengambil user berdasarkan IDs.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	voucherCode := &models.VoucherCode{
		Name:               request.Name,
		DiscountPercentage: request.DiscountPercentage,
		UsageLimitPerUser:  request.UsageLimitPerUser,
		IsGlobal:           request.IsGlobal,
		Users:              users,
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

func (h *voucherCodeHandlers) GetVoucherCodeByName(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	voucherCode, err := h.VoucherCodeRepositories.GetVoucherCodeByName(name)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No voucher code found",
			"id": "Kode voucher tidak ditemukan",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Voucher code retrieved successfully", ConvertToVoucherCodeResponse(voucherCode))
}

func (h *voucherCodeHandlers) GetVoucherCodes(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "You do not have permission to access this resource.",
			"id": "Anda tidak memiliki izin untuk mengakses sumber daya ini.",
		}
		ErrorResponse(w, http.StatusForbidden, messages, lang)
		return
	}

	voucherCodes, err := h.VoucherCodeRepositories.GetVoucherCodes()
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while fetching voucher codes.",
			"id": "Terjadi kesalahan saat mengambil kode voucher.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	if len(voucherCodes) == 0 {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No voucher codes available at the moment.",
			"id": "Tidak ada kode voucher yang tersedia saat ini.",
		}
		SuccessResponse(w, http.StatusOK, messages[lang], []voucher_code_dto.VoucherCodeResponse{})
		return
	}

	var responses []voucher_code_dto.VoucherCodeResponse
	for _, voucherCode := range voucherCodes {
		responses = append(responses, ConvertToVoucherCodeResponse(&voucherCode))
	}

	SuccessResponse(w, http.StatusOK, "Reviews retrieved successfully", responses)
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

	users, err := h.UserRepositories.GetUsersByIDs(request.UserIDs)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while fetching users by IDs.",
			"id": "Terjadi kesalahan saat mengambil user berdasarkan IDs.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	if request.Name != "" && request.DiscountPercentage != 0 {
		voucherCode.Name = request.Name
		voucherCode.DiscountPercentage = request.DiscountPercentage
	}
	voucherCode.UsageLimitPerUser = request.UsageLimitPerUser
	voucherCode.IsGlobal = request.IsGlobal
	voucherCode.Users = users

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
