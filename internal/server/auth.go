// Copyright (c) 2021 Braden Nicholson

package server

import (
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwt"
	"net/http"
	"os"
)

var tokenAuth *jwtauth.JWTAuth

func LoadKeys() {
	privateKey := os.Getenv("private")
	tokenAuth = jwtauth.New("HS512", []byte(privateKey), nil)
}

func AuthToken(token string) (string, error) {
	content, err := jwtauth.VerifyToken(tokenAuth, token)
	if err != nil {
		return "", err
	}

	val, ok := content.Get("id")
	if !ok {
		return "", fmt.Errorf("malformed jwt... This is a serious concern")
	}

	s := val.(string)

	return s, nil
}

func verifyToken() func(http.Handler) http.Handler {
	return jwtauth.Verifier(tokenAuth)
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())
		req, _ := NewRequest(w, r)

		if err != nil {
			req.Reject(err.Error(), http.StatusBadRequest)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			req.Reject(err, http.StatusUnauthorized)
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

type Claims map[string]interface{}

// SignUUID will generate and sign a JWT key with a set of claims. Use wisely.
func SignUUID(id string) (string, error) {
	claims := Claims{}
	claims["id"] = id
	_, s, err := tokenAuth.Encode(claims)
	if err != nil {
		return s, err
	}
	return s, nil
}