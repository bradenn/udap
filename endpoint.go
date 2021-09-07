package main

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
)

type Endpoint struct {
	Persistent
	Name  string `json:"name" gorm:"unique"`
	Token string `json:"token"`

	Groups []Group `json:"groups" gorm:"many2many:endpointGroup;"`
}

type tokenClaims struct {
	*jwt.RegisteredClaims
	Token string
	Id    string
}

func (e *Endpoint) BeforeCreate(tx *gorm.DB) error {

	return nil
}

func init() {
	claims := tokenClaims{
		&jwt.RegisteredClaims{
			Issuer:    "",
			Subject:   "",
			Audience:  nil,
			ExpiresAt: nil,
			NotBefore: nil,
			IssuedAt:  nil,
			ID:        "",
		},
		"Test",
		"Crest",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	keyFile, err := ioutil.ReadFile("jwtRS512.pem")
	if err != nil {
		fmt.Println(err)
	}

	signedString, err := token.SignedString(keyFile)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(signedString)

	// pubFile, err := ioutil.ReadFile("jwtRS512.pub")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	token, err = jwt.ParseWithClaims(signedString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return keyFile, nil
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(token, "Hmm")

}
