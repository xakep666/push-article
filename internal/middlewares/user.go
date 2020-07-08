package middlewares

import (
	"context"
	"net/http"
	"strconv"
)

func Authorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDCookie, err := r.Cookie("example-push-user")
		if err != nil {
			http.Error(w, "cookie not found", http.StatusUnauthorized)
			return
		}

		userID, err := strconv.ParseUint(userIDCookie.Value, 10, 64)
		if err != nil {
			http.Error(w, "bad user id: "+err.Error(), http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r.WithContext(ContextWithUser(r.Context(), userID)))
	})
}

type userIDContextKeyT struct{}

var userIDContextKey userIDContextKeyT

func ContextWithUser(ctx context.Context, userID uint64) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func UserFromContext(ctx context.Context) uint64 {
	u, _ := ctx.Value(userIDContextKey).(uint64)
	return u
}
