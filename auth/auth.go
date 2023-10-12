package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"github.com/farismfirdaus/plant-nursery/entity"
	apperr "github.com/farismfirdaus/plant-nursery/errors"
)

type Auth interface {
	Sign(context.Context, *entity.Customer) (string, error)
	Verify(context.Context, string) (int, error)
}

type CustomerAuth struct {
	privateKey []byte
	publicKey  []byte
}

func NewCustomerAuth(privateKey []byte, publicKey []byte) (*CustomerAuth, error) {
	return &CustomerAuth{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, &privateKey.PublicKey, nil
}

func (c *CustomerAuth) Sign(ctx context.Context, customer *entity.Customer) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(c.privateKey)
	if err != nil {
		return "", fmt.Errorf("error while parsing private key: %w", err)
	}

	now := time.Now().UTC()
	ttl := 15 * time.Minute

	claims := make(jwt.MapClaims)
	claims["dat"] = customer.ID         // customer data
	claims["exp"] = now.Add(ttl).Unix() // The expiration time after which the token must be disregarded.
	claims["iat"] = now.Unix()          // The time at which the token was issued.
	claims["nbf"] = now.Unix()          // The time before which the token must be disregarded.

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("error while signing token: %w", err)
	}

	return token, nil
}

func (c *CustomerAuth) Verify(ctx context.Context, token string) (int, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(c.publicKey)
	if err != nil {
		log.Err(err).Msgf("error parsing public key")
		return 0, apperr.Unauthorized
	}

	tok, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			log.Err(err).Msgf("unexpected method: %s", t.Header["alg"])
			return nil, apperr.Unauthorized
		}

		return key, nil
	})
	if err != nil {
		log.Err(err).Msgf("error parsing jwt")
		return 0, apperr.Unauthorized
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return 0, apperr.Unauthorized
	}

	return int(claims["dat"].(float64)), nil
}
