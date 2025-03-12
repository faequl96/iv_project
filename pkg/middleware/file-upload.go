package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

type Result struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type dataFileContextKey string

const dataFileKey dataFileContextKey = "dataFile"

func UploadFile(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("image")

		// condition not uploading file when update
		if err != nil && r.Method == "PATCH" {
			ctx := context.WithValue(r.Context(), dataFileKey, "false")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode("Error Retrieving the File")
			return
		}

		defer file.Close()

		const MAX_UPLOAD_SIZE = 10 << 20

		r.ParseMultipartForm(MAX_UPLOAD_SIZE)
		if r.ContentLength > MAX_UPLOAD_SIZE {
			w.WriteHeader(http.StatusBadRequest)
			response := Result{Code: http.StatusBadRequest, Message: "Max size in 10mb"}
			json.NewEncoder(w).Encode(response)
			return
		}

		tempFile, err := os.CreateTemp("uploads", "image-*.png")
		if err != nil {
			fmt.Println(err)
			fmt.Println("path upload error")
			json.NewEncoder(w).Encode(err)
			return
		}
		defer tempFile.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}

		tempFile.Write(fileBytes)

		data := tempFile.Name()
		// filename := data[8:]

		ctx := context.WithValue(r.Context(), dataFileKey, data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Context key untuk menyimpan hasil upload
type contextKey string

const UploadedFilesKey contextKey = "uploadedFiles"

// Kompres gambar agar ukurannya di bawah 300 KB
func compressImage(src io.Reader) ([]byte, error) {
	// Decode gambar
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("gagal decode gambar: %w", err)
	}

	// Resize gambar (misalnya maksimal lebar 1024)
	newImage := resize.Resize(1024, 0, img, resize.Lanczos3)

	// Kompres dengan kualitas 75
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, newImage, &jpeg.Options{Quality: 75})
	if err != nil {
		return nil, fmt.Errorf("gagal mengompres gambar: %w", err)
	}

	// Jika masih lebih dari 300 KB, turunkan kualitas
	quality := 75
	for buf.Len() > 300*1024 && quality > 10 {
		buf.Reset()
		quality -= 5
		err = jpeg.Encode(&buf, newImage, &jpeg.Options{Quality: quality})
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
		fields := []string{"gallery_image_url_1", "gallery_image_url_2", "gallery_image_url_3",
			"gallery_image_url_4", "gallery_image_url_5", "gallery_image_url_6",
			"gallery_image_url_7", "gallery_image_url_8", "gallery_image_url_9",
			"gallery_image_url_10", "gallery_image_url_11", "gallery_image_url_12"}
		uploadedFiles := make(map[string]string)

		for _, fieldName := range fields {
			file, header, err := r.FormFile(fieldName)
			if err != nil {
				if err == http.ErrMissingFile {
					continue
				}
				http.Error(w, fmt.Sprintf("Gagal mengambil file %s: %v", fieldName, err), http.StatusBadRequest)
				return
			}
			defer file.Close()

			// Baca file untuk kompresi
			var buf bytes.Buffer
			tee := io.TeeReader(file, &buf)

			// Kompres gambar
			compressedData, err := compressImage(tee)
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
