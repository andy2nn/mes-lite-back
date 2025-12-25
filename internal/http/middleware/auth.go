package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey   contextKey = "userID"
	UserRoleKey contextKey = "userRole"
)

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			rawHeader := r.Header.Get("Authorization")
			slog.Debug("AuthMiddleware: Authorization header received",
				slog.String("header", rawHeader),
				slog.String("path", r.URL.Path),
			)

			if rawHeader == "" {
				slog.Warn("AuthMiddleware: missing Authorization header",
					slog.String("path", r.URL.Path))
				http.Error(w, "no token", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(rawHeader, "Bearer ") {
				slog.Warn("AuthMiddleware: invalid Authorization format",
					slog.String("header", rawHeader))
				http.Error(w, "invalid token format", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(rawHeader, "Bearer ")
			slog.Debug("AuthMiddleware: extracted token",
				slog.String("token", tokenString))

			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
				return []byte(secret), nil
			})

			if err != nil {
				slog.Warn("AuthMiddleware: JWT parse failed",
					slog.String("token", tokenString),
					slog.Any("error", err))
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				slog.Warn("AuthMiddleware: token invalid",
					slog.String("token", tokenString))
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			slog.Info("AuthMiddleware: token accepted",
				slog.Any("claims", claims),
			)

			// Put data into context
			ctx := context.WithValue(r.Context(), UserIDKey, claims["id"])
			ctx = context.WithValue(ctx, UserRoleKey, claims["role"])

			slog.Debug("AuthMiddleware: user context applied",
				slog.Any("user_id", claims["id"]),
				slog.Any("role", claims["role"]),
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
