package middleware

import "net/http"

func RequireRole(role string) func(http.Handler) http.Handler {

	return RequireRoles(role)
}

// RequireRoles allows any of the given roles.
// Notes:
// - "superadmin" is treated as an alias of "admin" for compatibility.
func RequireRoles(roles ...string) func(http.Handler) http.Handler {

	roleSet := map[string]struct{}{}
	for _, r := range roles {
		roleSet[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userRole, ok := r.Context().Value(RoleKey).(string)

			if !ok || userRole == "" {
				http.Error(w, "Forbidden", 403)
				return
			}

			if userRole == "superadmin" {
				userRole = "admin"
			}

			if _, allowed := roleSet[userRole]; !allowed {
				http.Error(w, "Forbidden", 403)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
