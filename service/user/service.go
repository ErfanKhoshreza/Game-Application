package user

import (
	"Game-Application/entity"
	"Game-Application/pkg/phonenumber"
	"Game-Application/repository/mongo"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (entity.User, error)
	Login(params mongo.LoginParams) (entity.User, error)
	GetUserByID(userID string) (entity.User, error)
}
type AuthGenerator interface {
	createAccessToken(user entity.User) (string, error)
	createRefreshToken(user entity.User) (string, error)
	parseToken(accessToken string) (*Claims, error)
}
type Service struct {
	repo Repository
	auth AuthGenerator
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
	token, TErr := s.auth.createAccessToken(user)
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
	claim, PAError := s.auth.parseToken(token)
	if PAError != nil {
		return "", PAError
	}
	return claim.UserID, nil

}
