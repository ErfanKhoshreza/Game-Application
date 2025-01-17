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
