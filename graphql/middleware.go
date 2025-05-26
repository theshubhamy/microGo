package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/theshubhamy/microGo/services/account"
)

type contextKey string

const (
	ctxKeyIP        contextKey = "ip"
	ctxKeyUserAgent contextKey = "user-agent"
	UserIDKey       contextKey = "userID"
)

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

func AuthMiddleware(redisClient *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			opName := extractOperationName(r)
			log.Println(opName)
			if opName == "loginAccount" || opName == "createAccount" {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			token := extractBearerToken(authHeader)
			if token == "" {
				http.Error(w, `Missing token-${opName}`, http.StatusUnauthorized)
				return
			}

			claims, err := account.VerifyJWT(token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			sessionKey := fmt.Sprintf("user-sessions:%s", claims.UserID)
			exists, err := redisClient.Exists(r.Context(), sessionKey).Result()
			if err != nil || exists == 0 {
				http.Error(w, "Session not found", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractBearerToken(header string) string {
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer ")
	}
	return ""
}

func extractOperationName(r *http.Request) string {
	type gqlReq struct {
		OperationName string `json:"operationName"`
	}

	var body gqlReq

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		return ""
	}

	// Restore body so next handler can read it again
	r.Body = io.NopCloser(bytes.NewBuffer(buf))

	if err := json.Unmarshal(buf, &body); err != nil {
		return ""
	}

	return body.OperationName
}
