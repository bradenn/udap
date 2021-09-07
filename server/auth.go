package server

import (
	"context"
	"github.com/go-chi/jwtauth/v5"
	"os"
)

var tokenAuth *jwtauth.JWTAuth

func init() {

	privateKey := os.Getenv("private")

	tokenAuth = jwtauth.New("HS512", []byte(privateKey), nil)

}

func CheckJwt(jwt string) interface{} {

	s, err := tokenAuth.Decode(jwt)

	if err != nil {
		return ""
	}
	p, err := s.AsMap(context.Background())
	return p
}

func WriteJwt(claims map[string]interface{}) string {

	_, s, err := tokenAuth.Encode(claims)
	if err != nil {
		return ""
	}

	return s
}
