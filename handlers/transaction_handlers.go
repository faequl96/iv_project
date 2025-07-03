package handlers

import (
	"encoding/json"
	"errors"
	query_dto "iv_project/dto/query"
	transaction_dto "iv_project/dto/transaction"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/pkg/utils"
	"iv_project/repositories"
	"net/http"
	"time"

	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type transactionHandlers struct {
	TransactionRepositories          repositories.TransactionRepositories
	InvitationRepositories           repositories.InvitationRepositories
	IVCoinPackageRepositories        repositories.IVCoinPackageRepositories
	IVCoinRepositories               repositories.IVCoinRepositories
	UserRepositories                 repositories.UserRepositories
	VoucherCodeRepositories          repositories.VoucherCodeRepositories
	UserVoucherCodeUsageRepositories repositories.UserVoucherCodeUsageRepositories
}

func TransactionHandler(
	TransactionRepositories repositories.TransactionRepositories,
	InvitationRepositories repositories.InvitationRepositories,
	IVCoinPackageRepositories repositories.IVCoinPackageRepositories,
	IVCoinRepositories repositories.IVCoinRepositories,
	UserRepositories repositories.UserRepositories,
	VoucherCodeRepositories repositories.VoucherCodeRepositories,
	UserVoucherCodeUsageRepositories repositories.UserVoucherCodeUsageRepositories,
) *transactionHandlers {
	return &transactionHandlers{
		TransactionRepositories,
		InvitationRepositories,
		IVCoinPackageRepositories,
		IVCoinRepositories,
		UserRepositories,
		VoucherCodeRepositories,
		UserVoucherCodeUsageRepositories,
	}
}

func ConvertToTransactionResponse(transaction *models.Transaction) transaction_dto.TransactionResponse {
	transactionResponse := transaction_dto.TransactionResponse{
		ID:                     transaction.ID,
		TransactionCode:        transaction.TransactionCode,
		ProductType:            transaction.ProductType,
		ProductName:            transaction.ProductName,
		ProductDescription:     transaction.ProductDescription,
		Status:                 transaction.Status,
		PaymentMethod:          transaction.PaymentMethod,
		ReferenceNumber:        transaction.ReferenceNumber,
		PaymentURL:             transaction.PaymentURL,
		Acquirer:               transaction.Acquirer,
		MidtransStatus:         transaction.MidtransStatus,
		TimeLimitAt:            transaction.TimeLimitAt,
		IDRPrice:               transaction.IDRPrice,
		IDRDiscount:            transaction.IDRDiscount,
		IDRTotalPrice:          transaction.IDRTotalPrice,
		IVCPrice:               transaction.IVCPrice,
		IVCDiscount:            transaction.IVCDiscount,
		IVCTotalPrice:          transaction.IVCTotalPrice,
		VoucherCodeID:          transaction.VoucherCodeID,
		VoucherCodeName:        transaction.VoucherCodeName,
		IDRVoucherCodeDiscount: transaction.IDRVoucherCodeDiscount,
		IVCVoucherCodeDiscount: transaction.IVCVoucherCodeDiscount,
		PaymentProofImageUrl:   transaction.PaymentProofImageUrl,
		CreatedAt:              transaction.CreatedAt.Format(time.RFC3339),
	}

	return transactionResponse
}

func (h *transactionHandlers) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var request transaction_dto.CreateTransactionRequest
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

	transaction := &models.Transaction{
		ID:              uuid.New().String(),
		TransactionCode: utils.GenerateTransactionCode(),
		ProductType:     request.ProductType,
		Status:          models.TransactionStatusCreated,
		MidtransStatus:  models.MidtransTransactionStatusUnknown,
		ProductID:       request.ProductID,
		UserID:          request.UserID,
	}

	if request.ProductType == models.ProductInvitation {
		invitation, err := h.InvitationRepositories.GetInvitationByID(uint(request.ProductID))
		if err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "No invitation found with the provided ID.",
				"id": "Undangan tidak ditemukan dengan ID yang diberikan.",
			}
			ErrorResponse(w, http.StatusNotFound, messages, lang)
			return
		}

		transaction.IDRPrice = invitation.InvitationTheme.IDRPrice
		transaction.IDRDiscount = invitation.InvitationTheme.IDRPrice - invitation.InvitationTheme.IDRDiscountPrice
		transaction.IDRTotalPrice = invitation.InvitationTheme.IDRDiscountPrice
		transaction.IVCPrice = invitation.InvitationTheme.IVCPrice
		transaction.IVCDiscount = invitation.InvitationTheme.IVCPrice - invitation.InvitationTheme.IVCDiscountPrice
		transaction.IVCTotalPrice = invitation.InvitationTheme.IVCDiscountPrice

		transaction.ProductName = invitation.InvitationTheme.Name

		ivCoin, _ := h.IVCoinRepositories.GetIVCoinByUserID(transaction.UserID)
		if ivCoin != nil {
			if ivCoin.Balance > transaction.IVCTotalPrice {
				transaction.PaymentMethod = models.PaymentMethodIVCoin
				transaction.ReferenceNumber = utils.GenerateReferenceNumber(transaction.PaymentMethod.String())
			} else {
				transaction.PaymentMethod = models.PaymentMethodGopay
				transaction.ReferenceNumber = utils.GenerateReferenceNumber(transaction.PaymentMethod.String())
			}
		} else {
			transaction.PaymentMethod = models.PaymentMethodGopay
			transaction.ReferenceNumber = utils.GenerateReferenceNumber(transaction.PaymentMethod.String())
		}
	}

	if request.ProductType == models.ProductIVCoinPackage {
		ivCoinPackage, err := h.IVCoinPackageRepositories.GetIVCoinPackageByID(uint(request.ProductID))
		if err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "No iv coin package found with the provided ID.",
				"id": "Paket iv coin tidak ditemukan dengan ID yang diberikan.",
			}
			ErrorResponse(w, http.StatusNotFound, messages, lang)
			return
		}

		transaction.IDRPrice = ivCoinPackage.IDRPrice
		transaction.IDRDiscount = ivCoinPackage.IDRPrice - ivCoinPackage.IDRDiscountPrice
		transaction.IDRTotalPrice = ivCoinPackage.IDRDiscountPrice

		transaction.ProductName = ivCoinPackage.Name
		transaction.ProductDescription = fmt.Sprint(ivCoinPackage.CoinAmount)

		transaction.PaymentMethod = models.PaymentMethodGopay
		transaction.ReferenceNumber = utils.GenerateReferenceNumber(transaction.PaymentMethod.String())
	}

	err := h.TransactionRepositories.CreateTransaction(transaction)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Failed to create transaction.",
			"id": "Gagal membuat transaksi.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusCreated, "Transaction created successfully", ConvertToTransactionResponse(transaction))
}

