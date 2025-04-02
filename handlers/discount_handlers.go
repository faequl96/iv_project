package handlers

import (
	"encoding/json"
	discount_dto "iv_project/dto/discount"
	invitation_theme_dto "iv_project/dto/invitation_theme"
	iv_coin_package_dto "iv_project/dto/iv_coin_package"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/pkg/utils"
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
	discountResponse := discount_dto.DiscountResponse{}

	if len(invitationThemes) != 0 {
		var invitationThemeResponses []invitation_theme_dto.InvitationThemeResponse
		for _, invitationTheme := range invitationThemes {
			invitationThemeResponses = append(invitationThemeResponses, ConvertToInvitationThemeResponse(&invitationTheme))
		}
		discountResponse.InvitationThemes = invitationThemeResponses
	}

	if len(ivCoinPackages) != 0 {
		var ivCoinPackageResponses []iv_coin_package_dto.IVCoinPackageResponse
		for _, ivCoinPackage := range ivCoinPackages {
			ivCoinPackageResponses = append(ivCoinPackageResponses, ConvertToIVCoinPackageResponse(&ivCoinPackage))
		}
		discountResponse.IVCoinPackages = ivCoinPackageResponses
	}

	return discountResponse
}

func (h *dicountHandlers) SetProductPrices(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() && role != models.UserRoleAdmin.String() {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "You do not have permission to access this resource.",
			"id": "Anda tidak memiliki izin untuk mengakses sumber daya ini.",
		}
		ErrorResponse(w, http.StatusForbidden, messages, lang)
		return
	}

	var request discount_dto.DiscountRequest
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

	invitationThemes, err := h.InvitationThemeRepositories.GetInvitationThemesByDiscountCategoryID(request.DiscountCategoryID)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while fetching invitation themes by discount category.",
			"id": "Terjadi kesalahan saat mengambil tema undangan berdasarkan kategori diskon.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	for index, invitationTheme := range invitationThemes {
		invitationThemes[index].IDRDiscountPrice = utils.CalculateDiscountedPrice(invitationTheme.IDRPrice, request.Percentage)
		invitationThemes[index].IVCDiscountPrice = utils.CalculateDiscountedPrice(invitationTheme.IVCPrice, request.Percentage)
		invitationTheme.IDRDiscountPrice = utils.CalculateDiscountedPrice(invitationTheme.IDRPrice, request.Percentage)
		invitationTheme.IVCDiscountPrice = utils.CalculateDiscountedPrice(invitationTheme.IVCPrice, request.Percentage)
		if err := h.InvitationThemeRepositories.UpdateInvitationTheme(&invitationTheme); err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "An error occurred while updating the invitation theme.",
				"id": "Terjadi kesalahan saat mengupdate tema undangan.",
			}
			ErrorResponse(w, http.StatusInternalServerError, messages, lang)
			return
		}
	}

	ivCoinPackages, err := h.IVCoinPackageRepositories.GetIVCoinPackagesByDiscountCategoryID(request.DiscountCategoryID)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while fetching iv coins by discount category.",
			"id": "Terjadi kesalahan saat mengambil iv coin berdasarkan kategori diskon.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	for index, ivCoinPackage := range ivCoinPackages {
		ivCoinPackages[index].IDRDiscountPrice = utils.CalculateDiscountedPrice(ivCoinPackage.IDRPrice, request.Percentage)
		ivCoinPackage.IDRDiscountPrice = utils.CalculateDiscountedPrice(ivCoinPackage.IDRPrice, request.Percentage)
		if err := h.IVCoinPackageRepositories.UpdateIVCoinPackage(&ivCoinPackage); err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "An error occurred while updating the iv coin package.",
				"id": "Terjadi kesalahan saat mengupdate paket iv coin.",
			}
			ErrorResponse(w, http.StatusInternalServerError, messages, lang)
			return
		}
	}

	SuccessResponse(w, http.StatusOK, "Prices updated successfully", ConvertToDiscountResponse(invitationThemes, ivCoinPackages))
}
