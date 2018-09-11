package middleware

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs"
)

// Auth middleware only lets threw requests with a valid session.
// Additionally it sets user_id key in context
func Auth(sm *scs.Manager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			session := sm.Load(r)
			userID, err := session.GetInt("user_id")
			if err != nil || userID < 1 {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
