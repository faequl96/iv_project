package handlers

import (
	"encoding/json"
	category_dto "iv_project/dto/category"
	discount_category_dto "iv_project/dto/discount_category"
	invitation_theme_dto "iv_project/dto/invitation_theme"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type invitationThemeHandlers struct {
	InvitationThemeRepositories  repositories.InvitationThemeRepositories
	CategoryRepositories         repositories.CategoryRepositories
	DiscountCategoryRepositories repositories.DiscountCategoryRepositories
}

func InvitationThemeHandler(
	InvitationThemeRepositories repositories.InvitationThemeRepositories,
	CategoryRepositories repositories.CategoryRepositories,
	DiscountCategoryRepositories repositories.DiscountCategoryRepositories,
) *invitationThemeHandlers {
	return &invitationThemeHandlers{InvitationThemeRepositories, CategoryRepositories, DiscountCategoryRepositories}
}

func ConvertToInvitationThemeResponse(invitationTheme *models.InvitationTheme) invitation_theme_dto.InvitationThemeResponse {
	invitationThemeResponse := invitation_theme_dto.InvitationThemeResponse{
		ID:               invitationTheme.ID,
		Name:             invitationTheme.Name,
		IDRPrice:         invitationTheme.IDRPrice,
		IDRDiscountPrice: invitationTheme.IDRDiscountPrice,
		IVCPrice:         invitationTheme.IVCPrice,
		IVCDiscountPrice: invitationTheme.IVCDiscountPrice,
	}

	if len(invitationTheme.Categories) != 0 {
		var categoryResponses []category_dto.CategoryResponse
		for _, category := range invitationTheme.Categories {
			categoryCopy := ConvertToCategoryResponse(&category)
			categoryResponses = append(categoryResponses, categoryCopy)
		}
		invitationThemeResponse.Categories = categoryResponses
	}

	if len(invitationTheme.DiscountCategories) != 0 {
		var discountCategoryResponses []discount_category_dto.DiscountCategoryResponse
		for _, discountCategory := range invitationTheme.DiscountCategories {
			discountCategoryCopy := ConvertToDiscountCategoryResponse(&discountCategory)
			discountCategoryResponses = append(discountCategoryResponses, discountCategoryCopy)
		}
		invitationThemeResponse.DiscountCategories = discountCategoryResponses
	}

	return invitationThemeResponse
}

func (h *invitationThemeHandlers) CreateInvitationTheme(w http.ResponseWriter, r *http.Request) {
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

	var request invitation_theme_dto.CreateInvitationThemeRequest
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

	categories, err := h.CategoryRepositories.GetCategoriesByIDs(request.CategoryIDs)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while fetching categories by IDs.",
			"id": "Terjadi kesalahan saat mengambil kategori berdasarkan IDs.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
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

	invitationTheme := &models.InvitationTheme{
		Name:               request.Name,
		IDRPrice:           request.IDRPrice,
		IDRDiscountPrice:   request.IDRPrice,
		IVCPrice:           request.IVCPrice,
		IVCDiscountPrice:   request.IVCPrice,
		Categories:         categories,
		DiscountCategories: discountCategories,
	}

	if err := h.InvitationThemeRepositories.CreateInvitationTheme(invitationTheme); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Error occurred while creating invitation theme. Please try again later.",
			"id": "Terjadi kesalahan saat membuat tema undangan. Coba lagi nanti.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusCreated, "Invitation theme created successfully", ConvertToInvitationThemeResponse(invitationTheme))
}

func (h *invitationThemeHandlers) GetInvitationThemeByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid Invitation theme ID format. Please provide a numeric ID.",
			"id": "Format ID tema undangan tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	invitationTheme, err := h.InvitationThemeRepositories.GetInvitationThemeByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No invitation theme found with the provided ID.",
			"id": "Tema undangan tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitation theme retrieved successfully", ConvertToInvitationThemeResponse(invitationTheme))
}

