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
)

const UploadsKey MiddlewareKey = "uploadedFiles"

func saveImage(src io.Reader, filename string) (string, error) {
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

	filePath := filepath.Join(uploadDir, filename)
	err = os.WriteFile(filePath, fileBytes, 0644)
	if err != nil {
		return "", fmt.Errorf("gagal menyimpan file")
	}

	return filePath, nil
}

func InvitationImagesUploader(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Batasi ukuran maksimal file yang bisa diunggah (misalnya 10MB)
		r.ParseMultipartForm(10 << 20)

		fields := []string{"image_url_1", "image_url_2", "image_url_3",
			"image_url_4", "image_url_5", "image_url_6",
			"image_url_7", "image_url_8", "image_url_9",
			"image_url_10", "image_url_11", "image_url_12"}

		uploadedFiles := make(map[string]string)

		for _, fieldName := range fields {
			file, header, err := r.FormFile(fieldName)
			if err != nil {
				if err == http.ErrMissingFile {
					continue // Lewati jika user tidak mengunggah gambar ini
				}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusBadRequest, Message: "Gagal mengambil file"})
				return
			}
			defer file.Close()

			filePath, err := saveImage(file, header.Filename)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
				return
			}

			uploadedFiles[fieldName] = filePath
		}

		ctx := context.WithValue(r.Context(), UploadsKey, uploadedFiles)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
