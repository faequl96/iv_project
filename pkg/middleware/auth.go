package middleware

import (
	"context"
	"encoding/json"
	dto "iv_project/dto"
	jwtToken "iv_project/pkg/jwt"
	"net/http"
	"strings"
)

const authKey middlewareKey = "userInfo"

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusUnauthorized, Message: "unauthorized"})
			return
		}

		claims, err := jwtToken.DecodeToken(strings.Split(token, " ")[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusUnauthorized, Message: "unauthorized"})
			return
		}

		ctx := context.WithValue(r.Context(), authKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
