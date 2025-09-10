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

	invitation := &models.Invitation{
		Status:              models.InvitationStatusDraft,
		InvitationThemeID:   request.InvitationThemeID,
		InvitationThemeName: invitationTheme.Name,
		InvitationData: &models.InvitationData{
			Bride: models.Bridegroom{
				Nickname:    request.InvitationData.Bride.Nickname,
				FullName:    request.InvitationData.Bride.FullName,
				Title:       request.InvitationData.Bride.Title,
				FatherName:  request.InvitationData.Bride.FatherName,
				FatherTitle: request.InvitationData.Bride.FatherTitle,
				MotherName:  request.InvitationData.Bride.MotherName,
				MotherTitle: request.InvitationData.Bride.MotherTitle,
			},
			Groom: models.Bridegroom{
				Nickname:    request.InvitationData.Groom.Nickname,
				FullName:    request.InvitationData.Groom.FullName,
				Title:       request.InvitationData.Groom.Title,
				FatherName:  request.InvitationData.Groom.FatherName,
				FatherTitle: request.InvitationData.Groom.FatherTitle,
				MotherName:  request.InvitationData.Groom.MotherName,
				MotherTitle: request.InvitationData.Groom.MotherTitle,
			},
			ContractEvent: models.Event{
				Place:   request.InvitationData.ContractEvent.Place,
				Address: request.InvitationData.ContractEvent.Address,
				MapsURL: request.InvitationData.ContractEvent.MapsURL,
			},
			ReceptionEvent: models.Event{
				Place:   request.InvitationData.ReceptionEvent.Place,
				Address: request.InvitationData.ReceptionEvent.Address,
				MapsURL: request.InvitationData.ReceptionEvent.MapsURL,
			},
			Gallery: &models.Gallery{},
		},
		UserID: request.UserID,
	}

	for _, bankAccount := range request.InvitationData.BankAccounts {
		invitation.InvitationData.BankAccounts = append(invitation.InvitationData.BankAccounts, models.BankAccount{
			BankName:    bankAccount.BankName,
			AccountName: bankAccount.AccountName,
			Number:      bankAccount.Number,
		})
	}

	contractStartTime, err := time.Parse(time.RFC3339, request.InvitationData.ContractEvent.StartTime)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid Contract event start time format. Use RFC3339.",
			"id": "Contract event start time format tidak valid. Gunakan RFC3339.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}
	invitation.InvitationData.ContractEvent.StartTime = contractStartTime

	if request.InvitationData.ContractEvent.EndTime != "" {
		contractEndTime, err := time.Parse(time.RFC3339, request.InvitationData.ContractEvent.EndTime)
		if err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "Invalid Contract event end time format. Use RFC3339.",
				"id": "Contract event end time format tidak valid. Gunakan RFC3339.",
			}
			ErrorResponse(w, http.StatusBadRequest, messages, lang)
			return
		}
		invitation.InvitationData.ContractEvent.EndTime = &contractEndTime
	}

	receptionStartTime, err := time.Parse(time.RFC3339, request.InvitationData.ReceptionEvent.StartTime)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid Reception event start time format. Use RFC3339.",
			"id": "Reception event start time format tidak valid. Gunakan RFC3339.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}
	invitation.InvitationData.ReceptionEvent.StartTime = receptionStartTime

	if request.InvitationData.ReceptionEvent.EndTime != "" {
		receptionEndTime, err := time.Parse(time.RFC3339, request.InvitationData.ReceptionEvent.EndTime)
		if err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "Invalid Reception event start time format. Use RFC3339.",
				"id": "Reception event start time format tidak valid. Gunakan RFC3339.",
			}
			ErrorResponse(w, http.StatusBadRequest, messages, lang)
			return
		}
		invitation.InvitationData.ReceptionEvent.EndTime = &receptionEndTime
	}

	uploadedFiles, ok := r.Context().Value(middleware.UploadsKey).(map[string]string)
	if ok {
		if val, exists := uploadedFiles["cover_image"]; exists {
			invitation.InvitationData.CoverImageURL = val
		}

		if val, exists := uploadedFiles["bride_image"]; exists {
			invitation.InvitationData.Bride.ImageURL = val
		}
		if val, exists := uploadedFiles["groom_image"]; exists {
			invitation.InvitationData.Groom.ImageURL = val
		}

		if invitation.InvitationData.Gallery != nil {
			if val, exists := uploadedFiles["image_1"]; exists {
				invitation.InvitationData.Gallery.ImageURL1 = val
			}
			if val, exists := uploadedFiles["image_2"]; exists {
				invitation.InvitationData.Gallery.ImageURL2 = val
			}
			if val, exists := uploadedFiles["image_3"]; exists {
				invitation.InvitationData.Gallery.ImageURL3 = val
			}
			if val, exists := uploadedFiles["image_4"]; exists {
				invitation.InvitationData.Gallery.ImageURL4 = val
			}
			if val, exists := uploadedFiles["image_5"]; exists {
				invitation.InvitationData.Gallery.ImageURL5 = val
			}
			if val, exists := uploadedFiles["image_6"]; exists {
				invitation.InvitationData.Gallery.ImageURL6 = val
			}
			if val, exists := uploadedFiles["image_7"]; exists {
				invitation.InvitationData.Gallery.ImageURL7 = val
			}
			if val, exists := uploadedFiles["image_8"]; exists {
				invitation.InvitationData.Gallery.ImageURL8 = val
			}
			if val, exists := uploadedFiles["image_9"]; exists {
				invitation.InvitationData.Gallery.ImageURL9 = val
			}
			if val, exists := uploadedFiles["image_10"]; exists {
				invitation.InvitationData.Gallery.ImageURL10 = val
			}
			if val, exists := uploadedFiles["image_11"]; exists {
				invitation.InvitationData.Gallery.ImageURL11 = val
			}
			if val, exists := uploadedFiles["image_12"]; exists {
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

	invitation.Status = request.Status

	if request.InvitationData.Bride.Nickname != "" {
		invitation.InvitationData.Bride.Nickname = request.InvitationData.Bride.Nickname
	}
	if request.InvitationData.Bride.FullName != "" {
		invitation.InvitationData.Bride.FullName = request.InvitationData.Bride.FullName
	}
	if request.InvitationData.Bride.Title != "" {
		invitation.InvitationData.Bride.Title = request.InvitationData.Bride.Title
	}
	if request.InvitationData.Bride.FatherName != "" {
		invitation.InvitationData.Bride.FatherName = request.InvitationData.Bride.FatherName
	}
	if request.InvitationData.Bride.FatherTitle != "" {
		invitation.InvitationData.Bride.FatherTitle = request.InvitationData.Bride.FatherTitle
	}
	if request.InvitationData.Bride.MotherName != "" {
		invitation.InvitationData.Bride.MotherName = request.InvitationData.Bride.MotherName
	}
	if request.InvitationData.Bride.MotherTitle != "" {
		invitation.InvitationData.Bride.MotherTitle = request.InvitationData.Bride.MotherTitle
	}

	if request.InvitationData.Groom.Nickname != "" {
		invitation.InvitationData.Groom.Nickname = request.InvitationData.Groom.Nickname
	}
	if request.InvitationData.Groom.FullName != "" {
		invitation.InvitationData.Groom.FullName = request.InvitationData.Groom.FullName
	}
	if request.InvitationData.Groom.Title != "" {
		invitation.InvitationData.Groom.Title = request.InvitationData.Groom.Title
	}
	if request.InvitationData.Groom.FatherName != "" {
		invitation.InvitationData.Groom.FatherName = request.InvitationData.Groom.FatherName
	}
	if request.InvitationData.Groom.FatherTitle != "" {
		invitation.InvitationData.Groom.FatherTitle = request.InvitationData.Groom.FatherTitle
	}
	if request.InvitationData.Groom.MotherName != "" {
		invitation.InvitationData.Groom.MotherName = request.InvitationData.Groom.MotherName
	}
	if request.InvitationData.Groom.FatherTitle != "" {
		invitation.InvitationData.Groom.FatherTitle = request.InvitationData.Groom.FatherTitle
	}

	if request.InvitationData.ContractEvent.Place != "" {
		invitation.InvitationData.ContractEvent.Place = request.InvitationData.ContractEvent.Place
	}
	if request.InvitationData.ContractEvent.Address != "" {
		invitation.InvitationData.ContractEvent.Address = request.InvitationData.ContractEvent.Address
	}
	if request.InvitationData.ContractEvent.MapsURL != "" {
		invitation.InvitationData.ContractEvent.MapsURL = request.InvitationData.ContractEvent.MapsURL
	}

	if request.InvitationData.ReceptionEvent.Place != "" {
		invitation.InvitationData.ReceptionEvent.Place = request.InvitationData.ReceptionEvent.Place
	}
	if request.InvitationData.ReceptionEvent.Address != "" {
		invitation.InvitationData.ReceptionEvent.Address = request.InvitationData.ReceptionEvent.Address
	}
	if request.InvitationData.ReceptionEvent.MapsURL != "" {
		invitation.InvitationData.ReceptionEvent.MapsURL = request.InvitationData.ReceptionEvent.MapsURL
	}

	if request.InvitationData.ContractEvent.StartTime != "" {
		contractStartTime, err := time.Parse(time.RFC3339, request.InvitationData.ContractEvent.StartTime)
		if err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "Invalid Contract event start time format. Use RFC3339.",
				"id": "Contract event start time format tidak valid. Gunakan RFC3339.",
			}
			ErrorResponse(w, http.StatusBadRequest, messages, lang)
			return
		}
		invitation.InvitationData.ContractEvent.StartTime = contractStartTime
	}

	if request.InvitationData.ContractEvent.EndTime != "" {
		contractEndTime, err := time.Parse(time.RFC3339, request.InvitationData.ContractEvent.EndTime)
		if err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "Invalid Contract event end time format. Use RFC3339.",
				"id": "Contract event end time format tidak valid. Gunakan RFC3339.",
			}
			ErrorResponse(w, http.StatusBadRequest, messages, lang)
			return
		}
		invitation.InvitationData.ContractEvent.EndTime = &contractEndTime
	}

	if request.InvitationData.ReceptionEvent.StartTime != "" {
		receptionStartTime, err := time.Parse(time.RFC3339, request.InvitationData.ReceptionEvent.StartTime)
		if err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "Invalid Reception event start time format. Use RFC3339.",
				"id": "Reception event start time format tidak valid. Gunakan RFC3339.",
			}
			ErrorResponse(w, http.StatusBadRequest, messages, lang)
			return
		}
		invitation.InvitationData.ReceptionEvent.StartTime = receptionStartTime
	}

	if request.InvitationData.ReceptionEvent.EndTime != "" {
		receptionEndTime, err := time.Parse(time.RFC3339, request.InvitationData.ReceptionEvent.EndTime)
		if err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "Invalid Reception event start time format. Use RFC3339.",
				"id": "Reception event start time format tidak valid. Gunakan RFC3339.",
			}
			ErrorResponse(w, http.StatusBadRequest, messages, lang)
			return
		}
		invitation.InvitationData.ReceptionEvent.EndTime = &receptionEndTime
	}

	uploadedFiles, ok := r.Context().Value(middleware.UploadsKey).(map[string]string)
	if ok {
		if val, exists := uploadedFiles["cover_image"]; exists {
			invitation.InvitationData.CoverImageURL = val
		}

		if val, exists := uploadedFiles["bride_image"]; exists {
			invitation.InvitationData.Bride.ImageURL = val
		}
		if val, exists := uploadedFiles["groom_image"]; exists {
			invitation.InvitationData.Groom.ImageURL = val
		}

		if invitation.InvitationData.Gallery != nil {
			if val, exists := uploadedFiles["image_1"]; exists {
				invitation.InvitationData.Gallery.ImageURL1 = val
			}
			if val, exists := uploadedFiles["image_2"]; exists {
				invitation.InvitationData.Gallery.ImageURL2 = val
			}
			if val, exists := uploadedFiles["image_3"]; exists {
				invitation.InvitationData.Gallery.ImageURL3 = val
			}
			if val, exists := uploadedFiles["image_4"]; exists {
				invitation.InvitationData.Gallery.ImageURL4 = val
			}
			if val, exists := uploadedFiles["image_5"]; exists {
				invitation.InvitationData.Gallery.ImageURL5 = val
			}
			if val, exists := uploadedFiles["image_6"]; exists {
				invitation.InvitationData.Gallery.ImageURL6 = val
			}
			if val, exists := uploadedFiles["image_7"]; exists {
				invitation.InvitationData.Gallery.ImageURL7 = val
			}
			if val, exists := uploadedFiles["image_8"]; exists {
				invitation.InvitationData.Gallery.ImageURL8 = val
			}
			if val, exists := uploadedFiles["image_9"]; exists {
				invitation.InvitationData.Gallery.ImageURL9 = val
			}
			if val, exists := uploadedFiles["image_10"]; exists {
				invitation.InvitationData.Gallery.ImageURL10 = val
			}
			if val, exists := uploadedFiles["image_11"]; exists {
				invitation.InvitationData.Gallery.ImageURL11 = val
			}
			if val, exists := uploadedFiles["image_12"]; exists {
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
