package handlers

import (
	"encoding/json"
	rsvp_dto "iv_project/dto/rsvp"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type rsvpHandlers struct {
	RSVPRepositories repositories.RSVPRepositories
}

func RSVPHandlers(
	RSVPRepositories repositories.RSVPRepositories,
) *rsvpHandlers {
	return &rsvpHandlers{RSVPRepositories}
}

func ConvertToRSVPResponse(rsvp *models.RSVP) rsvp_dto.RSVPResponse {
	rsvpResponse := rsvp_dto.RSVPResponse{
		ID:             rsvp.ID,
		InvitationID:   rsvp.InvitationID,
		InvitedGuestID: rsvp.InvitedGuestID,
		Nickname:       rsvp.Nickname,
		Avatar:         rsvp.Avatar,
		Invited:        rsvp.Invited,
		Attendance:     rsvp.Attendance,
		Message:        rsvp.Message,
	}

	return rsvpResponse
}

func (h *rsvpHandlers) CreateRSVP(w http.ResponseWriter, r *http.Request) {
	var request rsvp_dto.RSVPRequest
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

	rsvp := &models.RSVP{
		InvitationID:   request.InvitationID,
		InvitedGuestID: request.InvitedGuestID,
		Nickname:       request.Nickname,
		Avatar:         request.Avatar,
		Invited:        request.Invited,
		Attendance:     request.Attendance,
		Message:        request.Message,
	}

	if err := h.RSVPRepositories.CreateRSVP(rsvp); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An Error occurred while creating the rsvp. Please try again later.",
			"id": "Terjadi kesalahan saat membuat rsvp. Coba lagi nanti.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusCreated, "RSVP successfully created", ConvertToRSVPResponse(rsvp))
}

func (h *rsvpHandlers) GetRSVPsByInvitationID(w http.ResponseWriter, r *http.Request) {
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

	rsvps, err := h.RSVPRepositories.GetRSVPsByInvitationID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while retrieving rsvps. Please try again.",
			"id": "Terjadi kesalahan saat mengambil rsvp. Silahkan coba lagi nanti",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	if len(rsvps) == 0 {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No rsvps found for the specified invitation.",
			"id": "Tidak ditemukan rsvp berdasarkan undangan yang dimaksud.",
		}
		SuccessResponse(w, http.StatusOK, messages[lang], []rsvp_dto.RSVPResponse{})
		return
	}

	var response []rsvp_dto.RSVPResponse
	for _, rsvp := range rsvps {
		response = append(response, ConvertToRSVPResponse(&rsvp))
	}

	SuccessResponse(w, http.StatusOK, "RSVPs retrieved successfully", response)
}
