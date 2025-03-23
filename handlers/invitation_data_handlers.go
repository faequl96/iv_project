package handlers

import (
	"encoding/json"
	invitation_data_dto "iv_project/dto/invitation_data"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type invitationDataHandlers struct {
	InvitationDataRepositories repositories.InvitationDataRepositories
}

func InvitationDataHandler(InvitationDataRepositories repositories.InvitationDataRepositories) *invitationDataHandlers {
	return &invitationDataHandlers{InvitationDataRepositories}
}

func ConvertToInvitationDataResponse(invitationData *models.InvitationData) invitation_data_dto.InvitationDataResponse {
	invitationDataResponse := invitation_data_dto.InvitationDataResponse{
		ID:           invitationData.ID,
		EventName:    invitationData.EventName,
		EventDate:    invitationData.EventDate.Format(time.RFC3339),
		Location:     invitationData.Location,
		MainImageURL: invitationData.MainImageURL,
	}

	if invitationData.Gallery != nil {
		galleryResponse := ConvertToGalleryResponse(invitationData.Gallery)
		invitationDataResponse.Gallery = &galleryResponse
	}

	return invitationDataResponse
}

func (h *invitationDataHandlers) GetInvitationDataByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid invitation data ID format. Please provide a numeric ID.")
		return
	}

	invitationData, err := h.InvitationDataRepositories.GetInvitationDataByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No invitation data found with the provided ID.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitation data retrieved successfully", ConvertToInvitationDataResponse(invitationData))
}

func (h *invitationDataHandlers) UpdateInvitationDataByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	invitationDataJSON := r.FormValue("invitation_data")
	var request invitation_data_dto.UpdateInvitationDataRequest
	if err := json.Unmarshal([]byte(invitationDataJSON), &request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid JSON format.")
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid invitation data ID format. Please provide a numeric ID.")
		return
	}

	invitationData, err := h.InvitationDataRepositories.GetInvitationDataByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No invitation data found with the provided ID.")
		return
	}

	if request.EventName != "" {
		invitationData.EventName = request.EventName
	}
	eventDate, err := time.Parse(time.RFC3339, request.EventDate)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid EventDate format. Use RFC3339.")
		return
	}
	invitationData.EventDate = eventDate
	if request.Location != "" {
		invitationData.Location = request.Location
	}

	uploadedFiles, ok := r.Context().Value(middleware.UploadsKey).(map[string]string)
	if ok {
		if val, exists := uploadedFiles["main_image_url"]; exists {
			invitationData.MainImageURL = val
		}

		if invitationData.Gallery != nil {
			if val, exists := uploadedFiles["image_url_1"]; exists {
				invitationData.Gallery.ImageURL1 = val
			}
			if val, exists := uploadedFiles["image_url_2"]; exists {
				invitationData.Gallery.ImageURL2 = val
			}
			if val, exists := uploadedFiles["image_url_3"]; exists {
				invitationData.Gallery.ImageURL3 = val
			}
			if val, exists := uploadedFiles["image_url_4"]; exists {
				invitationData.Gallery.ImageURL4 = val
			}
			if val, exists := uploadedFiles["image_url_5"]; exists {
				invitationData.Gallery.ImageURL5 = val
			}
			if val, exists := uploadedFiles["image_url_6"]; exists {
				invitationData.Gallery.ImageURL6 = val
			}
			if val, exists := uploadedFiles["image_url_7"]; exists {
				invitationData.Gallery.ImageURL7 = val
			}
			if val, exists := uploadedFiles["image_url_8"]; exists {
				invitationData.Gallery.ImageURL8 = val
			}
			if val, exists := uploadedFiles["image_url_9"]; exists {
				invitationData.Gallery.ImageURL9 = val
			}
			if val, exists := uploadedFiles["image_url_10"]; exists {
				invitationData.Gallery.ImageURL10 = val
			}
			if val, exists := uploadedFiles["image_url_11"]; exists {
				invitationData.Gallery.ImageURL11 = val
			}
			if val, exists := uploadedFiles["image_url_12"]; exists {
				invitationData.Gallery.ImageURL12 = val
			}
		}
	}

	err = h.InvitationDataRepositories.UpdateInvitationData(invitationData)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to update invitation data.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitation data updated successfully", ConvertToInvitationDataResponse(invitationData))
}
