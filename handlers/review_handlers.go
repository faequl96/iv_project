package handlers

import (
	"encoding/json"
	review_dto "iv_project/dto/review"
	"iv_project/models"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type reviewHandlers struct {
	ReviewRepositories repositories.ReviewRepositories
}

func ReviewHandlers(ReviewRepositories repositories.ReviewRepositories) *reviewHandlers {
	return &reviewHandlers{ReviewRepositories}
}

func ConvertToReviewResponse(review *models.Review) review_dto.ReviewResponse {
	return review_dto.ReviewResponse{
		ID:      review.ID,
		Star:    review.Star,
		Comment: review.Comment,
	}
}

func (h *reviewHandlers) CreateReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request review_dto.CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request format. Please check your input.")
		return
	}

	review := &models.Review{
		UserID:            request.UserID,
		InvitationThemeID: request.InvitationThemeID,
		Star:              request.Star,
		Comment:           request.Comment,
	}

	if err := h.ReviewRepositories.CreateReview(review); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while create the review. Please try again.")
		return
	}

	SuccessResponse(w, http.StatusCreated, "Review successfully created", ConvertToReviewResponse(review))
}

func (h *reviewHandlers) GetReviewByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID format. Please use a number.")
		return
	}

	review, err := h.ReviewRepositories.GetReviewByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Review not found with the given ID.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Review found successfully", ConvertToReviewResponse(review))
}

func (h *reviewHandlers) GetReviewsByInvitationThemeID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["invitationThemeId"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID format. Please use a number.")
		return
	}

	reviews, err := h.ReviewRepositories.GetReviewsByInvitationThemeID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while retrieving reviews. Please try again.")
		return
	}

	if len(reviews) == 0 {
		SuccessResponse(w, http.StatusOK, "No reviews available for this invitation theme", []review_dto.ReviewResponse{})
		return
	}

	var response []review_dto.ReviewResponse
	for _, review := range reviews {
		response = append(response, ConvertToReviewResponse(&review))
	}

	SuccessResponse(w, http.StatusOK, "Reviews retrieved successfully", response)
}

func (h *reviewHandlers) UpdateReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID format. Please use a number.")
		return
	}

	review, err := h.ReviewRepositories.GetReviewByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Review not found with the given ID.")
		return
	}

	var request review_dto.UpdateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request format. Please check your input.")
		return
	}

	if request.Star != 0 {
		review.Star = request.Star
	}
	if request.Comment != "" {
		review.Comment = request.Comment
	}

	if err := h.ReviewRepositories.UpdateReview(review); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while updating the review. Please try again.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Review successfully updated", ConvertToReviewResponse(review))
}

func (h *reviewHandlers) DeleteReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID format. Please use a number.")
		return
	}

	if err := h.ReviewRepositories.DeleteReview(uint(id)); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while deleting the review. Please try again.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Review successfully deleted", nil)
}
