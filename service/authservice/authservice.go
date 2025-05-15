package authservice

import (
	"Game-Application/entity"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type Service struct {
	signKey               string
	AccessExpirationTime  time.Duration
	RefreshExpirationTime time.Duration
}

func (s Service) CreateAccessToken(user entity.User) (string, error) {
	t := jwt.New(jwt.SigningMethodRS256) // Use RS256 correctly

	// Assign claims
	t.Claims = &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		UserID: user.ID.Hex(),
	}
	privateKeyPath := os.Getenv("PRIVATE_KEY_PATH")
	privateKey, PRErr := loadPrivateKey(privateKeyPath)
	if PRErr != nil {
		return "", PRErr
	}
	// Sign the token with the private key
	return t.SignedString(privateKey)
}
func (s Service) CreateRefreshToken(user entity.User) {

}

func (s Service) ParseToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		PublicKeyPath := os.Getenv("PUBLIC_KEY_PATH")
		publicKey, PError := loadPublicKey(PublicKeyPath)
		if PError != nil {
			return "", PError
		}
		// Ensure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Validate the claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
func loadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	parsedKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return parsedKey, nil
}

func loadPublicKey(filename string) (*rsa.PublicKey, error) {
	keyBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	parsedKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return parsedKey, nil
}
