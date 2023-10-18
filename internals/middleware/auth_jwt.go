package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
)

type key int

const (
	ClaimJWTKey key = iota
)

var jwtKey string

func SetTokenizeConfig(jwtK string) {
	jwtKey = jwtK
}

type MapClaimResponse struct {
	Id    ulid.ULID `json:"id"`
	Email string    `json:"email"`
	jwt.MapClaims
}

func AuthJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			idToken := strings.TrimSpace(strings.Replace(authorization, "Bearer", "", 1))
			if idToken == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			claims := &MapClaimResponse{}
			token, err := jwt.ParseWithClaims(
				idToken,
				claims,
				func(t *jwt.Token) (interface{}, error) {
					return []byte(jwtKey), nil
				},
			)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if !token.Valid {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			claims = token.Claims.(*MapClaimResponse)
			ctx := r.Context()
			req := r.WithContext(
				context.WithValue(
					ctx,
					ClaimJWTKey,
					claims,
				),
			)
			next.ServeHTTP(w, req)
		},
	)
}
