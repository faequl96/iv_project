package handlers

import (
	"encoding/json"
	gallery_dto "iv_project/dto/gallery"
	invitation_dto "iv_project/dto/invitation"
	invitation_data_dto "iv_project/dto/invitation_data"
	invitation_theme_dto "iv_project/dto/invitation_theme"
	user_dto "iv_project/dto/user"
	user_profile_dto "iv_project/dto/user_profile"
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
	return invitation_dto.InvitationResponse{
		ID: invitation.ID,
		User: &user_dto.UserResponse{
			ID: invitation.User.ID,
			UserProfile: &user_profile_dto.UserProfileResponse{
				ID:        invitation.User.UserProfile.ID,
				FirstName: invitation.User.UserProfile.FirstName,
				LastName:  invitation.User.UserProfile.LastName,
			},
		},
		InvitationTheme: &invitation_theme_dto.InvitationThemeResponse{
			ID:               invitation.InvitationTheme.ID,
			Title:            invitation.InvitationTheme.Title,
			IDRPrice:         invitation.InvitationTheme.IDRPrice,
			IDRDiscountPrice: invitation.InvitationTheme.IDRDiscountPrice,
			IVCPrice:         invitation.InvitationTheme.IVCPrice,
			IVCDiscountPrice: invitation.InvitationTheme.IVCDiscountPrice,
		},
		Status: invitation.Status,
		InvitationData: &invitation_data_dto.InvitationDataResponse{
			ID:           invitation.InvitationData.ID,
			EventName:    invitation.InvitationData.EventName,
			EventDate:    invitation.InvitationData.EventDate.Format(time.RFC3339),
			Location:     invitation.InvitationData.Location,
			MainImageURL: invitation.InvitationData.MainImageURL,
			Gallery: &gallery_dto.GalleryResponse{
				ID:         invitation.InvitationData.Gallery.ID,
				ImageURL1:  invitation.InvitationData.Gallery.ImageURL1,
				ImageURL2:  invitation.InvitationData.Gallery.ImageURL2,
				ImageURL3:  invitation.InvitationData.Gallery.ImageURL3,
				ImageURL4:  invitation.InvitationData.Gallery.ImageURL4,
				ImageURL5:  invitation.InvitationData.Gallery.ImageURL5,
				ImageURL6:  invitation.InvitationData.Gallery.ImageURL6,
				ImageURL7:  invitation.InvitationData.Gallery.ImageURL7,
				ImageURL8:  invitation.InvitationData.Gallery.ImageURL8,
				ImageURL9:  invitation.InvitationData.Gallery.ImageURL9,
				ImageURL10: invitation.InvitationData.Gallery.ImageURL10,
				ImageURL11: invitation.InvitationData.Gallery.ImageURL11,
				ImageURL12: invitation.InvitationData.Gallery.ImageURL12,
			},
		},
	}
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

	uploadedFiles, ok := r.Context().Value(middleware.UploadedFilesKey).(map[string]string)
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

	request := new(invitation_dto.UpdateInvitationRequest)
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

	uploadedFiles, ok := r.Context().Value(middleware.UploadedFilesKey).(map[string]string)
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
