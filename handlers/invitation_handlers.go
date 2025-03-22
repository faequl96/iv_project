package handlers

import (
	"encoding/json"
	invitation_dto "iv_project/dto/invitation"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type invitationHandlers struct {
	InvitationRepositories repositories.InvitationRepositories
}

func InvitationHandler(InvitationRepositories repositories.InvitationRepositories) *invitationHandlers {
	return &invitationHandlers{InvitationRepositories}
}

func ConvertToInvitationResponse(invitation *models.Invitation) invitation_dto.InvitationResponse {
	invitationResponse := invitation_dto.InvitationResponse{
		ID:     invitation.ID,
		Status: invitation.Status,
	}

	if invitation.InvitationTheme != nil {
		invitationThemeResponse := ConvertToInvitationThemeResponse(invitation.InvitationTheme)
		invitationResponse.InvitationTheme = &invitationThemeResponse
	}

	if invitation.InvitationData != nil {
		invitationDataResponse := ConvertToInvitationDataResponse(invitation.InvitationData)
		invitationResponse.InvitationData = &invitationDataResponse
		if invitation.InvitationData.Gallery != nil {
			galleryResponse := ConvertToGalleryResponse(invitation.InvitationData.Gallery)
			invitationResponse.InvitationData.Gallery = &galleryResponse
		}
	}

	return invitationResponse
}

func (h *invitationHandlers) CreateInvitation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request invitation_dto.CreateInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid JSON format.")
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	eventDate, err := time.Parse(time.RFC3339, request.InvitationData.EventDate)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid EventDate format. Use RFC3339.")
		return
	}

	invitation := &models.Invitation{
		UserID: request.UserID,
		Status: models.InvitationStatusDraft,
		InvitationData: &models.InvitationData{
			EventName: request.InvitationData.EventName,
			EventDate: eventDate,
			Location:  request.InvitationData.Location,
		},
	}

	uploadedFiles, ok := r.Context().Value(middleware.UploadsKey).(map[string]string)
	if ok {
		if val, exists := uploadedFiles["main_image_url"]; exists {
			invitation.InvitationData.MainImageURL = val
		}

		if invitation.InvitationData.Gallery != nil {
			if val, exists := uploadedFiles["image_url_1"]; exists {
				invitation.InvitationData.Gallery.ImageURL1 = val
			}
			if val, exists := uploadedFiles["image_url_2"]; exists {
				invitation.InvitationData.Gallery.ImageURL2 = val
			}
			if val, exists := uploadedFiles["image_url_3"]; exists {
				invitation.InvitationData.Gallery.ImageURL3 = val
			}
			if val, exists := uploadedFiles["image_url_4"]; exists {
				invitation.InvitationData.Gallery.ImageURL4 = val
			}
			if val, exists := uploadedFiles["image_url_5"]; exists {
				invitation.InvitationData.Gallery.ImageURL5 = val
			}
			if val, exists := uploadedFiles["image_url_6"]; exists {
				invitation.InvitationData.Gallery.ImageURL6 = val
			}
			if val, exists := uploadedFiles["image_url_7"]; exists {
				invitation.InvitationData.Gallery.ImageURL7 = val
			}
			if val, exists := uploadedFiles["image_url_8"]; exists {
				invitation.InvitationData.Gallery.ImageURL8 = val
			}
			if val, exists := uploadedFiles["image_url_9"]; exists {
				invitation.InvitationData.Gallery.ImageURL9 = val
			}
			if val, exists := uploadedFiles["image_url_10"]; exists {
				invitation.InvitationData.Gallery.ImageURL10 = val
			}
			if val, exists := uploadedFiles["image_url_11"]; exists {
				invitation.InvitationData.Gallery.ImageURL11 = val
			}
			if val, exists := uploadedFiles["image_url_12"]; exists {
				invitation.InvitationData.Gallery.ImageURL12 = val
			}
		}
	}

	err = h.InvitationRepositories.CreateInvitation(invitation)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to create invitation.")
		return
	}

	SuccessResponse(w, http.StatusCreated, "Invitation created successfully", ConvertToInvitationResponse(invitation))
}

func (h *invitationHandlers) GetInvitationByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid invitation ID format. Please provide a numeric ID.")
		return
	}

	invitation, err := h.InvitationRepositories.GetInvitationByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No invitation found with the provided ID.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitation retrieved successfully", ConvertToInvitationResponse(invitation))
}

