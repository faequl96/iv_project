package middleware

import (
	"context"
	"net/http"
)

const LanguageKey MiddlewareKey = "language"

func Language(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := r.Header.Get("Accept-Language")
		if lang == "" {
			lang = "id"
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, LanguageKey, lang)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
