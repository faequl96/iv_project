package middleware

import (
	"context"
	"encoding/json"
	dto "iv_project/dto"
	jwtToken "iv_project/pkg/jwt"
	"net/http"
	"strings"
)

const userIdKey middlewareKey = "userID"
const roleKey middlewareKey = "role"

func Auth(jwtServices jwtToken.JWTServices, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusUnauthorized, Message: "Authorization header is required"})
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusUnauthorized, Message: "Invalid token format"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		claims, err := jwtServices.DecodeToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusUnauthorized, Message: "Invalid or expired token"})
			return
		}

		userID, ok := claims["id"].(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusUnauthorized, Message: "Invalid token payload"})
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusUnauthorized, Message: "Invalid token role"})
			return
		}

		ctx := context.WithValue(r.Context(), userIdKey, userID)
		ctx = context.WithValue(ctx, roleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
