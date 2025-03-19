package handlers

import (
	"encoding/json"
	iv_coin_package_dto "iv_project/dto/iv_coin_package"
	"iv_project/models"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ivCoinPackageHandlers struct {
	IVCoinPackageRepositories repositories.IVCoinPackageRepositories
}

func IVCoinPackageHandler(IVCoinPackageRepositories repositories.IVCoinPackageRepositories) *ivCoinPackageHandlers {
	return &ivCoinPackageHandlers{IVCoinPackageRepositories}
}

func ConvertToIVCoinPackageResponse(ivCoinPackage *models.IVCoinPackage) iv_coin_package_dto.IVCoinPackageResponse {
	return iv_coin_package_dto.IVCoinPackageResponse{
		ID:               ivCoinPackage.ID,
		Name:             ivCoinPackage.Name,
		CoinAmount:       ivCoinPackage.CoinAmount,
		IDRPrice:         ivCoinPackage.IDRPrice,
		IDRDiscountPrice: ivCoinPackage.IDRDiscountPrice,
		IVCPrice:         ivCoinPackage.IVCPrice,
		IVCDiscountPrice: ivCoinPackage.IVCDiscountPrice,
	}
}

func (h *ivCoinPackageHandlers) CreateIVCoinPackage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request iv_coin_package_dto.CreateIVCoinPackageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Failed to parse request: invalid JSON format")
		return
	}

	ivCoinPackage := &models.IVCoinPackage{
		Name:             request.Name,
		CoinAmount:       request.CoinAmount,
		IDRPrice:         request.IDRPrice,
		IDRDiscountPrice: request.IDRDiscountPrice,
		IVCPrice:         request.IVCPrice,
		IVCDiscountPrice: request.IVCDiscountPrice,
	}

	if err := h.IVCoinPackageRepositories.CreateIVCoinPackage(ivCoinPackage); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Error occurred while creating IV coin package. Please try again later.")
		return
	}

	SuccessResponse(w, http.StatusCreated, "IV coin package created successfully", ConvertToIVCoinPackageResponse(ivCoinPackage))
}

func (h *ivCoinPackageHandlers) GetIVCoinPackageByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid IV coin package ID format. Please provide a numeric ID.")
		return
	}

	ivCoinPackage, err := h.IVCoinPackageRepositories.GetIVCoinPackageByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No IV coin package found with the provided ID.")
		return
	}

	SuccessResponse(w, http.StatusOK, "IV coin package retrieved successfully", ConvertToIVCoinPackageResponse(ivCoinPackage))
}

func (h *ivCoinPackageHandlers) GetIVCoinPackages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ivCoinPackages, err := h.IVCoinPackageRepositories.GetIVCoinPackages()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while fetching IV coin packages.")
		return
	}

	if len(ivCoinPackages) == 0 {
		SuccessResponse(w, http.StatusOK, "No iv coin packages available at this moment", []iv_coin_package_dto.IVCoinPackageResponse{})
		return
	}

	var responses []iv_coin_package_dto.IVCoinPackageResponse
	for _, ivCoinPackage := range ivCoinPackages {
		responses = append(responses, ConvertToIVCoinPackageResponse(&ivCoinPackage))
	}

	SuccessResponse(w, http.StatusOK, "IV coin packages retrieved successfully", responses)
}

func (h *ivCoinPackageHandlers) UpdateIVCoinPackage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid IV coin package ID format. Please provide a numeric ID.")
		return
	}

	ivCoinPackage, err := h.IVCoinPackageRepositories.GetIVCoinPackageByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No IV coin package found with the provided ID.")
		return
	}

	var request iv_coin_package_dto.UpdateIVCoinPackageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Failed to parse request: invalid JSON format")
		return
	}

	if request.Name != "" {
		ivCoinPackage.Name = request.Name
	}
	if request.CoinAmount != 0 {
		ivCoinPackage.CoinAmount = request.CoinAmount
	}
	if request.IDRPrice != 0 {
		ivCoinPackage.IDRPrice = request.IDRPrice
	}
	ivCoinPackage.IDRDiscountPrice = request.IDRDiscountPrice
	if request.IVCPrice != 0 {
		ivCoinPackage.IVCPrice = request.IVCPrice
	}
	ivCoinPackage.IVCDiscountPrice = request.IVCDiscountPrice

	if err := h.IVCoinPackageRepositories.UpdateIVCoinPackage(ivCoinPackage); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while updating the IV coin package.")
		return
	}

	SuccessResponse(w, http.StatusOK, "IV coin package updated successfully", ConvertToIVCoinPackageResponse(ivCoinPackage))
}

func (h *ivCoinPackageHandlers) DeleteIVCoinPackage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid IV coin package ID format. Please provide a numeric ID.")
		return
	}

	if _, err = h.IVCoinPackageRepositories.GetIVCoinPackageByID(uint(id)); err != nil {
		ErrorResponse(w, http.StatusNotFound, "No IV coin package found with the provided ID.")
		return
	}

	if err := h.IVCoinPackageRepositories.DeleteIVCoinPackage(uint(id)); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while deleting the IV coin package.")
		return
	}

	SuccessResponse(w, http.StatusOK, "IV coin package deleted successfully", nil)
}
