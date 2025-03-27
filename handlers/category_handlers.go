package handlers

import (
	"encoding/json"
	category_dto "iv_project/dto/category"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type categoryHandlers struct {
	CategoryRepositories repositories.CategoryRepositories
}

func CategoryHandler(CategoryRepositories repositories.CategoryRepositories) *categoryHandlers {
	return &categoryHandlers{CategoryRepositories}
}

func ConvertToCategoryResponse(category *models.Category) category_dto.CategoryResponse {
	categoryResponse := category_dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}

	return categoryResponse
}

func (h *categoryHandlers) CreateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() && role != models.UserRoleAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	var request category_dto.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Failed to parse request: invalid JSON format")
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	category := &models.Category{
		Name: request.Name,
	}

	if err := h.CategoryRepositories.CreateCategory(category); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Error occurred while creating category. Please try again later.")
		return
	}

	SuccessResponse(w, http.StatusCreated, "Category created successfully", ConvertToCategoryResponse(category))
}

func (h *categoryHandlers) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid category ID format. Please provide a numeric ID.")
		return
	}

	category, err := h.CategoryRepositories.GetCategoryByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No category found with the provided ID.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Category retrieved successfully", ConvertToCategoryResponse(category))
}

func (h *categoryHandlers) GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	categories, err := h.CategoryRepositories.GetCategories()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while fetching categories.")
		return
	}

	var categoryResponses []category_dto.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, ConvertToCategoryResponse(&category))
	}

	if len(categories) == 0 {
		SuccessResponse(w, http.StatusOK, "No categories available at the moment.", categoryResponses)
		return
	}

	SuccessResponse(w, http.StatusOK, "Categories retrieved successfully", categoryResponses)
}

func (h *categoryHandlers) UpdateCategoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() && role != models.UserRoleAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	var request category_dto.CategoryRequest
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
		ErrorResponse(w, http.StatusBadRequest, "Invalid category ID format. Please provide a numeric ID.")
		return
	}

	category, err := h.CategoryRepositories.GetCategoryByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No category found with the provided ID.")
		return
	}

	if request.Name != "" {
		category.Name = request.Name
	}

	if err := h.CategoryRepositories.UpdateCategory(category); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while updating the category.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Category updated successfully", ConvertToCategoryResponse(category))
}

func (h *categoryHandlers) DeleteCategoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() && role != models.UserRoleAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid category ID format. Please provide a numeric ID.")
		return
	}

	if _, err = h.CategoryRepositories.GetCategoryByID(uint(id)); err != nil {
		ErrorResponse(w, http.StatusNotFound, "No category found with the provided ID.")
		return
	}

	if err := h.CategoryRepositories.DeleteCategory(uint(id)); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while deleting the category.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Category deleted successfully", nil)
}
