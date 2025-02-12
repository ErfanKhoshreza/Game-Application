package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Name        string
	PhoneNumber string
	Avatar      string
	Password    string
	ID          primitive.ObjectID `bson:"_id,omitempty"`
}
