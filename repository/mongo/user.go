package mongo

import (
	"Game-Application/entity"
	"Game-Application/pkg/password"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

//func (d DB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
//	d.Database.Collection("gameUser").FindOne()
//}

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
	registeredUser, CErr := d.Database.Collection("gameUser").InsertOne(ctx, u)
	if CErr != nil {
		return entity.User{}, err
	}
	if oid, ok := registeredUser.InsertedID.(primitive.ObjectID); ok {
		u.ID = uint(oid.Timestamp().Unix())
	} else {
		return entity.User{}, fmt.Errorf("failed to convert inserted ID to ObjectID")
	}
	return user, nil
}
