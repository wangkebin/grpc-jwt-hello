package myjwt

import (
	"crypto"
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

type Validator struct {
	pkey crypto.PublicKey
}

func NewValidator(pkPath string) (*Validator, error) {
	kb, err := os.ReadFile(pkPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to generate validator: %w", err))
	}

	key, err := jwt.ParseEdPublicKeyFromPEM(kb)
	if err != nil {
		return nil, fmt.Errorf("unable to parse as ed private key: %w", err)
	}

	return &Validator{
		pkey: key,
	}, nil

}

func (v *Validator) ValidateTkString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return v.pkey, nil
		})
	if err != nil {
		return nil, fmt.Errorf("unable to validate token string: %w", err)
	}
	return token, nil
}
