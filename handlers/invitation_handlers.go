package handlers

import (
	"encoding/json"
	"iv_project/dto"
	invitation_dto "iv_project/dto/invitation"
	"iv_project/models"
	"iv_project/pkg/middleware"
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

	uploadedFiles, ok := r.Context().Value(middleware.UploadedFilesKey).(map[string]string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "Tidak ada file yang diunggah"}
		json.NewEncoder(w).Encode(response)
		return
	}

	invitation := models.Invitation{
		UserID: request.UserID,
		Status: request.Status,
		InvitationData: &models.InvitationData{
			EventName: request.InvitationData.EventName,
			EventDate: eventDate,
			Location:  request.InvitationData.Location,
			Gallery: &models.Gallery{
				ImageURL1:  uploadedFiles["image_url_1"],
				ImageURL2:  uploadedFiles["image_url_2"],
				ImageURL3:  uploadedFiles["image_url_3"],
				ImageURL4:  uploadedFiles["image_url_4"],
				ImageURL5:  uploadedFiles["image_url_5"],
				ImageURL6:  uploadedFiles["image_url_6"],
				ImageURL7:  uploadedFiles["image_url_7"],
				ImageURL8:  uploadedFiles["image_url_8"],
				ImageURL9:  uploadedFiles["image_url_9"],
				ImageURL10: uploadedFiles["image_url_10"],
				ImageURL11: uploadedFiles["image_url_11"],
				ImageURL12: uploadedFiles["image_url_12"],
			},
		},
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

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	invitation, err := h.InvitationRepositories.GetInvitationByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
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

	uploadedFiles, ok := r.Context().Value(middleware.UploadedFilesKey).(map[string]string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "Tidak ada file yang diunggah"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Update fields of the invitation
	invitation.Status = request.Status
	invitation.InvitationData.EventName = request.InvitationData.EventName
	invitation.InvitationData.EventDate = eventDate
	invitation.InvitationData.Location = request.InvitationData.Location
	invitation.InvitationData.Gallery.ImageURL1 = uploadedFiles["image_url_1"]
	invitation.InvitationData.Gallery.ImageURL2 = uploadedFiles["image_url_2"]
	invitation.InvitationData.Gallery.ImageURL3 = uploadedFiles["image_url_3"]
	invitation.InvitationData.Gallery.ImageURL4 = uploadedFiles["image_url_4"]
	invitation.InvitationData.Gallery.ImageURL5 = uploadedFiles["image_url_5"]
	invitation.InvitationData.Gallery.ImageURL6 = uploadedFiles["image_url_6"]
	invitation.InvitationData.Gallery.ImageURL7 = uploadedFiles["image_url_7"]
	invitation.InvitationData.Gallery.ImageURL8 = uploadedFiles["image_url_8"]
	invitation.InvitationData.Gallery.ImageURL9 = uploadedFiles["image_url_9"]
	invitation.InvitationData.Gallery.ImageURL10 = uploadedFiles["image_url_10"]
	invitation.InvitationData.Gallery.ImageURL11 = uploadedFiles["image_url_11"]
	invitation.InvitationData.Gallery.ImageURL12 = uploadedFiles["image_url_12"]

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
