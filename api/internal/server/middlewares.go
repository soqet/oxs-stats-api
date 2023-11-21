package server

import (
	"context"
	"crypto/ecdsa"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const (
	tokenKey ctxKey = "token"
)

func AuthMiddleware(pubKey *ecdsa.PublicKey) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			at, err := jwt.ParseWithClaims(
				token,
				&jwt.RegisteredClaims{},
				jwt.Keyfunc(func(_ *jwt.Token) (any, error) { return pubKey, nil }),
				jwt.WithValidMethods([]string{"ES256"}),
				jwt.WithLeeway(time.Second*15),
			)
			if err != nil || !at.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), tokenKey, at.Claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
