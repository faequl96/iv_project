package handlers

import (
	"errors"
	"fmt"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/repositories"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

type transactionPaymentHandlers struct {
	TransactionRepositories          repositories.TransactionRepositories
	InvitationRepositories           repositories.InvitationRepositories
	IVCoinPackageRepositories        repositories.IVCoinPackageRepositories
	IVCoinRepositories               repositories.IVCoinRepositories
	UserRepositories                 repositories.UserRepositories
	VoucherCodeRepositories          repositories.VoucherCodeRepositories
	UserVoucherCodeUsageRepositories repositories.UserVoucherCodeUsageRepositories
}

func TransactionPaymentHandler(
	TransactionRepositories repositories.TransactionRepositories,
	InvitationRepositories repositories.InvitationRepositories,
	IVCoinPackageRepositories repositories.IVCoinPackageRepositories,
	IVCoinRepositories repositories.IVCoinRepositories,
	UserRepositories repositories.UserRepositories,
	VoucherCodeRepositories repositories.VoucherCodeRepositories,
	UserVoucherCodeUsageRepositories repositories.UserVoucherCodeUsageRepositories,
) *transactionPaymentHandlers {
	return &transactionPaymentHandlers{
		TransactionRepositories,
		InvitationRepositories,
		IVCoinPackageRepositories,
		IVCoinRepositories,
		UserRepositories,
		VoucherCodeRepositories,
		UserVoucherCodeUsageRepositories,
	}
}

func (h *transactionPaymentHandlers) IssueByID(w http.ResponseWriter, r *http.Request) {
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

	if transaction.PaymentMethod == models.PaymentMethodIVCoin {
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

		if ivCoin.Balance < transaction.IVCTotalPrice {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": fmt.Sprintf("Insufficient IVCoin balance: %d/%d IVC.", ivCoin.Balance, transaction.IVCTotalPrice),
				"id": fmt.Sprintf("IV koin tidak cukup: %d/%d IVC.", ivCoin.Balance, transaction.IVCTotalPrice),
			}
			ErrorResponse(w, http.StatusPaymentRequired, messages, lang)
			return
		}
		transaction.Status = models.TransactionStatusConfirmed

		if transaction.ProductType == models.ProductInvitation {
			invitation, err := h.InvitationRepositories.GetInvitationByID(uint(transaction.ProductID))
			if err != nil {
				lang, _ := r.Context().Value(middleware.LanguageKey).(string)
				messages := map[string]string{
					"en": "No transaction found with the provided ID.",
					"id": "Transaksi tidak ditemukan dengan ID yang diberikan.",
				}
				ErrorResponse(w, http.StatusNotFound, messages, lang)
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

		ivCoin.Balance = ivCoin.Balance - transaction.IVCTotalPrice

		err = h.IVCoinRepositories.UpdateIVCoin(ivCoin)
		if err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "Failed to update iv coin.",
				"id": "Gagal mengupdate iv coin.",
			}
			ErrorResponse(w, http.StatusInternalServerError, messages, lang)
			return
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

		SuccessResponse(w, http.StatusOK, "Transaction issue successfully", ConvertToTransactionResponse(transaction))
		return
	}

	user, err := h.UserRepositories.GetUserByID(transaction.UserID)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No user found with the provided ID.",
			"id": "Pengguna tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	var url string

	if transaction.PaymentMethod == models.PaymentMethodBankTransfer {
		var s = snap.Client{}
		s.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

		snapReq := &snap.Request{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  transaction.ReferenceNumber,
				GrossAmt: int64(transaction.IDRTotalPrice),
			},
			CustomerDetail: &midtrans.CustomerDetails{
				FName: user.UserProfile.FirstName + " " + user.UserProfile.LastName,
				Email: user.UserProfile.Email,
			},
			Expiry:          &snap.ExpiryDetails{Duration: 24, Unit: "hour"},
			EnabledPayments: []snap.SnapPaymentType{},
		}

		snapReq.EnabledPayments = append(snapReq.EnabledPayments, snap.PaymentTypeBankTransfer)

		snapResp, _ := s.CreateTransaction(snapReq)
		url = snapResp.RedirectURL
		timeLimitAt := time.Now().Add(24 * time.Hour)
		transaction.TimeLimitAt = &timeLimitAt
	}

	if transaction.PaymentMethod == models.PaymentMethodGopay {
		var c coreapi.Client
		c.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

		chargeReq := &coreapi.ChargeReq{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  transaction.ReferenceNumber,
				GrossAmt: int64(transaction.IDRTotalPrice),
			},
			CustomerDetails: &midtrans.CustomerDetails{
				FName: user.UserProfile.FirstName + " " + user.UserProfile.LastName,
				Email: user.UserProfile.Email,
			},
			CustomExpiry: &coreapi.CustomExpiry{ExpiryDuration: 15, Unit: "minute"},
			PaymentType:  coreapi.PaymentTypeGopay,
			Gopay: &coreapi.GopayDetails{
				EnableCallback: true,
				CallbackUrl:    "ivprojectapp://payment-result",
			},
			Qris: &coreapi.QrisDetails{},
		}

		chargeResp, err := c.ChargeTransaction(chargeReq)
		if err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "Failed to connect Payment Gateway",
				"id": "Gagal menghubungkan ke Payment Gateway",
			}
			ErrorResponse(w, http.StatusInternalServerError, messages, lang)
			return
		}
		for _, action := range chargeResp.Actions {
			if action.Name == "deeplink-redirect" {
				url = action.URL
			}
		}
		timeLimitAt := time.Now().Add(15 * time.Minute)
		transaction.TimeLimitAt = &timeLimitAt
	}

	if transaction.PaymentMethod == models.PaymentMethodQRIS {
		var c coreapi.Client
		c.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

		chargeReq := &coreapi.ChargeReq{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  transaction.ReferenceNumber,
				GrossAmt: int64(transaction.IDRTotalPrice),
			},
			CustomerDetails: &midtrans.CustomerDetails{
				FName: user.UserProfile.FirstName + " " + user.UserProfile.LastName,
				Email: user.UserProfile.Email,
			},
			CustomExpiry: &coreapi.CustomExpiry{ExpiryDuration: 15, Unit: "minute"},
			PaymentType:  coreapi.PaymentTypeQris,
			Qris:         &coreapi.QrisDetails{},
		}

		chargeResp, err := c.ChargeTransaction(chargeReq)
		if err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "Failed to connect Payment Gateway",
				"id": "Gagal menghubungkan ke Payment Gateway",
			}
			ErrorResponse(w, http.StatusInternalServerError, messages, lang)
			return
		}
		for _, action := range chargeResp.Actions {
			if action.Name == "generate-qr-code" {
				url = action.URL
			}
		}
		transaction.Acquirer = chargeResp.Acquirer
		timeLimitAt := time.Now().Add(15 * time.Minute)
		transaction.TimeLimitAt = &timeLimitAt
	}

	if url != "" {
		transaction.Status = models.TransactionStatusPending
		transaction.PaymentURL = url
		transaction.MidtransStatus = models.MidtransTransactionStatusPending
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

	SuccessResponse(w, http.StatusCreated, "Connection to payment gateway successfully", ConvertToTransactionResponse(transaction))
}
