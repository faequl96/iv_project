package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
