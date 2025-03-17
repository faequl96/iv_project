package handlers

import (
	"encoding/json"
	invitation_theme_dto "iv_project/dto/invitation_theme"
	"iv_project/models"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type invitationThemeHandlers struct {
	InvitationThemeRepositories repositories.InvitationThemeRepositories
}

// InvitationThemeHandler initializes the handler with the given repository.
func InvitationThemeHandler(InvitationThemeRepositories repositories.InvitationThemeRepositories) *invitationThemeHandlers {
	return &invitationThemeHandlers{InvitationThemeRepositories}
}

// convertToInvitationThemeResponse maps the InvitationTheme model to DTO
func convertToInvitationThemeResponse(theme *models.InvitationTheme) invitation_theme_dto.InvitationThemeResponse {
	return invitation_theme_dto.InvitationThemeResponse{
		ID:          theme.ID,
		Title:       theme.Title,
		NormalPrice: theme.NormalPrice,
		DiskonPrice: theme.DiskonPrice,
		Category:    theme.Category,
		Reviews:     theme.Reviews,
	}
}

// CreateInvitationTheme handles the creation of a new invitation theme.
func (h *invitationThemeHandlers) CreateInvitationTheme(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode the incoming JSON request body
	var request invitation_theme_dto.CreateInvitationThemeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Failed to parse request: invalid JSON format")
		return
	}

	// Create a new invitation theme instance
	invitationTheme := &models.InvitationTheme{
		Title:       request.Title,
		NormalPrice: request.NormalPrice,
		DiskonPrice: request.DiskonPrice,
		Category:    request.Category,
	}

	// Store the new theme in the database
	if err := h.InvitationThemeRepositories.CreateInvitationTheme(invitationTheme); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Error occurred while creating invitation theme. Please try again later.")
		return
	}

	// Respond with success message
	SuccessResponse(w, http.StatusCreated, "Invitation theme created successfully", convertToInvitationThemeResponse(invitationTheme))
}

// GetInvitationThemeByID retrieves an invitation theme by its ID.
func (h *invitationThemeHandlers) GetInvitationThemeByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the ID from the request URL parameters
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid theme ID format. Please provide a numeric ID.")
		return
	}

	// Retrieve the invitation theme from the database
	invitationTheme, err := h.InvitationThemeRepositories.GetInvitationThemeByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No invitation theme found with the provided ID.")
		return
	}

	// Respond with the retrieved invitation theme
	SuccessResponse(w, http.StatusOK, "Invitation theme retrieved successfully", convertToInvitationThemeResponse(invitationTheme))
}

// GetInvitationThemes retrieves all available invitation themes.
func (h *invitationThemeHandlers) GetInvitationThemes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Fetch all invitation themes
	invitationThemes, err := h.InvitationThemeRepositories.GetInvitationThemes()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while fetching invitation themes.")
		return
	}

	// Convert slice of users ke slice InvitationThemeResponse DTO
	var invitationThemeResponses []invitation_theme_dto.InvitationThemeResponse
	for _, invitationTheme := range invitationThemes {
		invitationThemeResponses = append(invitationThemeResponses, convertToInvitationThemeResponse(&invitationTheme))
	}

	// If no themes are found, return an empty success response
	if len(invitationThemes) == 0 {
		SuccessResponse(w, http.StatusOK, "No invitation themes available at the moment.", invitationThemeResponses)
		return
	}

	// Respond with the retrieved themes
	SuccessResponse(w, http.StatusOK, "Invitation themes retrieved successfully", invitationThemeResponses)
}

// GetInvitationThemesByCategory retrieves invitation themes filtered by category.
func (h *invitationThemeHandlers) GetInvitationThemesByCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the category from the URL parameters
	category := mux.Vars(r)["category"]

	// Fetch themes by category
	invitationThemes, err := h.InvitationThemeRepositories.GetInvitationThemesByCategory(category)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while fetching themes by category.")
		return
	}

	// Convert slice of users ke slice InvitationThemeResponse DTO
	var invitationThemeResponses []invitation_theme_dto.InvitationThemeResponse
	for _, invitationTheme := range invitationThemes {
		invitationThemeResponses = append(invitationThemeResponses, convertToInvitationThemeResponse(&invitationTheme))
	}

	// If no themes are found, return an empty success response
	if len(invitationThemes) == 0 {
		SuccessResponse(w, http.StatusOK, "No invitation themes found for the specified category.", invitationThemeResponses)
		return
	}

	// Respond with the retrieved themes
	SuccessResponse(w, http.StatusOK, "Invitation themes retrieved successfully", invitationThemeResponses)
}

// UpdateInvitationTheme modifies an existing invitation theme.
func (h *invitationThemeHandlers) UpdateInvitationTheme(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the ID from the request URL parameters
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid theme ID format. Please provide a numeric ID.")
		return
	}

	// Fetch the existing invitation theme
	invitationTheme, err := h.InvitationThemeRepositories.GetInvitationThemeByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No invitation theme found with the provided ID.")
		return
	}

	// Decode the incoming JSON request body
	var request invitation_theme_dto.UpdateInvitationThemeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Failed to parse request: invalid JSON format")
		return
	}

	// Update fields if they are provided in the request
	if request.Title != "" {
		invitationTheme.Title = request.Title
	}
	invitationTheme.NormalPrice = request.NormalPrice
	invitationTheme.DiskonPrice = request.DiskonPrice
	if request.Category != "" {
		invitationTheme.Category = request.Category
	}

	// Save the updated invitation theme
	if err := h.InvitationThemeRepositories.UpdateInvitationTheme(invitationTheme); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while updating the invitation theme.")
		return
	}

	// Respond with success message
	SuccessResponse(w, http.StatusOK, "Invitation theme updated successfully", convertToInvitationThemeResponse(invitationTheme))
}

// DeleteInvitationTheme removes an invitation theme by its ID.
func (h *invitationThemeHandlers) DeleteInvitationTheme(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the ID from the request URL parameters
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid theme ID format. Please provide a numeric ID.")
		return
	}

	// Check if the invitation theme exists before attempting to delete
	if _, err = h.InvitationThemeRepositories.GetInvitationThemeByID(uint(id)); err != nil {
		ErrorResponse(w, http.StatusNotFound, "No invitation theme found with the provided ID.")
		return
	}

	// Delete the invitation theme from the database
	if err := h.InvitationThemeRepositories.DeleteInvitationTheme(uint(id)); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while deleting the invitation theme.")
		return
	}

	// Respond with success message
	SuccessResponse(w, http.StatusOK, "Invitation theme deleted successfully", nil)
}
