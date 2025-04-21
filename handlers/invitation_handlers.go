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
	InvitationRepositories      repositories.InvitationRepositories
	InvitationThemeRepositories repositories.InvitationThemeRepositories
}

func InvitationHandler(
	InvitationRepositories repositories.InvitationRepositories,
	InvitationThemeRepositories repositories.InvitationThemeRepositories,
) *invitationHandlers {
	return &invitationHandlers{InvitationRepositories, InvitationThemeRepositories}
}

func ConvertToInvitationResponse(invitation *models.Invitation) invitation_dto.InvitationResponse {
	invitationResponse := invitation_dto.InvitationResponse{
		ID:                  invitation.ID,
		Status:              invitation.Status,
		InvitationThemeID:   invitation.InvitationThemeID,
		InvitationThemeName: invitation.InvitationThemeName,
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
	invitationJSON := r.FormValue("invitation")
	var request invitation_dto.CreateInvitationRequest
	if err := json.Unmarshal([]byte(invitationJSON), &request); err != nil {
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

	invitationTheme, err := h.InvitationThemeRepositories.GetInvitationThemeByID(request.InvitationThemeID)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No invitation theme found with the provided ID.",
			"id": "Tema undangan tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	eventDate, err := time.Parse(time.RFC3339, request.InvitationData.EventDate)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid EventDate format. Use RFC3339.",
			"id": "EventDate format tidak valid. Gunakan RFC3339.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	invitation := &models.Invitation{
		Status:              models.InvitationStatusDraft,
		InvitationThemeID:   request.InvitationThemeID,
		InvitationThemeName: invitationTheme.Name,
		InvitationData: &models.InvitationData{
			EventName: request.InvitationData.EventName,
			EventDate: eventDate,
			Location:  request.InvitationData.Location,
			Gallery:   &models.Gallery{},
		},
		UserID: request.UserID,
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
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Failed to create invitation.",
			"id": "Gagal membuat undangan.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusCreated, "Invitation created successfully", ConvertToInvitationResponse(invitation))
}

func (h *invitationHandlers) GetInvitationByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid Invitation ID format. Please provide a numeric ID.",
			"id": "Format ID undangan tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	invitation, err := h.InvitationRepositories.GetInvitationByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No invitation found with the provided ID.",
			"id": "Undangan tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitation retrieved successfully", ConvertToInvitationResponse(invitation))
}

func (h *invitationHandlers) GetInvitations(w http.ResponseWriter, r *http.Request) {
	invitations, err := h.InvitationRepositories.GetInvitations()
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while fetching invitations.",
			"id": "Terjadi kesalahan saat mengambil undangan.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	if len(invitations) == 0 {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No invitations available at the moment.",
			"id": "Tidak ada undangan yang tersedia saat ini.",
		}
		SuccessResponse(w, http.StatusOK, messages[lang], []invitation_dto.InvitationResponse{})
		return
	}

	var invitationResponses []invitation_dto.InvitationResponse
	for _, invitation := range invitations {
		invitationResponses = append(invitationResponses, ConvertToInvitationResponse(&invitation))
	}

	SuccessResponse(w, http.StatusOK, "Invitations retrieved successfully", invitationResponses)
}

func (h *invitationHandlers) GetInvitationsByUserID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["userId"]

	invitations, err := h.InvitationRepositories.GetInvitationsByUserID(id)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while fetching invitations by user ID.",
			"id": "Terjadi kesalahan saat mengambil undangan berdasarkan ID user.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	if len(invitations) == 0 {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No invitations found for the specified user.",
			"id": "Tidak ditemukan undangan untuk pengguna yang dimaksud.",
		}
		SuccessResponse(w, http.StatusOK, messages[lang], []invitation_dto.InvitationResponse{})
		return
	}

	var invitationResponses []invitation_dto.InvitationResponse
	for _, invitation := range invitations {
		invitationResponses = append(invitationResponses, ConvertToInvitationResponse(&invitation))
	}

	SuccessResponse(w, http.StatusOK, "Invitations retrieved successfully", invitationResponses)
}

func (h *invitationHandlers) UpdateInvitationByID(w http.ResponseWriter, r *http.Request) {
	invitationJSON := r.FormValue("invitation")
	var request invitation_dto.UpdateInvitationRequest
	if err := json.Unmarshal([]byte(invitationJSON), &request); err != nil {
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
			"en": "Invalid Invitation ID format. Please provide a numeric ID.",
			"id": "Format ID undangan tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	invitation, err := h.InvitationRepositories.GetInvitationByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No invitation found with the provided ID.",
			"id": "Undangan tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if request.InvitationData.EventName != "" {
		invitation.InvitationData.EventName = request.InvitationData.EventName
	}
	eventDate, err := time.Parse(time.RFC3339, request.InvitationData.EventDate)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid EventDate format. Use RFC3339.",
			"id": "EventDate format tidak valid. Gunakan RFC3339.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}
	invitation.Status = request.Status
	invitation.InvitationData.EventDate = eventDate
	if request.InvitationData.Location != "" {
		invitation.InvitationData.Location = request.InvitationData.Location
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

	err = h.InvitationRepositories.UpdateInvitation(invitation)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Failed to update invitation",
			"id": "Gagal mengupdate undangan",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitation updated successfully", ConvertToInvitationResponse(invitation))
}

func (h *invitationHandlers) DeleteInvitationByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid Invitation ID format. Please provide a numeric ID.",
			"id": "Format ID undangan tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	if _, err = h.InvitationRepositories.GetInvitationByID(uint(id)); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No invitation found with the provided ID.",
			"id": "Undangan tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if err := h.InvitationRepositories.DeleteInvitation(uint(id)); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while deleting the invitation.",
			"id": "Terjadi kesalahan saat menghapus undangan.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Invitation deleted successfully", nil)
}
