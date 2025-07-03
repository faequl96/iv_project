package handlers

import (
	"errors"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/pkg/utils"
	"iv_project/repositories"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"gorm.io/gorm"
)

type transactionStatusHandlers struct {
	TransactionRepositories          repositories.TransactionRepositories
	InvitationRepositories           repositories.InvitationRepositories
	IVCoinPackageRepositories        repositories.IVCoinPackageRepositories
	IVCoinRepositories               repositories.IVCoinRepositories
	VoucherCodeRepositories          repositories.VoucherCodeRepositories
	UserVoucherCodeUsageRepositories repositories.UserVoucherCodeUsageRepositories
}

func TransactionStatusHandler(
	TransactionRepositories repositories.TransactionRepositories,
	InvitationRepositories repositories.InvitationRepositories,
	IVCoinPackageRepositories repositories.IVCoinPackageRepositories,
	IVCoinRepositories repositories.IVCoinRepositories,
	VoucherCodeRepositories repositories.VoucherCodeRepositories,
	UserVoucherCodeUsageRepositories repositories.UserVoucherCodeUsageRepositories,
) *transactionStatusHandlers {
	return &transactionStatusHandlers{
		TransactionRepositories,
		InvitationRepositories,
		IVCoinPackageRepositories,
		IVCoinRepositories,
		VoucherCodeRepositories,
		UserVoucherCodeUsageRepositories,
	}
}

