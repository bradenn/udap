package server

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"net/http"
	"os"
)

var tokenAuth *jwtauth.JWTAuth

func Init() {
	privateKey := os.Getenv("private")
	tokenAuth = jwtauth.New("HS512", []byte(privateKey), nil)
}

func VerifyToken() func(http.Handler) http.Handler {
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
func SignUUID(uuid uuid.UUID) (string, error) {
	claims := Claims{}
	claims["id"] = uuid
	_, s, err := tokenAuth.Encode(claims)
	if err != nil {
		return s, err
	}
	return s, nil
}
