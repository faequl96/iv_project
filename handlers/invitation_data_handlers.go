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
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid invitation data ID format. Please provide a numeric ID.",
			"id": "Format ID data undangan tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	invitationData, err := h.InvitationDataRepositories.GetInvitationDataByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No invitation data found with the provided ID.",
			"id": "Data undangan tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitation data retrieved successfully", ConvertToInvitationDataResponse(invitationData))
}

func (h *invitationDataHandlers) UpdateInvitationDataByID(w http.ResponseWriter, r *http.Request) {
	invitationDataJSON := r.FormValue("invitation_data")
	var request invitation_data_dto.UpdateInvitationDataRequest
	if err := json.Unmarshal([]byte(invitationDataJSON), &request); err != nil {
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
			"en": "Invalid invitation data ID format. Please provide a numeric ID.",
			"id": "Format ID data undangan tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	invitationData, err := h.InvitationDataRepositories.GetInvitationDataByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No invitation data found with the provided ID.",
			"id": "Data undangan tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if request.EventName != "" {
		invitationData.EventName = request.EventName
	}
	eventDate, err := time.Parse(time.RFC3339, request.EventDate)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid EventDate format. Use RFC3339.",
			"id": "EventDate format tidak valid. Gunakan RFC3339.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
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
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Failed to update invitation data",
			"id": "Gagal mengupdate data undangan",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitation data updated successfully", ConvertToInvitationDataResponse(invitationData))
}
