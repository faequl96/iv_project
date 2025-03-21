package handlers

import (
	"encoding/json"
	discount_dto "iv_project/dto/discount"
	invitation_theme_dto "iv_project/dto/invitation_theme"
	iv_coin_package_dto "iv_project/dto/iv_coin_package"
	"iv_project/models"
	"iv_project/repositories"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type dicountHandlers struct {
	InvitationThemeRepositories repositories.InvitationThemeRepositories
	IVCoinPackageRepositories   repositories.IVCoinPackageRepositories
}

func DiscountHandlers(
	InvitationThemeRepositories repositories.InvitationThemeRepositories,
	IVCoinPackageRepositories repositories.IVCoinPackageRepositories,
) *dicountHandlers {
	return &dicountHandlers{InvitationThemeRepositories, IVCoinPackageRepositories}
}

func ConvertToDiscountResponse(
	invitationThemes []models.InvitationTheme,
	ivCoinPackages []models.IVCoinPackage,
) discount_dto.DiscountResponse {
	var invitationThemeResponses []invitation_theme_dto.InvitationThemeResponse
	for _, invitationTheme := range invitationThemes {
		invitationThemeResponses = append(invitationThemeResponses, ConvertToInvitationThemeResponse(&invitationTheme))
	}

	var ivCoinPackageResponses []iv_coin_package_dto.IVCoinPackageResponse
	for _, ivCoinPackage := range ivCoinPackages {
		ivCoinPackageResponses = append(ivCoinPackageResponses, ConvertToIVCoinPackageResponse(&ivCoinPackage))
	}

	return discount_dto.DiscountResponse{
		InvitationThemes: invitationThemeResponses,
		IVCoinPackages:   ivCoinPackageResponses,
	}
}

func (h *dicountHandlers) SetProductPrices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request discount_dto.DiscountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	invitationThemes, err := h.InvitationThemeRepositories.GetInvitationThemesByDiscountCategory(request.DiscountCategory)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while fetching invitation themes by discount category.")
		return
	}

	for _, invitationTheme := range invitationThemes {
		invitationTheme.IDRDiscountPrice = CalculateDiscountedPrice(invitationTheme.IDRPrice, request.Percentage)
		if err := h.InvitationThemeRepositories.UpdateInvitationTheme(&invitationTheme); err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "An error occurred while updating the invitation theme.")
			return
		}
	}

	ivCoinPackages, err := h.IVCoinPackageRepositories.GetIVCoinPackagesByDiscountCategory(request.DiscountCategory)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while fetching iv coins by discount category.")
		return
	}

	for _, ivCoinPackage := range ivCoinPackages {
		ivCoinPackage.IDRDiscountPrice = CalculateDiscountedPrice(ivCoinPackage.IDRPrice, request.Percentage)
		if err := h.IVCoinPackageRepositories.UpdateIVCoinPackage(&ivCoinPackage); err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "An error occurred while updating the iv coin package.")
			return
		}
	}

	SuccessResponse(w, http.StatusOK, "Prices updated successfully", ConvertToDiscountResponse(invitationThemes, ivCoinPackages))
}
