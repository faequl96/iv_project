package handlers

import (
	"encoding/json"
	discount_category_dto "iv_project/dto/discount_category"
	iv_coin_package_dto "iv_project/dto/iv_coin_package"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type ivCoinPackageHandlers struct {
	IVCoinPackageRepositories    repositories.IVCoinPackageRepositories
	DiscountCategoryRepositories repositories.DiscountCategoryRepositories
}

func IVCoinPackageHandler(
	IVCoinPackageRepositories repositories.IVCoinPackageRepositories,
	DiscountCategoryRepositories repositories.DiscountCategoryRepositories,
) *ivCoinPackageHandlers {
	return &ivCoinPackageHandlers{IVCoinPackageRepositories, DiscountCategoryRepositories}
}

func ConvertToIVCoinPackageResponse(ivCoinPackage *models.IVCoinPackage) iv_coin_package_dto.IVCoinPackageResponse {
	ivCoinPackageResponse := iv_coin_package_dto.IVCoinPackageResponse{
		ID:               ivCoinPackage.ID,
		Name:             ivCoinPackage.Name,
		CoinAmount:       ivCoinPackage.CoinAmount,
		IDRPrice:         ivCoinPackage.IDRPrice,
		IDRDiscountPrice: ivCoinPackage.IDRDiscountPrice,
	}

	if len(ivCoinPackage.DiscountCategories) != 0 {
		var discountCategoryResponses []discount_category_dto.DiscountCategoryResponse
		for _, discountCategory := range ivCoinPackage.DiscountCategories {
			discountCategoryCopy := ConvertToDiscountCategoryResponse(&discountCategory)
			discountCategoryResponses = append(discountCategoryResponses, discountCategoryCopy)
		}
		ivCoinPackageResponse.DiscountCategories = discountCategoryResponses
	} else {
		ivCoinPackageResponse.DiscountCategories = []discount_category_dto.DiscountCategoryResponse{}
	}

	return ivCoinPackageResponse
}

func (h *ivCoinPackageHandlers) CreateIVCoinPackage(w http.ResponseWriter, r *http.Request) {
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

	var request iv_coin_package_dto.CreateIVCoinPackageRequest
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

	discountCategories, err := h.DiscountCategoryRepositories.GetDiscountCategoriesByIDs(request.DiscountCategoryIDs)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while fetching discount categories by IDs.",
			"id": "Terjadi kesalahan saat mengambil kategori diskon berdasarkan IDs.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	ivCoinPackage := &models.IVCoinPackage{
		Name:               request.Name,
		CoinAmount:         request.CoinAmount,
		IDRPrice:           request.IDRPrice,
		IDRDiscountPrice:   request.IDRPrice,
		DiscountCategories: discountCategories,
	}

	if err := h.IVCoinPackageRepositories.CreateIVCoinPackage(ivCoinPackage); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Error occurred while creating iv coin package. Please try again later.",
			"id": "Terjadi kesalahan saat membuat paket iv coin. Coba lagi nanti.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusCreated, "IV coin package created successfully", ConvertToIVCoinPackageResponse(ivCoinPackage))
}

func (h *ivCoinPackageHandlers) GetIVCoinPackageByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid iv coin package ID format. Please provide a numeric ID.",
			"id": "Format ID paket iv coin tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	ivCoinPackage, err := h.IVCoinPackageRepositories.GetIVCoinPackageByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No iv coin package found with the provided ID.",
			"id": "Paket iv coin tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "IV coin package retrieved successfully", ConvertToIVCoinPackageResponse(ivCoinPackage))
}

func (h *ivCoinPackageHandlers) GetIVCoinPackages(w http.ResponseWriter, r *http.Request) {
	ivCoinPackages, err := h.IVCoinPackageRepositories.GetIVCoinPackages()
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while fetching iv coin packages.",
			"id": "Terjadi kesalahan saat mengambil paket iv coin.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	if len(ivCoinPackages) == 0 {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No iv coin packages available at the moment.",
			"id": "Tidak ada paket iv coin yang tersedia saat ini.",
		}
		SuccessResponse(w, http.StatusOK, messages[lang], []iv_coin_package_dto.IVCoinPackageResponse{})
		return
	}

	var responses []iv_coin_package_dto.IVCoinPackageResponse
	for _, ivCoinPackage := range ivCoinPackages {
		responses = append(responses, ConvertToIVCoinPackageResponse(&ivCoinPackage))
	}

	SuccessResponse(w, http.StatusOK, "IV coin packages retrieved successfully", responses)
}

func (h *ivCoinPackageHandlers) UpdateIVCoinPackageByID(w http.ResponseWriter, r *http.Request) {
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

	var request iv_coin_package_dto.UpdateIVCoinPackageRequest
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
			"en": "Invalid iv coin package ID format. Please provide a numeric ID.",
			"id": "Format ID paket iv coin tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	ivCoinPackage, err := h.IVCoinPackageRepositories.GetIVCoinPackageByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No iv coin package found with the provided ID.",
			"id": "Paket iv coin tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	discountCategories, err := h.DiscountCategoryRepositories.GetDiscountCategoriesByIDs(request.DiscountCategoryIDs)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while fetching discount categories by IDs.",
			"id": "Terjadi kesalahan saat mengambil kategori diskon berdasarkan IDs.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	if request.Name != "" {
		ivCoinPackage.Name = request.Name
	}
	if request.CoinAmount != 0 {
		ivCoinPackage.CoinAmount = request.CoinAmount
	}
	if request.IDRPrice != 0 {
		ivCoinPackage.IDRPrice = request.IDRPrice
		ivCoinPackage.IDRDiscountPrice = request.IDRPrice
	}
	if len(discountCategories) != 0 {
		ivCoinPackage.DiscountCategories = discountCategories
	}

	if err := h.IVCoinPackageRepositories.UpdateIVCoinPackage(ivCoinPackage); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while updating the iv coin package.",
			"id": "Terjadi kesalahan saat mengupdate paket iv coin.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "IV coin package updated successfully", ConvertToIVCoinPackageResponse(ivCoinPackage))
}

func (h *ivCoinPackageHandlers) DeleteIVCoinPackageByID(w http.ResponseWriter, r *http.Request) {
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

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid iv coin package ID format. Please provide a numeric ID.",
			"id": "Format ID paket iv coin tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	if _, err = h.IVCoinPackageRepositories.GetIVCoinPackageByID(uint(id)); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No iv coin package found with the provided ID.",
			"id": "Paket iv coin tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if err := h.IVCoinPackageRepositories.DeleteIVCoinPackage(uint(id)); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while deleting the iv coin package.",
			"id": "Terjadi kesalahan saat menghapus paket iv coin.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "IV coin package deleted successfully", nil)
}
