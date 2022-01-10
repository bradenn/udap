// Copyright (c) 2022 Braden Nicholson

package auth

import (
	"fmt"
	"github.com/go-chi/jwtauth/v5"
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

func VerifyToken() func(http.Handler) http.Handler {
	return jwtauth.Verifier(tokenAuth)
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
