package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/KinyaElGrande/GraphQL-example/internal/pkg/jwt"
	"github.com/KinyaElGrande/GraphQL-example/internal/users"
)

var UserCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			//Allow unauthenticated users
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			//validate jwt token
			tokenStr := header
			username, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			//creates user and check if he exists in the Db
			user := users.User{Username: username}
			id, err := users.GetUserIdByUsername(username)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			user.ID = strconv.Itoa(id)
			ctx := context.WithValue(r.Context(), UserCtxKey, &user)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *users.User {
	raw, _ := ctx.Value(UserCtxKey).(*users.User)
	return raw
}
