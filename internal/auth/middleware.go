package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"api-with-graphql/internal/pkg/jwt"
	"api-with-graphql/internal/users"
)

type contextKey struct {
	name string
}

var userCtxKey = &contextKey{"user"}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			tokenStr := header
			username, err := jwt.ParseToken(tokenStr)

			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			fmt.Println("username", username)

			user := users.User{Username: username}
			id, err := users.GetUserIdByUserName(username)

			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			user.Id = strconv.Itoa(id)
			ctx := context.WithValue(r.Context(), userCtxKey, &user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) *users.User {
	raw, _ := ctx.Value(userCtxKey).(*users.User)
	return raw
}
