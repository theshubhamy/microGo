package graphql

import (
	"context"
	"net/http"
)

type contextKey string

const (
	ctxKeyIP        contextKey = "ip"
	ctxKeyUserAgent contextKey = "user-agent"
)

// ✅ Middleware to inject IP and User-Agent into context
func InjectRequestMeta(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}
		userAgent := r.Header.Get("User-Agent")

		ctx := context.WithValue(r.Context(), ctxKeyIP, ip)
		ctx = context.WithValue(ctx, ctxKeyUserAgent, userAgent)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
