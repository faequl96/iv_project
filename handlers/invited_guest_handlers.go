package handlers

import (
	"encoding/json"
	invited_guest_dto "iv_project/dto/invited_guest"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type invitedGuestHandlers struct {
	InvitedGuestRepositories repositories.InvitedGuestRepositories
}

func InvitedGuestHandlers(
	InvitedGuestRepositories repositories.InvitedGuestRepositories,
) *invitedGuestHandlers {
	return &invitedGuestHandlers{InvitedGuestRepositories}
}

func ConvertToInvitedGuestResponse(invitedGuest *models.InvitedGuest) invited_guest_dto.InvitedGuestResponse {
	invitedGuestResponse := invited_guest_dto.InvitedGuestResponse{
		ID:           invitedGuest.ID,
		InvitationID: invitedGuest.InvitationID,
		NameInstance: invitedGuest.NameInstance,
		Name:         invitedGuest.Name,
		Instance:     invitedGuest.Instance,
		Nickname:     invitedGuest.Nickname,
		Avatar:       invitedGuest.Avatar,
		Attendance:   invitedGuest.Attendance,
	}

	return invitedGuestResponse
}

func (h *invitedGuestHandlers) CreateInvitedGuest(w http.ResponseWriter, r *http.Request) {
	var request invited_guest_dto.CreateInvitedGuestRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
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

	invitedGuest := &models.InvitedGuest{
		InvitationID: request.InvitationID,
		NameInstance: request.NameInstance,
		Name:         request.Name,
		Instance:     request.Instance,
		Nickname:     request.Nickname,
		Avatar:       request.Avatar,
		Attendance:   request.Attendance,
	}

	if err := h.InvitedGuestRepositories.CreateInvitedGuest(invitedGuest); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An Error occurred while creating the invited guest. Please try again later.",
			"id": "Terjadi kesalahan saat membuat tamu undangan. Coba lagi nanti.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusCreated, "Invited Guest successfully created", ConvertToInvitedGuestResponse(invitedGuest))
}

func (h *invitedGuestHandlers) GetInvitedGuestsByInvitationID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["invitationId"])
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid invitation ID format. Please provide a numeric ID.",
			"id": "Format ID undangan tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	invitedGuests, err := h.InvitedGuestRepositories.GetInvitedGuestsByInvitationID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while retrieving invited guests. Please try again.",
			"id": "Terjadi kesalahan saat mengambil tamu undangan. Silahkan coba lagi nanti",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	if len(invitedGuests) == 0 {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No invited guests found for the specified invitation.",
			"id": "Tidak ditemukan tamu undangan berdasarkan undangan yang dimaksud.",
		}
		SuccessResponse(w, http.StatusOK, messages[lang], []invited_guest_dto.InvitedGuestResponse{})
		return
	}

	var response []invited_guest_dto.InvitedGuestResponse
	for _, invitedGuest := range invitedGuests {
		response = append(response, ConvertToInvitedGuestResponse(&invitedGuest))
	}

	SuccessResponse(w, http.StatusOK, "Invited Guests retrieved successfully", response)
}

func (h *invitedGuestHandlers) UpdateInvitedGuestByID(w http.ResponseWriter, r *http.Request) {
	var request invited_guest_dto.UpdateInvitedGuestRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
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
			"en": "Invalid invited guest ID format. Please provide a numeric ID.",
			"id": "Format ID tamu undangan tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	invitedGuest, err := h.InvitedGuestRepositories.GetInvitedGuestByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No invited guest found with the provided ID.",
			"id": "Tamu undangan tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if request.Nickname != "" {
		invitedGuest.Nickname = request.Nickname
	}
	if request.Avatar != "" {
		invitedGuest.Avatar = request.Avatar
	}
	if request.Attendance != "" {
		invitedGuest.Attendance = request.Attendance
	}

	if err := h.InvitedGuestRepositories.UpdateInvitedGuest(invitedGuest); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while updating the invited guest.",
			"id": "Terjadi kesalahan saat mengupdate tamu undangan.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Invited Guest successfully updated", ConvertToInvitedGuestResponse(invitedGuest))
}
