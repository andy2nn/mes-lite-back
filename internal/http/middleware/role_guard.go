package middleware

import (
	"net/http"
	"slices"
)

func RoleGuard(allowed ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roleVal := r.Context().Value(UserRoleKey)
			if roleVal == nil {
				http.Error(w, "no role", http.StatusForbidden)
				return
			}

			role, ok := roleVal.(string)
			if !ok {
				http.Error(w, "invalid role type", http.StatusForbidden)
				return
			}

			if slices.Contains(allowed, role) {
				next.ServeHTTP(w, r)
				return
			}

			http.Error(w, "forbidden", http.StatusForbidden)
		})
	}
}
