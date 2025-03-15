package handlers

import (
	"encoding/json"
	"iv_project/dto"
	invitation_dto "iv_project/dto/invitation"
	"iv_project/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type invitationDataHandlers struct {
	InvitationDataRepositories repositories.InvitationDataRepositories
}

func InvitationDataHandler(InvitationDataRepositories repositories.InvitationDataRepositories) *invitationDataHandlers {
	return &invitationDataHandlers{InvitationDataRepositories}
}

func (h *invitationDataHandlers) UpdateInvitationData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(invitation_dto.InvitationDataRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	invitationData, err := h.InvitationDataRepositories.GetInvitationDataByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.EventName != "" {
		invitationData.EventName = request.EventName
	}
	if request.EventDate != "" {
		layout := "2006-01-02 15:04:05"
		newTime, err := time.Parse(layout, request.EventDate)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
		invitationData.EventDate = newTime
	}
	if request.Location != "" {
		invitationData.Location = request.Location
	}
	if request.GalleryImageURL1 != "" {
		invitationData.GalleryImageURL1 = request.GalleryImageURL1
	}
	if request.GalleryImageURL2 != "" {
		invitationData.GalleryImageURL2 = request.GalleryImageURL2
	}
	if request.GalleryImageURL3 != "" {
		invitationData.GalleryImageURL3 = request.GalleryImageURL3
	}
	if request.GalleryImageURL4 != "" {
		invitationData.GalleryImageURL4 = request.GalleryImageURL4
	}
	if request.GalleryImageURL5 != "" {
		invitationData.GalleryImageURL5 = request.GalleryImageURL5
	}
	if request.GalleryImageURL6 != "" {
		invitationData.GalleryImageURL6 = request.GalleryImageURL6
	}
	if request.GalleryImageURL7 != "" {
		invitationData.GalleryImageURL7 = request.GalleryImageURL7
	}
	if request.GalleryImageURL8 != "" {
		invitationData.GalleryImageURL8 = request.GalleryImageURL8
	}
	if request.GalleryImageURL9 != "" {
		invitationData.GalleryImageURL9 = request.GalleryImageURL9
	}
	if request.GalleryImageURL10 != "" {
		invitationData.GalleryImageURL10 = request.GalleryImageURL10
	}
	if request.GalleryImageURL11 != "" {
		invitationData.GalleryImageURL11 = request.GalleryImageURL11
	}
	if request.GalleryImageURL12 != "" {
		invitationData.GalleryImageURL12 = request.GalleryImageURL12
	}

	err = h.InvitationDataRepositories.UpdateInvitationData(invitationData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: "Update Success"}
	json.NewEncoder(w).Encode(response)
}
