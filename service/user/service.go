package user

import (
	"Game-Application/entity"
	"Game-Application/pkg/phonenumber"
	"Game-Application/repository/mongo"
	"errors"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (entity.User, error)
	Login(params mongo.LoginParams) (entity.User, error)
	GetUserByID(userID uint) (entity.User, error)
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
	User entity.User
}
type ProfileRequest struct {
	UserID uint
}
type ProfileResponse struct {
	Name string `json:"name"`
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
	return LoginRespond{User: user}, nil

}

func (s Service) getProfile(req ProfileRequest) (ProfileResponse, error) {
	// Get User By Id
	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return ProfileResponse{}, err
	}
	//	Retun IT
	return ProfileResponse{Name: user.Name}, nil
}
