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
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	eventDate, err := time.Parse(time.RFC3339, request.InvitationData.EventDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	invitation := models.Invitation{
		UserID: request.UserID,
		Status: request.Status,
		InvitationData: &models.InvitationData{
			EventName: request.InvitationData.EventName,
			EventDate: eventDate,
			Location:  request.InvitationData.Location,
			Gallery:   []*models.Gallery{}, // Inisialisasi agar tidak nil
		},
	}

	// If there are uploaded files, update the Gallery
	for _, gallery := range request.InvitationData.Gallery {
		invitation.InvitationData.Gallery = append(invitation.InvitationData.Gallery, &models.Gallery{
			Position: gallery.Position,
			ImageURL: gallery.ImageURL,
		})
	}

	// Start a database transaction to ensure data consistency
	err = h.InvitationRepositories.CreateInvitation(invitation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Create Success")
}

func (h *invitationHandlers) GetInvitationByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	invitation, err := h.InvitationRepositories.GetInvitationByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: invitation}
	json.NewEncoder(w).Encode(response)
}

func (h *invitationHandlers) GetInvitations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	invitations, err := h.InvitationRepositories.GetInvitations()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: invitations}
	json.NewEncoder(w).Encode(response)
}

func (h *invitationHandlers) GetInvitationsByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["userId"]
	invitations, err := h.InvitationRepositories.GetInvitationsByUserID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: invitations}
	json.NewEncoder(w).Encode(response)
}

func (h *invitationHandlers) UpdateInvitation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode request body
	request := new(invitation_dto.InvitationRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// Parse event date
	eventDate, err := time.Parse(time.RFC3339, request.InvitationData.EventDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	invitation, err := h.InvitationRepositories.GetInvitationByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// Update fields of the invitation
	invitation.Status = request.Status
	invitation.InvitationData.EventName = request.InvitationData.EventName
	invitation.InvitationData.EventDate = eventDate
	invitation.InvitationData.Location = request.InvitationData.Location

	// If there are uploaded files, update the Gallery
	for _, gallery := range request.InvitationData.Gallery {
		invitation.InvitationData.Gallery = append(invitation.InvitationData.Gallery, &models.Gallery{
			Position: gallery.Position,
			ImageURL: gallery.ImageURL,
		})
	}

	// Start a database transaction to ensure data consistency
	err = h.InvitationRepositories.UpdateInvitation(invitation)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Update Success")
}

func (h *invitationHandlers) DeleteInvitation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	invitation, err := h.InvitationRepositories.GetInvitationByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	err = h.InvitationRepositories.DeleteInvitation(invitation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Delete Success")
}
