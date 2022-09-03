// Copyright (c) 2022 Braden Nicholson

package jwt

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

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

func AuthToken(token string) (string, error) {
	content, err := jwtauth.VerifyToken(tokenAuth, token)
	if err != nil {
		return "", err
	}

	val, ok := content.Get("id")
	if !ok {
		return "", fmt.Errorf("malformed jwt... This is a serious concern (security or ACAP?)")
	}

	s := val.(string)

	return s, nil
}

func VerifyToken() func(http.Handler) http.Handler {
	return jwtauth.Verifier(tokenAuth)
}

func SignUUID(id string) (string, error) {
	claims := map[string]any{}
	claims["id"] = id
	_, s, err := tokenAuth.Encode(claims)
	if err != nil {
		return s, err
	}
	return s, nil
}
