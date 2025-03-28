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
	w.Header().Set("Content-Type", "application/json")

	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	var request invitation_theme_dto.CreateInvitationThemeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Failed to parse request: invalid JSON format")
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	categories, err := h.CategoryRepositories.GetCategoriesByIDs(request.CategoryIDs)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while fetching categories by ids.")
		return
	}

	discountCategories, err := h.DiscountCategoryRepositories.GetDiscountCategoriesByIDs(request.DiscountCategoryIDs)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while fetching discount categories by ids.")
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
		ErrorResponse(w, http.StatusInternalServerError, "Error occurred while creating invitation theme. Please try again later.")
		return
	}

	SuccessResponse(w, http.StatusCreated, "Invitation theme created successfully", ConvertToInvitationThemeResponse(invitationTheme))
}

func (h *invitationThemeHandlers) GetInvitationThemeByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid invitation theme ID format. Please provide a numeric ID.")
		return
	}

	invitationTheme, err := h.InvitationThemeRepositories.GetInvitationThemeByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No invitation theme found with the provided ID.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitation theme retrieved successfully", ConvertToInvitationThemeResponse(invitationTheme))
}

func (h *invitationThemeHandlers) GetInvitationThemes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	invitationThemes, err := h.InvitationThemeRepositories.GetInvitationThemes()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while fetching invitation themes.")
		return
	}

	var invitationThemeResponses []invitation_theme_dto.InvitationThemeResponse
	for _, invitationTheme := range invitationThemes {
		invitationThemeResponses = append(invitationThemeResponses, ConvertToInvitationThemeResponse(&invitationTheme))
	}

	if len(invitationThemes) == 0 {
		SuccessResponse(w, http.StatusOK, "No invitation themes available at the moment.", invitationThemeResponses)
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitation themes retrieved successfully", invitationThemeResponses)
}

func (h *invitationThemeHandlers) GetInvitationThemesByCategoryID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	categoryID, err := strconv.Atoi(mux.Vars(r)["categoryId"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid category ID format. Please provide a numeric ID.")
		return
	}

	invitationThemes, err := h.InvitationThemeRepositories.GetInvitationThemesByCategoryID(uint(categoryID))
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while fetching invitation themes by category.")
		return
	}

	var invitationThemeResponses []invitation_theme_dto.InvitationThemeResponse
	for _, invitationTheme := range invitationThemes {
		invitationThemeResponses = append(invitationThemeResponses, ConvertToInvitationThemeResponse(&invitationTheme))
	}

	if len(invitationThemes) == 0 {
		SuccessResponse(w, http.StatusOK, "No invitation themes found for the specified category.", invitationThemeResponses)
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitation themes retrieved successfully", invitationThemeResponses)
}

func (h *invitationThemeHandlers) UpdateInvitationThemeByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	var request invitation_theme_dto.UpdateInvitationThemeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Failed to parse request: invalid JSON format")
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid invitation theme ID format. Please provide a numeric ID.")
		return
	}

	invitationTheme, err := h.InvitationThemeRepositories.GetInvitationThemeByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No invitation theme found with the provided ID.")
		return
	}

	categories, err := h.CategoryRepositories.GetCategoriesByIDs(request.CategoryIDs)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while fetching categories by ids.")
		return
	}

	discountCategories, err := h.DiscountCategoryRepositories.GetDiscountCategoriesByIDs(request.DiscountCategoryIDs)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while fetching discount categories by ids.")
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
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while updating the invitation theme.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitation theme updated successfully", ConvertToInvitationThemeResponse(invitationTheme))
}

func (h *invitationThemeHandlers) DeleteInvitationThemeByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid invitation theme ID format. Please provide a numeric ID.")
		return
	}

	if _, err = h.InvitationThemeRepositories.GetInvitationThemeByID(uint(id)); err != nil {
		ErrorResponse(w, http.StatusNotFound, "No invitation theme found with the provided ID.")
		return
	}

	if err := h.InvitationThemeRepositories.DeleteInvitationTheme(uint(id)); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while deleting the invitation theme.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitation theme deleted successfully", nil)
}
