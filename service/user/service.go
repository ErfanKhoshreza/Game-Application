package user

import (
	"Game-Application/entity"
	"Game-Application/pkg/phonenumber"
	"Game-Application/repository/mongo"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (entity.User, error)
	Login(params mongo.LoginParams) (entity.User, error)
	GetUserByID(userID string) (entity.User, error)
}
type Service struct {
	repo Repository
}
type RegisterRequest struct {
	Name        string
	PhoneNumber string
	Password    string
}
type LoginRequest struct {
	PhoneNumber string
	Password    string
}
type RegisterRespond struct {
	User entity.User
}
type LoginRespond struct {
	User  entity.User
	Token string
}
type ProfileRequest struct {
	UserID string
}
type ProfileResponse struct {
	Name string `json:"name"`
}
type Claims struct {
	UserID string
	jwt.RegisteredClaims
}

func New(repo Repository) Service {
	return Service{repo: repo}
}
func (s Service) Register(req RegisterRequest) (RegisterRespond, error) {

	if phonenumber.IsValid(req.PhoneNumber) {
		return RegisterRespond{}, errors.New("invalid phone")
	}
	// checkUniqueNess Phone number
	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			return RegisterRespond{}, err
		}
		return RegisterRespond{}, errors.New("phone Number is not Unique")
	}
	//	 Validate Name
	if len(req.Name) < 3 {
		return RegisterRespond{}, errors.New("name is too short")
	}

	//	 create New User In  storage
	user := entity.User{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	}
	registeredUser, err := s.repo.Register(user)
	if err != nil {
		return RegisterRespond{}, err
	}
	//	Retuurn creaTED User
	return RegisterRespond{User: registeredUser}, nil
}

func (s Service) Login(req LoginRequest) (LoginRespond, error) {
	if phonenumber.IsValid(req.PhoneNumber) {
		return LoginRespond{}, errors.New("invalid phone")
	}
	params := mongo.LoginParams{
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
	}
	user, err := s.repo.Login(params)
	if err != nil {
		return LoginRespond{}, errors.New(err.Error())
	}
	privateKeyPath := os.Getenv("PRIVATE_KEY_PATH")
	privateKey, PRErr := loadPrivateKey(privateKeyPath)
	if PRErr != nil {
		return LoginRespond{}, errors.New(PRErr.Error())
	}
	token, TErr := createToken(user.ID.Hex(), privateKey)
	if TErr != nil {
		return LoginRespond{}, TErr
	}
	return LoginRespond{User: user, Token: token}, nil

}

func (s Service) GetProfile(req ProfileRequest) (ProfileResponse, error) {
	// Get User By Id
	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return ProfileResponse{}, err
	}
	//	Retun IT
	return ProfileResponse{Name: user.Name}, nil
}

func (s Service) GetUserIDByToken(token string) (string, error) {
	PublicKeyPath := os.Getenv("PUBLIC_KEY_PATH")
	publicKey, PError := loadPublicKey(PublicKeyPath)
	if PError != nil {
		return "", PError
	}
	claim, PAError := parseJWTToken(token, publicKey)
	if PAError != nil {
		return "", PAError
	}
	return claim.UserID, nil

}
func createToken(userID string, privateKey *rsa.PrivateKey) (string, error) {
	t := jwt.New(jwt.SigningMethodRS256) // Use RS256 correctly

	// Assign claims
	t.Claims = &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		UserID: userID,
	}

	// Sign the token with the private key
	return t.SignedString(privateKey)
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
func parseJWTToken(tokenString string, publicKey *rsa.PublicKey) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
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