func (h *invitationThemeHandlers) GetInvitationThemes(w http.ResponseWriter, r *http.Request) {
	invitationThemes, err := h.InvitationThemeRepositories.GetInvitationThemes()
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while fetching invitation themes.",
			"id": "Terjadi kesalahan saat mengambil tema undangan.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	if len(invitationThemes) == 0 {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No invitation themes available at the moment.",
			"id": "Tidak ada tema undangan yang tersedia saat ini.",
		}
		SuccessResponse(w, http.StatusOK, messages[lang], []invitation_theme_dto.InvitationThemeResponse{})
		return
	}

	var invitationThemeResponses []invitation_theme_dto.InvitationThemeResponse
	for _, invitationTheme := range invitationThemes {
		invitationThemeResponses = append(invitationThemeResponses, ConvertToInvitationThemeResponse(&invitationTheme))
	}

	SuccessResponse(w, http.StatusOK, "Invitation themes retrieved successfully", invitationThemeResponses)
}

func (h *invitationThemeHandlers) GetInvitationThemesByCategoryID(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.Atoi(mux.Vars(r)["categoryId"])
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid Invitation theme ID format. Please provide a numeric ID.",
			"id": "Format ID tema undangan tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	invitationThemes, err := h.InvitationThemeRepositories.GetInvitationThemesByCategoryID(uint(categoryID))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while fetching invitation themes by category ID.",
			"id": "Terjadi kesalahan saat mengambil tema undangan berdasarkan ID kategori.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	if len(invitationThemes) == 0 {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No invitation themes found for the specified category.",
			"id": "Tidak ditemukan tema undangan berdasarkan kategori yang dimaksud.",
		}
		SuccessResponse(w, http.StatusOK, messages[lang], []invitation_theme_dto.InvitationThemeResponse{})
		return
	}

	var invitationThemeResponses []invitation_theme_dto.InvitationThemeResponse
	for _, invitationTheme := range invitationThemes {
		invitationThemeResponses = append(invitationThemeResponses, ConvertToInvitationThemeResponse(&invitationTheme))
	}

	SuccessResponse(w, http.StatusOK, "Invitation themes retrieved successfully", invitationThemeResponses)
}

func (h *invitationThemeHandlers) UpdateInvitationThemeByID(w http.ResponseWriter, r *http.Request) {
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

	var request invitation_theme_dto.UpdateInvitationThemeRequest
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
			"en": "Invalid Invitation theme ID format. Please provide a numeric ID.",
			"id": "Format ID tema undangan tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	invitationTheme, err := h.InvitationThemeRepositories.GetInvitationThemeByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No invitation theme found with the provided ID.",
			"id": "Tema undangan tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	categories, err := h.CategoryRepositories.GetCategoriesByIDs(request.CategoryIDs)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while fetching categories by IDs.",
			"id": "Terjadi kesalahan saat mengambil kategori berdasarkan IDs.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
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
		invitationTheme.Name = request.Name
	}
	if request.IDRPrice != 0 {
		invitationTheme.IDRPrice = request.IDRPrice
		invitationTheme.IDRDiscountPrice = request.IDRPrice
	}
	if request.IVCPrice != 0 {
		invitationTheme.IVCPrice = request.IVCPrice
		invitationTheme.IVCDiscountPrice = request.IVCPrice
	}
	if len(categories) != 0 {
		invitationTheme.Categories = categories
	}
	if len(discountCategories) != 0 {
		invitationTheme.DiscountCategories = discountCategories
	}

	if err := h.InvitationThemeRepositories.UpdateInvitationTheme(invitationTheme); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while updating the invitation theme.",
			"id": "Terjadi kesalahan saat mengupdate tema undangan.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitation theme updated successfully", ConvertToInvitationThemeResponse(invitationTheme))
}

func (h *invitationThemeHandlers) DeleteInvitationThemeByID(w http.ResponseWriter, r *http.Request) {
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
			"en": "Invalid Invitation theme ID format. Please provide a numeric ID.",
			"id": "Format ID tema undangan tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	if _, err = h.InvitationThemeRepositories.GetInvitationThemeByID(uint(id)); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No invitation theme found with the provided ID.",
			"id": "Tema undangan tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if err := h.InvitationThemeRepositories.DeleteInvitationTheme(uint(id)); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while deleting the invitation theme.",
			"id": "Terjadi kesalahan saat menghapus tema undangan.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitation theme deleted successfully", nil)
}