func (h *invitationHandlers) GetInvitations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	invitations, err := h.InvitationRepositories.GetInvitations()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while fetching invitations.")
		return
	}

	var invitationResponses []invitation_dto.InvitationResponse
	for _, invitation := range invitations {
		invitationResponses = append(invitationResponses, ConvertToInvitationResponse(&invitation))
	}

	if len(invitations) == 0 {
		SuccessResponse(w, http.StatusOK, "No invitations available at the moment.", invitationResponses)
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitations retrieved successfully", invitationResponses)
}

func (h *invitationHandlers) GetInvitationsByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["userId"]

	invitations, err := h.InvitationRepositories.GetInvitationsByUserID(id)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while fetching invitations by user id.")
		return
	}

	var invitationResponses []invitation_dto.InvitationResponse
	for _, invitation := range invitations {
		invitationResponses = append(invitationResponses, ConvertToInvitationResponse(&invitation))
	}

	if len(invitations) == 0 {
		SuccessResponse(w, http.StatusOK, "No invitations found for the specified user.", invitationResponses)
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitations retrieved successfully", invitationResponses)
}

func (h *invitationHandlers) UpdateInvitation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request invitation_dto.UpdateInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid JSON format.")
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid invitation ID format. Please provide a numeric ID.")
		return
	}

	invitation, err := h.InvitationRepositories.GetInvitationByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No invitation found with the provided ID.")
		return
	}

	eventDate, err := time.Parse(time.RFC3339, request.InvitationData.EventDate)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid EventDate format. Use RFC3339.")
		return
	}

	invitation.Status = request.Status

	invitation.InvitationData.EventDate = eventDate

	uploadedFiles, ok := r.Context().Value(middleware.UploadsKey).(map[string]string)
	if ok {
		if val, exists := uploadedFiles["main_image_url"]; exists {
			invitation.InvitationData.MainImageURL = val
		}

		if invitation.InvitationData.Gallery != nil {
			if val, exists := uploadedFiles["image_url_1"]; exists {
				invitation.InvitationData.Gallery.ImageURL1 = val
			}
			if val, exists := uploadedFiles["image_url_2"]; exists {
				invitation.InvitationData.Gallery.ImageURL2 = val
			}
			if val, exists := uploadedFiles["image_url_3"]; exists {
				invitation.InvitationData.Gallery.ImageURL3 = val
			}
			if val, exists := uploadedFiles["image_url_4"]; exists {
				invitation.InvitationData.Gallery.ImageURL4 = val
			}
			if val, exists := uploadedFiles["image_url_5"]; exists {
				invitation.InvitationData.Gallery.ImageURL5 = val
			}
			if val, exists := uploadedFiles["image_url_6"]; exists {
				invitation.InvitationData.Gallery.ImageURL6 = val
			}
			if val, exists := uploadedFiles["image_url_7"]; exists {
				invitation.InvitationData.Gallery.ImageURL7 = val
			}
			if val, exists := uploadedFiles["image_url_8"]; exists {
				invitation.InvitationData.Gallery.ImageURL8 = val
			}
			if val, exists := uploadedFiles["image_url_9"]; exists {
				invitation.InvitationData.Gallery.ImageURL9 = val
			}
			if val, exists := uploadedFiles["image_url_10"]; exists {
				invitation.InvitationData.Gallery.ImageURL10 = val
			}
			if val, exists := uploadedFiles["image_url_11"]; exists {
				invitation.InvitationData.Gallery.ImageURL11 = val
			}
			if val, exists := uploadedFiles["image_url_12"]; exists {
				invitation.InvitationData.Gallery.ImageURL12 = val
			}
		}
	}

	err = h.InvitationRepositories.UpdateInvitation(invitation)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to update invitation.")
		return
	}

	SuccessResponse(w, http.StatusCreated, "Invitation updated successfully", ConvertToInvitationResponse(invitation))
}

func (h *invitationHandlers) DeleteInvitation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid invitation ID format. Please provide a numeric ID.")
		return
	}

	if _, err = h.InvitationRepositories.GetInvitationByID(uint(id)); err != nil {
		ErrorResponse(w, http.StatusNotFound, "No invitation found with the provided ID.")
		return
	}

	if err := h.InvitationRepositories.DeleteInvitation(uint(id)); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while deleting the invitation.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitation deleted successfully", nil)
}
