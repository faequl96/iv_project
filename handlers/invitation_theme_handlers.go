package handlers

import (
	"encoding/json"
	invitation_theme_dto "iv_project/dto/invitation_theme"
	review_dto "iv_project/dto/review"
	"iv_project/models"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type invitationThemeHandlers struct {
	InvitationThemeRepositories repositories.InvitationThemeRepositories
}

func InvitationThemeHandler(InvitationThemeRepositories repositories.InvitationThemeRepositories) *invitationThemeHandlers {
	return &invitationThemeHandlers{InvitationThemeRepositories}
}

func ConvertToInvitationThemeResponse(invitationTheme *models.InvitationTheme) invitation_theme_dto.InvitationThemeResponse {
	var reviewResponses []review_dto.ReviewResponse
	for _, review := range invitationTheme.Reviews {
		reviewCopy := ConvertToReviewResponse(&review)
		reviewResponses = append(reviewResponses, reviewCopy)
	}

	return invitation_theme_dto.InvitationThemeResponse{
		ID:          invitationTheme.ID,
		Title:       invitationTheme.Title,
		NormalPrice: invitationTheme.NormalPrice,
		DiskonPrice: invitationTheme.DiskonPrice,
		Category:    invitationTheme.Category,
		Reviews:     reviewResponses,
	}
}

func (h *invitationThemeHandlers) CreateInvitationTheme(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request invitation_theme_dto.CreateInvitationThemeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Failed to parse request: invalid JSON format")
		return
	}

	invitationTheme := &models.InvitationTheme{
		Title:       request.Title,
		NormalPrice: request.NormalPrice,
		DiskonPrice: request.DiskonPrice,
		Category:    request.Category,
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

func (h *invitationThemeHandlers) GetInvitationThemesByCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	category := mux.Vars(r)["category"]

	invitationThemes, err := h.InvitationThemeRepositories.GetInvitationThemesByCategory(category)
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

func (h *invitationThemeHandlers) UpdateInvitationTheme(w http.ResponseWriter, r *http.Request) {
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

	var request invitation_theme_dto.UpdateInvitationThemeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Failed to parse request: invalid JSON format")
		return
	}

	if request.Title != "" {
		invitationTheme.Title = request.Title
	}
	invitationTheme.NormalPrice = request.NormalPrice
	invitationTheme.DiskonPrice = request.DiskonPrice
	if request.Category != "" {
		invitationTheme.Category = request.Category
	}

	if err := h.InvitationThemeRepositories.UpdateInvitationTheme(invitationTheme); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while updating the invitation theme.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitation theme updated successfully", ConvertToInvitationThemeResponse(invitationTheme))
}

func (h *invitationThemeHandlers) DeleteInvitationTheme(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
