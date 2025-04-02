package handlers

import (
	"encoding/json"
	voucher_code_dto "iv_project/dto/voucher_code"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type voucherCodeHandlers struct {
	VoucherCodeRepositories repositories.VoucherCodeRepositories
}

func VoucherCodeHandler(VoucherCodeRepositories repositories.VoucherCodeRepositories) *voucherCodeHandlers {
	return &voucherCodeHandlers{VoucherCodeRepositories}
}

func ConvertToVoucherCodeResponse(VoucherCode *models.VoucherCode) voucher_code_dto.VoucherCodeResponse {
	voucherCodeResponse := voucher_code_dto.VoucherCodeResponse{
		ID:                 VoucherCode.ID,
		Name:               VoucherCode.Name,
		DiscountPercentage: VoucherCode.DiscountPercentage,
	}

	return voucherCodeResponse
}

func (h *voucherCodeHandlers) CreateVoucherCode(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() && role != models.UserRoleAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	var request voucher_code_dto.VoucherCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Failed to parse request: invalid JSON format")
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	voucherCode := &models.VoucherCode{
		Name:               request.Name,
		DiscountPercentage: request.DiscountPercentage,
	}

	if err := h.VoucherCodeRepositories.CreateVoucherCode(voucherCode); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Error occurred while creating voucher code. Please try again later.")
		return
	}

	SuccessResponse(w, http.StatusCreated, "Voucher code created successfully", ConvertToVoucherCodeResponse(voucherCode))
}

func (h *voucherCodeHandlers) GetVoucherCodeByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid voucher code ID format. Please provide a numeric ID.")
		return
	}

	voucherCode, err := h.VoucherCodeRepositories.GetVoucherCodeByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No voucher code found with the provided ID.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Voucher code retrieved successfully", ConvertToVoucherCodeResponse(voucherCode))
}

func (h *voucherCodeHandlers) UpdateVoucherCodeByID(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() && role != models.UserRoleAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	var request voucher_code_dto.VoucherCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Failed to parse request: invalid JSON format")
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid voucher code ID format. Please provide a numeric ID.")
		return
	}

	voucherCode, err := h.VoucherCodeRepositories.GetVoucherCodeByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No voucher code found with the provided ID.")
		return
	}

	if request.Name != "" && request.DiscountPercentage != 0 {
		voucherCode.Name = request.Name
		voucherCode.DiscountPercentage = request.DiscountPercentage
	}

	if err := h.VoucherCodeRepositories.UpdateVoucherCode(voucherCode); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while updating the voucher code.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Voucher code updated successfully", ConvertToVoucherCodeResponse(voucherCode))
}

func (h *voucherCodeHandlers) DeleteVoucherCodeByID(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() && role != models.UserRoleAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid voucher code ID format. Please provide a numeric ID.")
		return
	}

	if _, err = h.VoucherCodeRepositories.GetVoucherCodeByID(uint(id)); err != nil {
		ErrorResponse(w, http.StatusNotFound, "No voucher code found with the provided ID.")
		return
	}

	if err := h.VoucherCodeRepositories.DeleteVoucherCode(uint(id)); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while deleting the voucher code.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Voucher code deleted successfully", nil)
}
