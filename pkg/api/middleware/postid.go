package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

// PostID validates postID param and sets it as a context value
func PostID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		postID, err := strconv.Atoi(chi.URLParam(r, "postID"))
		if err != nil {
			http.Error(w, "Invalid PostID Param", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "post_id", postID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