func (h *transactionStatusHandlers) CheckByReferenceNumber(w http.ResponseWriter, r *http.Request) {
	referenceNumber := mux.Vars(r)["referenceNumber"]

	var c = coreapi.Client{}
	c.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)
	resp, midtransErr := c.CheckTransaction(referenceNumber)
	if midtransErr != nil {
		if strings.Contains(midtransErr.Message, "Transaction doesn't exist.") {
			transaction, err := h.TransactionRepositories.GetTransactionByReferenceNumber(referenceNumber)
			if err != nil {
				lang, _ := r.Context().Value(middleware.LanguageKey).(string)
				messages := map[string]string{
					"en": "No transaction found with the provided Reference Number.",
					"id": "Transaksi tidak ditemukan dengan Nomor Referensi yang diberikan.",
				}
				ErrorResponse(w, http.StatusNotFound, messages, lang)
				return
			}
			SuccessResponse(w, http.StatusOK, "Transaction check successfully", ConvertToTransactionResponse(transaction))
			return
		} else {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "An error occurred while connecting the payment gateway",
				"id": "Terjadi kesalahan saat menghubungkan payment gateway",
			}
			ErrorResponse(w, http.StatusNotFound, messages, lang)
			return
		}
	}

	transaction, err := h.TransactionRepositories.GetTransactionByReferenceNumber(referenceNumber)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No transaction found with the provided Reference Number.",
			"id": "Transaksi tidak ditemukan dengan Nomor Referensi yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if transaction.Status != models.TransactionStatusConfirmed {
		switch resp.TransactionStatus {
		case "settlement", "capture":
			if transaction.ProductType == models.ProductInvitation {
				invitation, err := h.InvitationRepositories.GetInvitationByID(uint(transaction.ProductID))
				if err != nil {
					lang, _ := r.Context().Value(middleware.LanguageKey).(string)
					messages := map[string]string{
						"en": "Invalid Invitation ID format. Please provide a numeric ID.",
						"id": "Format ID undangan tidak valid. Harap berikan ID dalam format angka.",
					}
					ErrorResponse(w, http.StatusBadRequest, messages, lang)
					return
				}
				invitation.Status = models.InvitationStatusActive
				err = h.InvitationRepositories.UpdateInvitation(invitation)
				if err != nil {
					lang, _ := r.Context().Value(middleware.LanguageKey).(string)
					messages := map[string]string{
						"en": "Failed to update invitation",
						"id": "Gagal mengupdate undangan",
					}
					ErrorResponse(w, http.StatusInternalServerError, messages, lang)
					return
				}
			}

			if transaction.ProductType == models.ProductIVCoinPackage {
				ivCoinPackage, err := h.IVCoinPackageRepositories.GetIVCoinPackageByID(uint(transaction.ProductID))
				if err != nil {
					lang, _ := r.Context().Value(middleware.LanguageKey).(string)
					messages := map[string]string{
						"en": "No iv coin package found with the provided ID.",
						"id": "Paket iv coin tidak ditemukan dengan ID yang diberikan.",
					}
					ErrorResponse(w, http.StatusNotFound, messages, lang)
					return
				}

				ivCoin, err := h.IVCoinRepositories.GetIVCoinByUserID(transaction.UserID)
				if err != nil {
					lang, _ := r.Context().Value(middleware.LanguageKey).(string)
					messages := map[string]string{
						"en": "No iv coin found with the provided user.",
						"id": "IV coin tidak ditemukan dengan pengguna yang diberikan.",
					}
					ErrorResponse(w, http.StatusNotFound, messages, lang)
					return
				}

				ivCoin.Balance = ivCoin.Balance + ivCoinPackage.CoinAmount

				err = h.IVCoinRepositories.UpdateIVCoin(ivCoin)
				if err != nil {
					lang, _ := r.Context().Value(middleware.LanguageKey).(string)
					messages := map[string]string{
						"en": "An error occurred while updating the iv coin.",
						"id": "Terjadi kesalahan saat mengupdate iv coin.",
					}
					ErrorResponse(w, http.StatusInternalServerError, messages, lang)
					return
				}
			}

			if resp.TransactionStatus == "settlement" {
				transaction.MidtransStatus = models.MidtransTransactionStatusSettlement
			}
			if resp.TransactionStatus == "capture" {
				transaction.MidtransStatus = models.MidtransTransactionStatusCapture
			}
			transaction.Status = models.TransactionStatusConfirmed
			transaction.PaymentURL = ""
			transaction.TimeLimitAt = nil

			if transaction.VoucherCodeName != "" {
				voucherCode, err := h.VoucherCodeRepositories.GetVoucherCodeByName(transaction.VoucherCodeName)
				if err != nil {
					lang, _ := r.Context().Value(middleware.LanguageKey).(string)
					messages := map[string]string{
						"en": "No voucher code found with the provided name.",
						"id": "Kode voucher tidak ditemukan dengan nama yang diberikan.",
					}
					ErrorResponse(w, http.StatusNotFound, messages, lang)
					return
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

				if errors.Is(err, gorm.ErrRecordNotFound) {
					userVoucherCodeUsage = &models.UserVoucherCodeUsage{
						UserID:        transaction.UserID,
						VoucherCodeID: voucherCode.ID,
						UsageCount:    1,
					}
					err = h.UserVoucherCodeUsageRepositories.CreateUserVoucherCodeUsage(userVoucherCodeUsage)
					if err != nil {
						lang, _ := r.Context().Value(middleware.LanguageKey).(string)
						messages := map[string]string{
							"en": "Failed to create voucher code usage.",
							"id": "Gagal membuat penggunaan kode voucher",
						}
						ErrorResponse(w, http.StatusInternalServerError, messages, lang)
						return
					}
				} else {
					userVoucherCodeUsage.UsageCount += 1
					err = h.UserVoucherCodeUsageRepositories.UpdateUserVoucherCodeUsage(userVoucherCodeUsage)
					if err != nil {
						lang, _ := r.Context().Value(middleware.LanguageKey).(string)
						messages := map[string]string{
							"en": "Failed to update voucher code usage.",
							"id": "Gagal memperbarui penggunaan kode voucher",
						}
						ErrorResponse(w, http.StatusInternalServerError, messages, lang)
						return
					}
				}
			}
		case "pending":
			transaction.MidtransStatus = models.MidtransTransactionStatusPending
			transaction.Status = models.TransactionStatusPending
		case "expire":
			transaction.MidtransStatus = models.MidtransTransactionStatusExpire
			transaction.Status = models.TransactionStatusCreated
			transaction.PaymentURL = ""
			transaction.TimeLimitAt = nil
			transaction.Acquirer = ""
		case "cancel":
			transaction.MidtransStatus = models.MidtransTransactionStatusCancel
			transaction.Status = models.TransactionStatusCreated
			transaction.PaymentURL = ""
			transaction.TimeLimitAt = nil
			transaction.Acquirer = ""
		case "deny":
			transaction.MidtransStatus = models.MidtransTransactionStatusDeny
			transaction.Status = models.TransactionStatusCreated
			transaction.PaymentURL = ""
			transaction.TimeLimitAt = nil
			transaction.Acquirer = ""
		default:
			transaction.MidtransStatus = models.MidtransTransactionStatusUnknown
			transaction.Status = models.TransactionStatusCreated
			transaction.PaymentURL = ""
			transaction.TimeLimitAt = nil
			transaction.Acquirer = ""
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

	SuccessResponse(w, http.StatusOK, "Transaction check successfully", ConvertToTransactionResponse(transaction))
}

func (h *transactionStatusHandlers) ResetByID(w http.ResponseWriter, r *http.Request) {
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

	transaction.ReferenceNumber = utils.GenerateReferenceNumber(transaction.PaymentMethod.String())
	transaction.MidtransStatus = models.MidtransTransactionStatusUnknown
	transaction.Status = models.TransactionStatusCreated
	transaction.PaymentURL = ""
	transaction.Acquirer = ""
	transaction.TimeLimitAt = nil

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

	SuccessResponse(w, http.StatusOK, "Transaction reset successfully", ConvertToTransactionResponse(transaction))

}
