package middleware

import (
	"context"
	"myapp/tools"
	"net/http"
	"strings"
)

var CtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

type User struct {
	ID int `json:"id"`
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			next.ServeHTTP(w, r)
			return
		}

		authTokens := strings.Split(authToken, " ")
		if authTokens[0] != "Bearer" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		jwtToken, err := tools.TokenValidate(authTokens[1])
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := jwtToken.Claims.(*tools.JwtClaim)
		if !ok || !jwtToken.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), CtxKey, &User{
			ID: claims.ID,
		})

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func AuthContext(ctx context.Context) *User {
	raw, _ := ctx.Value(CtxKey).(*User)
	return raw
}
