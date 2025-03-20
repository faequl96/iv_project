package handlers

import (
	"encoding/json"
	gallery_dto "iv_project/dto/gallery"
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
	return invitation_data_dto.InvitationDataResponse{
		ID:           invitationData.ID,
		EventName:    invitationData.EventName,
		EventDate:    invitationData.EventDate.Format(time.RFC3339),
		Location:     invitationData.Location,
		MainImageURL: invitationData.MainImageURL,
		Gallery: &gallery_dto.GalleryResponse{
			ID:         invitationData.Gallery.ID,
			ImageURL1:  invitationData.Gallery.ImageURL1,
			ImageURL2:  invitationData.Gallery.ImageURL2,
			ImageURL3:  invitationData.Gallery.ImageURL3,
			ImageURL4:  invitationData.Gallery.ImageURL4,
			ImageURL5:  invitationData.Gallery.ImageURL5,
			ImageURL6:  invitationData.Gallery.ImageURL6,
			ImageURL7:  invitationData.Gallery.ImageURL7,
			ImageURL8:  invitationData.Gallery.ImageURL8,
			ImageURL9:  invitationData.Gallery.ImageURL9,
			ImageURL10: invitationData.Gallery.ImageURL10,
			ImageURL11: invitationData.Gallery.ImageURL11,
			ImageURL12: invitationData.Gallery.ImageURL12,
		},
	}
}

func (h *invitationDataHandlers) CreateInvitationData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request invitation_data_dto.CreateInvitationDataRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid JSON format.")
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	eventDate, err := time.Parse(time.RFC3339, request.EventDate)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid EventDate format. Use RFC3339.")
		return
	}

	invitationData := &models.InvitationData{
		EventName: request.EventName,
		EventDate: eventDate,
		Location:  request.Location,
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

	err = h.InvitationDataRepositories.CreateInvitationData(invitationData)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to create invitation data.")
		return
	}

	SuccessResponse(w, http.StatusCreated, "Invitation data created successfully", ConvertToInvitationDataResponse(invitationData))
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

func (h *invitationDataHandlers) UpdateInvitationData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request invitation_data_dto.UpdateInvitationDataRequest
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
		ErrorResponse(w, http.StatusBadRequest, "Invalid invitation data ID format. Please provide a numeric ID.")
		return
	}

	invitationData, err := h.InvitationDataRepositories.GetInvitationDataByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No invitation data found with the provided ID.")
		return
	}

	eventDate, err := time.Parse(time.RFC3339, request.EventDate)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid EventDate format. Use RFC3339.")
		return
	}

	invitationData.EventDate = eventDate

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

	SuccessResponse(w, http.StatusCreated, "Invitation data updated successfully", ConvertToInvitationDataResponse(invitationData))
}

func (h *invitationDataHandlers) DeleteInvitationData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid invitation data ID format. Please provide a numeric ID.")
		return
	}

	if _, err = h.InvitationDataRepositories.GetInvitationDataByID(uint(id)); err != nil {
		ErrorResponse(w, http.StatusNotFound, "No invitation data found with the provided ID.")
		return
	}

	if err := h.InvitationDataRepositories.DeleteInvitationData(uint(id)); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while deleting the invitation data.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitation data deleted successfully", nil)
}
