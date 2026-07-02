package middleware

import (
	"context"
	"net/http"
	"strings"

	"anuario/interno/auth"
)

type contextKey string

const claimsContextKey contextKey = "claims"

func AutenticacaoJWT(manager *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			parts := strings.Fields(header)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				writeAuthError(w, http.StatusUnauthorized, "token bearer ausente")
				return
			}

			claims, err := manager.Validate(parts[1])
			if err != nil {
				writeAuthError(w, http.StatusUnauthorized, "token invalido")
				return
			}

			ctx := context.WithValue(r.Context(), claimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AutorizarRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsDoContexto(r.Context())
			if !ok || claims.Role != role {
				writeAuthError(w, http.StatusForbidden, "acesso negado")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func ClaimsDoContexto(ctx context.Context) (auth.Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(auth.Claims)
	return claims, ok
}

func writeAuthError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"erro":"` + message + `"}`))
}
