package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Firstname   string             `json:"firstname" bson:"firstname"`
	Lastname    string             `json:"lastname" bson:"lastname"`
	PhoneNumber string             `json:"phone_number" bson:"phone_number"`
	Password    string             `json:"password" bson:"password"`
}
