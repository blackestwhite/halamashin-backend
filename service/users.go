package service

import (
	"app/db"
	"app/entity"
	"app/utils"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService struct{}

func (u *UserService) Create(instance entity.User) (insertionID primitive.ObjectID, err error) {
	instance.ID = primitive.NewObjectID()
	instance.Password = utils.HashString(instance.Password)

	res, err := db.Client.Database("hala").Collection("users").InsertOne(context.TODO(), instance)
	return res.InsertedID.(primitive.ObjectID), err
}

func (u *UserService) GetByPhoneNumber(phoneNumber string) (user entity.User, err error) {
	filter := bson.M{
		"phone_number": phoneNumber,
	}

	res := db.Client.Database("hala").Collection("users").FindOne(context.TODO(), filter)
	if res.Err() != nil {
		return user, res.Err()
	}

	err = res.Decode(&user)
	return
}

func (u *UserService) Get(userID string) (user entity.User, err error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return user, err
	}

	res := db.Client.Database("hala").Collection("users").FindOne(context.TODO(), bson.M{
		"_id": objectID,
	})
	if res.Err() != nil {
		return user, res.Err()
	}

	err = res.Decode(&user)
	return
}

func (u *UserService) Update(user entity.User, update interface{}) (err error) {
	_, err = db.Client.Database("hala").Collection("users").UpdateOne(context.TODO(), bson.M{
		"_id": user.ID,
	}, bson.M{
		"$set": update,
	})
	return
}

func (u *UserService) Delete(userID string) (err error) {
	user, err := u.Get(userID)
	if err != nil {
		return err
	}

	_, err = db.Client.Database("hala").Collection("users").DeleteOne(context.TODO(), bson.M{
		"_id": user.ID,
	})
	return
}

func (u *UserService) List(limit, skip int64) (users []entity.User, err error) {
	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSkip(skip)

	cur, err := db.Client.Database("hala").Collection("users").Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var user entity.User
		decodeErr := cur.Decode(&user)
		if decodeErr == nil {
			users = append(users, user)
		} else {
			log.Println("error decoding user", err)
		}
	}

	return
}
