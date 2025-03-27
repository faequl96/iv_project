package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	dto "iv_project/dto"
	jwtToken "iv_project/pkg/jwt"
	"net/http"
	"strings"
)

const UserIdKey MiddlewareKey = "userID"
const RoleKey MiddlewareKey = "role"

func Auth(jwtServices jwtToken.JWTServices, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(dto.ErrorResult{StatusCode: http.StatusUnauthorized, Message: "Authorization header is required"})
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(dto.ErrorResult{StatusCode: http.StatusUnauthorized, Message: "Invalid token format"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		claims, err := jwtServices.DecodeToken(tokenString)
		if err != nil || claims == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(dto.ErrorResult{StatusCode: http.StatusUnauthorized, Message: "Invalid or expired token"})
			return
		}

		idValue, exists := claims["id"]
		if !exists {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(dto.ErrorResult{StatusCode: http.StatusUnauthorized, Message: "Invalid token claims: missing userID"})
			return
		}
		userID, ok := idValue.(string)
		if !ok {
			userID = fmt.Sprintf("%v", idValue)
		}

		roleValue, exists := claims["role"]
		if !exists {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(dto.ErrorResult{StatusCode: http.StatusUnauthorized, Message: "Invalid token claims: missing role"})
			return
		}
		role, ok := roleValue.(string)
		if !ok {
			role = fmt.Sprintf("%v", roleValue)
		}

		ctx := context.WithValue(r.Context(), UserIdKey, userID)
		ctx = context.WithValue(ctx, RoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
