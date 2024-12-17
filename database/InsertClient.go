package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertClient(name, email, phone string) (string, error) {
	collection := MongoClient.Database("Vr").Collection("Clients")

	client := map[string]interface{}{
		"name":         name,
		"email":        email,
		"phone number": phone,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, client)
	if err != nil {
		return "", err
	}

	log.Println("Вставленный пользователь:", result.InsertedID.(primitive.ObjectID).Hex())
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}
