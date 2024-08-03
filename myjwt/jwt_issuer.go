package myjwt

import (
	"crypto"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type Issuer struct {
	key crypto.PrivateKey
}

func NewIssuer(pkPath string) (*Issuer, error) {
	kb, err := os.ReadFile(pkPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to read private key file at %s : %w", pkPath, err))
	}
	key, err := jwt.ParseEdPrivateKeyFromPEM(kb)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse private key: %w", err))
	}
	return &Issuer{
		key: key,
	}, nil
}

func (issuer *Issuer) IssueToken(user string, roles []string) (string, error) {
	ct := time.Now()
	token := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, jwt.MapClaims{
		"iss":   "Kebin Test Server",
		"sub":   user,
		"aud":   "api",
		"exp":   ct.Add(time.Minute).Unix(),
		"nbf":   ct.Unix(),
		"iat":   ct.Unix(),
		"user":  user,
		"roles": roles,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(issuer.key)
	if err != nil {
		return "", errors.New(fmt.Sprintf("failed to issue token: %w", err))
	}
	return tokenString, nil
}
