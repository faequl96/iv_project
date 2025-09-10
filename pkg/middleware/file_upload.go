package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"iv_project/dto"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
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

		cld, err := getCloudinary()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ErrorResult{StatusCode: http.StatusBadRequest, Message: err.Error()})
			return
		}

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
			defer file.Close()

			uploadResp, err := cld.Upload.Upload(
				context.Background(),
				file,
				uploader.UploadParams{
					PublicID: fmt.Sprintf("%s_%s_%s", brideNickname, groomNickname, header.Filename),
					Folder:   "invitations",
				},
			)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(dto.ErrorResult{StatusCode: http.StatusInternalServerError, Message: "Gagal upload ke Cloudinary"})
				return
			}

			uploadedFiles[fieldName] = uploadResp.SecureURL
		}

		ctx := context.WithValue(r.Context(), UploadsKey, uploadedFiles)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getCloudinary() (*cloudinary.Cloudinary, error) {
	cloudName := os.Getenv("CLOUD_NAME")
	apiKey := os.Getenv("CLOUD_API_KEY")
	apiSecret := os.Getenv("CLOUD_API_SECRET")

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("gagal konek ke Cloudinary")
	}
	return cld, nil
}

// func saveImage(src io.Reader, filename, nicknameBride, nicknameGroom string) (string, error) {
// 	fileBytes, err := io.ReadAll(src)
// 	if err != nil {
// 		return "", fmt.Errorf("gagal membaca file: %w", err)
// 	}

// 	uploadDir := "uploads"
// 	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
// 		err := os.Mkdir(uploadDir, os.ModePerm)
// 		if err != nil {
// 			return "", fmt.Errorf("gagal membuat folder upload: %w", err)
// 		}
// 	}

// 	ext := filepath.Ext(filename)

// 	safeBride := strings.ReplaceAll(strings.ToLower(nicknameBride), " ", "_")
// 	safeGroom := strings.ReplaceAll(strings.ToLower(nicknameGroom), " ", "_")

// 	newFilename := fmt.Sprintf("%s_%s_%d%s", safeBride, safeGroom, time.Now().UnixNano(), ext)

// 	filePath := filepath.Join(uploadDir, newFilename)
// 	err = os.WriteFile(filePath, fileBytes, 0644)
// 	if err != nil {
// 		return "", fmt.Errorf("gagal menyimpan file")
// 	}

// 	return filePath, nil
// }
