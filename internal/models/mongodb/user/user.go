package user

import (
	"context"

	"github.com/cpustejovsky/mongogo/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Create(collection *mongo.Collection, user models.FormUser) (interface{}, error) {
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	return insertResult.InsertedID, nil
}

func Fetch(collection *mongo.Collection, id string) (models.User, error) {
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{
		"_id": id,
	}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func Update(collection *mongo.Collection, updatedUser map[string]interface{}) (models.User, error) {
	var user models.User
	filter := bson.M{
		"_id": updatedUser["_id"],
	}
	update := bson.M{
		"$set": updatedUser,
	}

	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	result := collection.FindOneAndUpdate(context.TODO(), filter, update, &opt)
	if result.Err() != nil {
		return user, result.Err()
	}
	decodeErr := result.Decode(&user)
	return user, decodeErr
}

func Delete(collection *mongo.Collection, id string) error {
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