func (h *transactionHandlers) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	transaction, err := h.TransactionRepositories.GetTransactionByID(id)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No transaction found with the provided ID.",
			"id": "Transaksi tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Transaction retrieved successfully", ConvertToTransactionResponse(transaction))
}

func (h *transactionHandlers) GetTransactionByReferenceNumber(w http.ResponseWriter, r *http.Request) {
	referenceNumber := mux.Vars(r)["referenceNumber"]

	transaction, err := h.TransactionRepositories.GetTransactionByReferenceNumber(referenceNumber)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No transaction found with the provided ID.",
			"id": "Transaksi tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Transaction retrieved successfully", ConvertToTransactionResponse(transaction))
}

func (h *transactionHandlers) GetTransactions(w http.ResponseWriter, r *http.Request) {
	var request *query_dto.QueryRequest
	json.NewDecoder(r.Body).Decode(&request)

	transactions, err := h.TransactionRepositories.GetTransactions(request)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while fetching transactions.",
			"id": "Terjadi kesalahan saat mengambil transaksi.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	if len(transactions) == 0 {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No transactions available at the moment.",
			"id": "Tidak ada transaksi saat ini.",
		}
		SuccessResponse(w, http.StatusOK, messages[lang], []transaction_dto.TransactionResponse{})
		return
	}

	var transactionResponses []transaction_dto.TransactionResponse
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, ConvertToTransactionResponse(&transaction))
	}

	SuccessResponse(w, http.StatusOK, "Transactions retrieved successfully", transactionResponses)
}

func (h *transactionHandlers) GetTransactionsByUserID(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userId"]

	transactions, err := h.TransactionRepositories.GetTransactionsByUserID(userID)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while fetching transactions.",
			"id": "Terjadi kesalahan saat mengambil transaksi.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	if len(transactions) == 0 {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No transactions available at the moment.",
			"id": "Tidak ada transaksi saat ini.",
		}
		SuccessResponse(w, http.StatusOK, messages[lang], []transaction_dto.TransactionResponse{})
		return
	}

	var transactionResponses []transaction_dto.TransactionResponse
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, ConvertToTransactionResponse(&transaction))
	}

	SuccessResponse(w, http.StatusOK, "Transactions retrieved successfully", transactionResponses)
}

