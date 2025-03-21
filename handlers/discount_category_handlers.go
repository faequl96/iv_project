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
	return discount_category_dto.DiscountCategoryResponse{
		ID:   discountCategory.ID,
		Name: discountCategory.Name,
	}
}

func (h *discountCategoryHandlers) CreateDiscountCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	var request discount_category_dto.CreateDiscountCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Failed to parse request: invalid JSON format")
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	discountCategory := &models.DiscountCategory{
		Name: request.Name,
	}

	if err := h.DiscountCategoryRepositories.CreateDiscountCategory(discountCategory); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Error occurred while creating discount category. Please try again later.")
		return
	}

	SuccessResponse(w, http.StatusCreated, "Discount category created successfully", ConvertToDiscountCategoryResponse(discountCategory))
}

func (h *discountCategoryHandlers) GetDiscountCategoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid discount category ID format. Please provide a numeric ID.")
		return
	}

	discountCategory, err := h.DiscountCategoryRepositories.GetDiscountCategoryByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No discount category found with the provided ID.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Discount category retrieved successfully", ConvertToDiscountCategoryResponse(discountCategory))
}

func (h *discountCategoryHandlers) GetDiscountCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	discountCategories, err := h.DiscountCategoryRepositories.GetDiscountCategories()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while fetching discount categories.")
		return
	}

	var discountCategoryResponses []discount_category_dto.DiscountCategoryResponse
	for _, discountCategory := range discountCategories {
		discountCategoryResponses = append(discountCategoryResponses, ConvertToDiscountCategoryResponse(&discountCategory))
	}

	if len(discountCategories) == 0 {
		SuccessResponse(w, http.StatusOK, "No discount categories available at the moment.", discountCategoryResponses)
		return
	}

	SuccessResponse(w, http.StatusOK, "Discount categories retrieved successfully", discountCategoryResponses)
}

func (h *discountCategoryHandlers) UpdateDiscountCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	var request discount_category_dto.UpdateDiscountCategoryRequest
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
		ErrorResponse(w, http.StatusBadRequest, "Invalid discount category ID format. Please provide a numeric ID.")
		return
	}

	discountCategory, err := h.DiscountCategoryRepositories.GetDiscountCategoryByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No discount category found with the provided ID.")
		return
	}

	if request.Name != "" {
		discountCategory.Name = request.Name
	}

	if err := h.DiscountCategoryRepositories.UpdateDiscountCategory(discountCategory); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while updating the discount category.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Discount category updated successfully", ConvertToDiscountCategoryResponse(discountCategory))
}

func (h *discountCategoryHandlers) DeleteDiscountCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid discount category ID format. Please provide a numeric ID.")
		return
	}

	if _, err = h.DiscountCategoryRepositories.GetDiscountCategoryByID(uint(id)); err != nil {
		ErrorResponse(w, http.StatusNotFound, "No discount category found with the provided ID.")
		return
	}

	if err := h.DiscountCategoryRepositories.DeleteDiscountCategory(uint(id)); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while deleting the discount category.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Discount category deleted successfully", nil)
}
