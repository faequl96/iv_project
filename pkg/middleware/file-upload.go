package middleware

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Context key untuk menyimpan hasil upload
type contextKey string

const UploadedFilesKey contextKey = "uploadedFiles"

// Kompres gambar agar ukurannya di bawah 300 KB tanpa package tambahan
func compressImage(src io.Reader) ([]byte, error) {
	// Decode gambar
	img, err := jpeg.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("gagal decode gambar: %w", err)
	}

	// Kompres dengan kualitas awal 75
	var buf bytes.Buffer
	quality := 75

	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, fmt.Errorf("gagal mengompres gambar: %w", err)
	}

	// Jika masih lebih dari 300 KB, turunkan kualitas
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

// UploaderMiddleware menangani upload dan kompresi gambar sebelum request mencapai handler
func InvitationImagesUploader(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Batasi ukuran maksimal file yang bisa diunggah (misalnya 10MB)
		r.ParseMultipartForm(10 << 20)

		// Nama field gambar yang akan diunggah
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
				http.Error(w, fmt.Sprintf("Gagal mengambil file %s: %v", fieldName, err), http.StatusBadRequest)
				return
			}
			defer file.Close()

			// Kompres gambar
			compressedData, err := compressImage(file)
			if err != nil {
				http.Error(w, fmt.Sprintf("Gagal mengompres file %s: %v", fieldName, err), http.StatusInternalServerError)
				return
			}

			// Simpan file yang sudah dikompres
			filePath := filepath.Join("uploads", header.Filename)
			err = os.WriteFile(filePath, compressedData, 0644)
			if err != nil {
				http.Error(w, "Gagal menyimpan file", http.StatusInternalServerError)
				return
			}

			// Simpan path gambar yang telah dikompres
			uploadedFiles[fieldName] = filePath
		}

		// Tambahkan hasil upload ke dalam context
		ctx := context.WithValue(r.Context(), UploadedFilesKey, uploadedFiles)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
