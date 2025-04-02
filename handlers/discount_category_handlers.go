package handlers

import (
	"encoding/json"
	discount_category_dto "iv_project/dto/discount_category"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type discountCategoryHandlers struct {
	DiscountCategoryRepositories repositories.DiscountCategoryRepositories
}

func DiscountCategoryHandler(DiscountCategoryRepositories repositories.DiscountCategoryRepositories) *discountCategoryHandlers {
	return &discountCategoryHandlers{DiscountCategoryRepositories}
}

func ConvertToDiscountCategoryResponse(discountCategory *models.DiscountCategory) discount_category_dto.DiscountCategoryResponse {
	discountCategoryResponse := discount_category_dto.DiscountCategoryResponse{
		ID:   discountCategory.ID,
		Name: discountCategory.Name,
	}

	return discountCategoryResponse
}

func (h *discountCategoryHandlers) CreateDiscountCategory(w http.ResponseWriter, r *http.Request) {
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

	var request discount_category_dto.DiscountCategoryRequest
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

	discountCategory := &models.DiscountCategory{
		Name: request.Name,
	}

	if err := h.DiscountCategoryRepositories.CreateDiscountCategory(discountCategory); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Error occurred while creating discount category. Please try again later.",
			"id": "Terjadi kesalahan saat membuat kategori diskon. Coba lagi nanti.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusCreated, "Discount category created successfully", ConvertToDiscountCategoryResponse(discountCategory))
}

func (h *discountCategoryHandlers) GetDiscountCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid discount category ID format. Please provide a numeric ID.",
			"id": "Format ID kategori diskon tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	discountCategory, err := h.DiscountCategoryRepositories.GetDiscountCategoryByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No discount category found with the provided ID.",
			"id": "Kategori diskon tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Discount category retrieved successfully", ConvertToDiscountCategoryResponse(discountCategory))
}

func (h *discountCategoryHandlers) GetDiscountCategories(w http.ResponseWriter, r *http.Request) {
	discountCategories, err := h.DiscountCategoryRepositories.GetDiscountCategories()
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while fetching discount categories.",
			"id": "Terjadi kesalahan saat mengambil kategori diskon.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	var discountCategoryResponses []discount_category_dto.DiscountCategoryResponse
	for _, discountCategory := range discountCategories {
		discountCategoryResponses = append(discountCategoryResponses, ConvertToDiscountCategoryResponse(&discountCategory))
	}

	if len(discountCategories) == 0 {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No discount categories available at the moment.",
			"id": "Tidak ada kategori diskon yang tersedia saat ini.",
		}
		SuccessResponse(w, http.StatusOK, messages[lang], discountCategoryResponses)
		return
	}

	SuccessResponse(w, http.StatusOK, "Discount categories retrieved successfully", discountCategoryResponses)
}

func (h *discountCategoryHandlers) UpdateDiscountCategory(w http.ResponseWriter, r *http.Request) {
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

	var request discount_category_dto.DiscountCategoryRequest
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
			"en": "Invalid discount category ID format. Please provide a numeric ID.",
			"id": "Format ID kategori diskon tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	discountCategory, err := h.DiscountCategoryRepositories.GetDiscountCategoryByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No discount category found with the provided ID.",
			"id": "Kategori diskon tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if request.Name != "" {
		discountCategory.Name = request.Name
	}

	if err := h.DiscountCategoryRepositories.UpdateDiscountCategory(discountCategory); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while updating the discount category.",
			"id": "Terjadi kesalahan saat memperbarui kategori diskon.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Discount category updated successfully", ConvertToDiscountCategoryResponse(discountCategory))
}

func (h *discountCategoryHandlers) DeleteDiscountCategory(w http.ResponseWriter, r *http.Request) {
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
			"en": "Invalid discount category ID format. Please provide a numeric ID.",
			"id": "Format ID kategori diskon tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	if _, err = h.DiscountCategoryRepositories.GetDiscountCategoryByID(uint(id)); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No discount category found with the provided ID.",
			"id": "Kategori diskon tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if err := h.DiscountCategoryRepositories.DeleteDiscountCategory(uint(id)); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while deleting the discount category.",
			"id": "Terjadi kesalahan saat menghapus kategori diskon.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Discount category deleted successfully", nil)
}
