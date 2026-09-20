package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/yourname/qr-ordering-system/internal/models"
	"github.com/yourname/qr-ordering-system/internal/services"
	"github.com/yourname/qr-ordering-system/internal/utils"
)

type contextKey string

const ClaimsContextKey contextKey = "claims"

func RequireAuth(authService *services.AuthService, allowedRoles ...models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				utils.RespondError(w, http.StatusUnauthorized, "missing or invalid authorization header")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := authService.ValidateToken(tokenString)
			if err != nil {
				utils.RespondError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			if len(allowedRoles) > 0 {
				allowed := false
				for _, role := range allowedRoles {
					if claims.Role == role {
						allowed = true
						break
					}
				}
				if !allowed {
					utils.RespondError(w, http.StatusForbidden, "insufficient permissions")
					return
				}
			}

			ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
