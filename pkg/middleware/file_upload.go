package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iv_project/dto"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const UploadsKey MiddlewareKey = "uploadedFiles"
const UploadKey MiddlewareKey = "uploadedFile"

func InvitationImagesUploader(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ErrorResult{StatusCode: http.StatusBadRequest, Message: "Gagal parse multipart form"})
			return
		}

		fields := []string{
			"cover_image",
			"bride_image",
			"groom_image",
			"image_1", "image_2", "image_3", "image_4", "image_5", "image_6",
			"image_7", "image_8", "image_9", "image_10", "image_11", "image_12",
		}

		uploadedFiles := make(map[string]string)

		brideNickname := r.FormValue("bride_nickname")
		groomNickname := r.FormValue("groom_nickname")

		for _, fieldName := range fields {
			file, header, err := r.FormFile(fieldName)
			if err != nil {
				if err == http.ErrMissingFile {
					continue
				}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(dto.ErrorResult{StatusCode: http.StatusBadRequest, Message: "Gagal mengambil file"})
				return
			}

			filePath, err := saveImage(file, header.Filename, brideNickname, groomNickname)
			file.Close()

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(dto.ErrorResult{StatusCode: http.StatusBadRequest, Message: err.Error()})
				return
			}

			uploadedFiles[fieldName] = filePath
		}

		ctx := context.WithValue(r.Context(), UploadsKey, uploadedFiles)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func PaymentProofImageUploader(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Batasi ukuran maksimal file yang bisa diunggah (misalnya 10MB)
		r.ParseMultipartForm(10 << 20)

		file, header, err := r.FormFile("image")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ErrorResult{StatusCode: http.StatusBadRequest, Message: "Gagal mengambil file"})
			return
		}
		defer file.Close()

		filePath, err := saveImage(file, header.Filename, r.FormValue("bride_nickname"), r.FormValue("groom_nickname"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ErrorResult{StatusCode: http.StatusBadRequest, Message: err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), UploadsKey, filePath)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func saveImage(src io.Reader, filename, nicknameBride, nicknameGroom string) (string, error) {
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("gagal membaca file: %w", err)
	}

	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("gagal membuat folder upload: %w", err)
		}
	}

	ext := filepath.Ext(filename)

	safeBride := strings.ReplaceAll(strings.ToLower(nicknameBride), " ", "_")
	safeGroom := strings.ReplaceAll(strings.ToLower(nicknameGroom), " ", "_")

	newFilename := fmt.Sprintf("%s_%s_%d%s", safeBride, safeGroom, time.Now().UnixNano(), ext)

	filePath := filepath.Join(uploadDir, newFilename)
	err = os.WriteFile(filePath, fileBytes, 0644)
	if err != nil {
		return "", fmt.Errorf("gagal menyimpan file")
	}

	return filePath, nil
}
