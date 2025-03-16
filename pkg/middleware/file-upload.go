package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	invitation_dto "iv_project/dto/invitation"
	"net/http"
	"os"
)

// UploaderMiddleware menangani upload dan kompresi gambar sebelum request mencapai handler
func InvitationImagesUploader(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20) // 10MB max upload size
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Ambil JSON dari form-data
		jsonData := r.FormValue("data")

		// Unmarshal JSON ke struct
		invitation := new(invitation_dto.InvitationRequest)
		err = json.Unmarshal([]byte(jsonData), &invitation)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Ambil files dari form-data
		files := r.MultipartForm.File["images"]
		if len(files) != len(invitation.InvitationData.Gallery) {
			http.Error(w, "Jumlah gambar tidak sesuai dengan JSON", http.StatusBadRequest)
			return
		}

		// Loop untuk menyimpan gambar dengan nama dari JSON
		for i, file := range files {
			src, err := file.Open()
			if err != nil {
				http.Error(w, "Error opening file", http.StatusInternalServerError)
				return
			}
			defer src.Close()

			// Ambil nama dari image_url di JSON
			filename := "uploads/" + invitation.InvitationData.Gallery[i].ImageURL

			// Kompres dan Simpan file
			err = compressAndSaveImage(src, filename, 300*1024) // 300 KB max
			if err != nil {
				http.Error(w, "Error compressing file", http.StatusInternalServerError)
				return
			}
		}

		next.ServeHTTP(w, r)
	}
}

// UploaderMiddleware menangani upload dan kompresi gambar sebelum request mencapai handler
func InvitationDataImagesUploader(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20) // 10MB max upload size
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Ambil JSON dari form-data
		jsonData := r.FormValue("data")

		// Unmarshal JSON ke struct
		invitationData := new(invitation_dto.InvitationDataRequest)
		err = json.Unmarshal([]byte(jsonData), &invitationData)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Ambil files dari form-data
		files := r.MultipartForm.File["images"]
		if len(files) != len(invitationData.Gallery) {
			http.Error(w, "Jumlah gambar tidak sesuai dengan JSON", http.StatusBadRequest)
			return
		}

		// Loop untuk menyimpan gambar dengan nama dari JSON
		for i, file := range files {
			src, err := file.Open()
			if err != nil {
				http.Error(w, "Error opening file", http.StatusInternalServerError)
				return
			}
			defer src.Close()

			// Ambil nama dari image_url di JSON
			filename := "uploads/" + invitationData.Gallery[i].ImageURL

			// Kompres dan Simpan file
			err = compressAndSaveImage(src, filename, 300*1024) // 300 KB max
			if err != nil {
				http.Error(w, "Error compressing file", http.StatusInternalServerError)
				return
			}
		}

		next.ServeHTTP(w, r)
	}
}

// Fungsi untuk mengompresi gambar tanpa package tambahan
func compressAndSaveImage(src io.Reader, outputPath string, maxSize int) error {
	// Decode gambar dari file upload
	img, format, err := image.Decode(src)
	if err != nil {
		return err
	}

	// Periksa apakah format didukung (hanya JPEG untuk sekarang)
	if format != "jpeg" && format != "jpg" {
		return fmt.Errorf("unsupported format: %s", format)
	}

	// Coba menyimpan dengan kualitas bertahap hingga di bawah 300 KB
	var buf bytes.Buffer
	quality := 80 // Mulai dari kualitas 80%
	for {
		buf.Reset() // Reset buffer setiap iterasi
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
		if err != nil {
			return err
		}

		// Jika ukuran sudah sesuai atau kualitas sudah sangat rendah, simpan file
		if buf.Len() <= maxSize || quality <= 30 {
			break
		}
		quality -= 5 // Kurangi kualitas
	}

	// Simpan ke file dengan nama yang sudah ditentukan
	dst, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, &buf)
	return err
}
