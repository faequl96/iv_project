package handlers

import (
	"encoding/json"
	"iv_project/dto"
	invitation_dto "iv_project/dto/invitation"
	"iv_project/models"
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

func (h *invitationHandlers) CreateInvitation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(invitation_dto.InvitationRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	eventDate, err := time.Parse(time.RFC3339, request.InvitationData.EventDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "Invalid event_date format. Use RFC3339 (e.g. 2025-06-15T14:00:00Z)"}
		json.NewEncoder(w).Encode(response)
		return
	}

	invitation := models.Invitation{
		UserID: request.UserID,
		Status: request.Status,
		InvitationData: models.InvitationData{
			EventName: request.InvitationData.EventName,
			EventDate: eventDate,
			Location:  request.InvitationData.Location,
		},
	}

	err = h.InvitationRepositories.CreateInvitation(invitation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: "Success"}
	json.NewEncoder(w).Encode(response)
}

func (h *invitationHandlers) GetInvitationByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	invitation, err := h.InvitationRepositories.GetInvitationByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: invitation}
	json.NewEncoder(w).Encode(response)
}

func (h *invitationHandlers) GetInvitationsByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	invitations, err := h.InvitationRepositories.GetInvitationsByUserID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: invitations}
	json.NewEncoder(w).Encode(response)
}
