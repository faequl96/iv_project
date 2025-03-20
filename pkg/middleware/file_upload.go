package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io"
	"iv_project/dto"
	"net/http"
	"os"
	"path/filepath"
)

const UploadsKey MiddlewareKey = "uploadedFiles"

// Kompres gambar agar ukurannya di bawah 300 KB tanpa package tambahan
func compressImage(src io.Reader) ([]byte, error) {
	img, err := jpeg.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("gagal decode gambar: %w", err)
	}

	var buf bytes.Buffer
	quality := 75
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, fmt.Errorf("gagal mengompres gambar: %w", err)
	}

	for buf.Len() > 300*1024 && quality > 10 {
		buf.Reset()
		quality -= 5
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
		if err != nil {
			return nil, fmt.Errorf("gagal mengompres gambar: %w", err)
		}
	}

	return buf.Bytes(), nil
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

			// Kompres gambar
			compressedData, err := compressImage(file)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
				return
			}

			// Simpan file yang sudah dikompres
			filePath := filepath.Join("uploads", header.Filename)
			err = os.WriteFile(filePath, compressedData, 0644)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusBadRequest, Message: "Gagal menyimpan file"})
				return
			}

			uploadedFiles[fieldName] = filePath
		}

		// Tambahkan hasil upload ke dalam context
		ctx := context.WithValue(r.Context(), UploadsKey, uploadedFiles)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
