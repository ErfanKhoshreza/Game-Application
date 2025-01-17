package user

import (
	"Game-Application/entity"
	"Game-Application/pkg/phonenumber"
	"errors"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (entity.User, error)
}
type Service struct {
	repo Repository
}
type RegisterRequest struct {
	Name        string
	PhoneNumber string
}
type RegisterRespond struct {
	User entity.User
}

func New(repo Repository) Service {
	return Service{repo: repo}
}
func (s Service) Register(req RegisterRequest) (RegisterRespond, error) {
	//		Validate Phone Number
	// TODO : We Should Verify Phone Number by Verfication Code
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
	}
	registeredUser, err := s.repo.Register(user)
	if err != nil {
		return RegisterRespond{}, err
	}
	//	Retuurn creaTED User
	return RegisterRespond{User: registeredUser}, nil
}
