package mongo

import (
	"Game-Application/entity"
	"Game-Application/pkg/password"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type LoginParams struct {
	PhoneNumber string
	Password    string
}

func (d DB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user := d.Database.Collection("gameUser").FindOne(ctx, bson.M{"phoneNumber": phoneNumber})
	if errors.Is(user.Err(), mongo.ErrNoDocuments) {
		return true, nil
	} else if user.Err() != nil {
		return false, user.Err()
	}
	return true, nil
}

func (d DB) Register(u entity.User) (entity.User, error) {
	hashedPassword, err := password.HashPassword(u.Password)
	if err != nil {
		return entity.User{}, err
	}
	user := entity.User{
		PhoneNumber: u.PhoneNumber,
		Name:        u.Name,
		Avatar:      u.Avatar,
		Password:    hashedPassword,
	}
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, CErr := d.Database.Collection("gameUser").InsertOne(ctx, user)
	if CErr != nil {
		return entity.User{}, err
	}
	return user, nil
}
func (d DB) Login(param LoginParams) (entity.User, error) {
	user, err := d.FindUserByPhoneNumber(param.PhoneNumber)
	if err != nil {
		return entity.User{}, err
	}
	if password.CheckPasswordHash(param.Password, user.Password) {
		return user, nil
	}
	return entity.User{}, errors.New("password Did Not Match")
}
func (d DB) GetUserByID(userID uint) (entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user entity.User
	err := d.Database.Collection("gameUser").FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, errors.New("user Not Found")
		}
		return entity.User{}, err // If any other error occurs
	}
	return user, nil
}

func (d DB) FindUserByPhoneNumber(phoneNumber string) (entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user entity.User
	err := d.Database.Collection("gameUser").FindOne(ctx, bson.M{"phonenumber": phoneNumber}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, errors.New("user Not Found")
		}
		return entity.User{}, err // If any other error occurs
	}
	return user, nil
}