func (h *transactionHandlers) UpdateTransactionByID(w http.ResponseWriter, r *http.Request) {
	var request transaction_dto.UpdateTransactionRequest
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

	id := mux.Vars(r)["id"]

	transaction, err := h.TransactionRepositories.GetTransactionByID(id)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No transaction found with the provided ID.",
			"id": "Transaksi tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if transaction.Status == models.TransactionStatusConfirmed || transaction.Status == models.TransactionStatusCanceled {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "The transaction cannot be updated, because the transaction has been completed.",
			"id": "Transaksi tidak dapat diupdate, karena transaksi telah selesai.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	transaction.MidtransStatus = models.MidtransTransactionStatusUnknown

	if request.PaymentMethod.String() != "" {
		transaction.PaymentMethod = request.PaymentMethod
		transaction.ReferenceNumber = utils.GenerateReferenceNumber(request.PaymentMethod.String())
	}

	if request.VoucherCodeName != "" {
		voucherCode, err := h.VoucherCodeRepositories.GetVoucherCodeByName(request.VoucherCodeName)
		if err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "Voucher code not available",
				"id": "Kode voucher tidak tersedia",
			}
			ErrorResponse(w, http.StatusNotFound, messages, lang)
			return
		}
		if !voucherCode.IsGlobal {
			allowed := false
			for _, user := range voucherCode.Users {
				if user.ID == transaction.UserID {
					allowed = true
					break
				}
			}
			if !allowed {
				lang, _ := r.Context().Value(middleware.LanguageKey).(string)
				messages := map[string]string{
					"en": "Voucher code is not valid",
					"id": "Kode voucher tidak berlaku",
				}
				ErrorResponse(w, http.StatusNotFound, messages, lang)
				return
			}
		}

		userVoucherCodeUsage, err := h.UserVoucherCodeUsageRepositories.GetUserVoucherCodeUsageByUserAndVoucherCodeID(transaction.UserID, voucherCode.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "An error occurred while retrieving voucher code usage data.",
				"id": "Terjadi kesalahan saat mengambil data penggunaan kode voucher.",
			}
			ErrorResponse(w, http.StatusNotFound, messages, lang)
			return
		}

		if userVoucherCodeUsage != nil {
			if userVoucherCodeUsage.UsageCount >= voucherCode.UsageLimitPerUser {
				lang, _ := r.Context().Value(middleware.LanguageKey).(string)
				messages := map[string]string{
					"en": "Voucher usage has reached the limit",
					"id": "Penggunaan voucher telah mencapai batas",
				}
				ErrorResponse(w, http.StatusNotFound, messages, lang)
				return
			}
		}

		transaction.VoucherCodeID = voucherCode.ID
		transaction.VoucherCodeName = voucherCode.Name

		if transaction.ProductType == models.ProductInvitation {
			totalIDRVoucherCodeDiscount := utils.CalculateDiscountedPrice((transaction.IDRPrice - transaction.IDRDiscount), voucherCode.DiscountPercentage)
			transaction.IDRVoucherCodeDiscount = (transaction.IDRPrice - transaction.IDRDiscount) - totalIDRVoucherCodeDiscount
			transaction.IDRTotalPrice = totalIDRVoucherCodeDiscount
			totalIVCVoucherCodeDiscount := utils.CalculateDiscountedPrice((transaction.IVCPrice - transaction.IVCDiscount), voucherCode.DiscountPercentage)
			transaction.IVCVoucherCodeDiscount = (transaction.IVCPrice - transaction.IVCDiscount) - totalIVCVoucherCodeDiscount
			transaction.IVCTotalPrice = totalIVCVoucherCodeDiscount
		}

		if transaction.ProductType == models.ProductIVCoinPackage {
			totalIDRVoucherCodeDiscount := utils.CalculateDiscountedPrice((transaction.IDRPrice - transaction.IDRDiscount), voucherCode.DiscountPercentage)
			transaction.IDRVoucherCodeDiscount = (transaction.IDRPrice - transaction.IDRDiscount) - totalIDRVoucherCodeDiscount
			transaction.IDRTotalPrice = totalIDRVoucherCodeDiscount
		}
	}

	if request.VoucherCodeName == "" {
		transaction.VoucherCodeID = 0
		transaction.VoucherCodeName = ""

		if transaction.ProductType == models.ProductInvitation {
			transaction.IDRVoucherCodeDiscount = 0
			transaction.IDRTotalPrice = transaction.IDRPrice - transaction.IDRDiscount
			transaction.IVCVoucherCodeDiscount = 0
			transaction.IVCTotalPrice = transaction.IVCPrice - transaction.IVCDiscount
		}

		if transaction.ProductType == models.ProductIVCoinPackage {
			transaction.IDRVoucherCodeDiscount = 0
			transaction.IDRTotalPrice = transaction.IDRPrice - transaction.IDRDiscount
		}
	}

	err = h.TransactionRepositories.UpdateTransaction(transaction)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Failed to update transaction.",
			"id": "Gagal mengupdate transaksi.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusCreated, "Transaction created successfully", ConvertToTransactionResponse(transaction))
}

func (h *transactionHandlers) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if _, err := h.TransactionRepositories.GetTransactionByID(id); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No transaction found with the provided ID.",
			"id": "Transaksi tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if err := h.TransactionRepositories.DeleteTransaction(id); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while deleting the transaction.",
			"id": "Terjadi kesalahan saat menghapus transaksi.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Transaction deleted successfully", nil)
}
