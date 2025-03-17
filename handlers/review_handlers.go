package handlers

import (
	"encoding/json"
	"iv_project/dto"
	review_dto "iv_project/dto/review"
	"iv_project/models"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type reviewHandlers struct {
	ReviewRepositories repositories.ReviewRepositories
}

func ReviewHandlers(ReviewRepositories repositories.ReviewRepositories) *reviewHandlers {
	return &reviewHandlers{ReviewRepositories}
}

func (h *reviewHandlers) CreateReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(review_dto.ReviewRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	review := models.Review{
		UserID:            request.UserID,
		InvitationThemeID: request.InvitationThemeID,
		Star:              request.Star,
		Comment:           request.Comment,
	}

	err = h.ReviewRepositories.CreateReview(review)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.SuccessResult{Code: http.StatusOK, Data: "Review created successfully"})
}

func (h *reviewHandlers) GetReviewByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	review, err := h.ReviewRepositories.GetReviewByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusNotFound, Message: "Review not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.SuccessResult{Code: http.StatusOK, Data: review})
}

func (h *reviewHandlers) GetReviewsByInvitationThemeID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["invitationThemeId"])
	reviews, err := h.ReviewRepositories.GetReviewsByInvitationThemeID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.SuccessResult{Code: http.StatusOK, Data: reviews})
}

func (h *reviewHandlers) UpdateReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(review_dto.ReviewRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	review, err := h.ReviewRepositories.GetReviewByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusNotFound, Message: "Review not found"})
		return
	}

	if request.Star != 0 {
		review.Star = request.Star
	}
	if request.Comment != "" {
		review.Comment = request.Comment
	}

	err = h.ReviewRepositories.UpdateReview(review)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.SuccessResult{Code: http.StatusOK, Data: "Review updated successfully"})
}

func (h *reviewHandlers) DeleteReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	review, err := h.ReviewRepositories.GetReviewByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusNotFound, Message: "Review not found"})
		return
	}

	err = h.ReviewRepositories.DeleteReview(review)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.SuccessResult{Code: http.StatusOK, Data: "Review deleted successfully"})
}
