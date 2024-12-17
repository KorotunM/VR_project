package database

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTariffIDByName(name string) (string, error) {
	collection := MongoClient.Database("Vr").Collection("Tariffs")

	var result struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	filter := bson.M{"name": name}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return "", errors.New("tariff not found")
	}
	log.Println("найденный Тариф:", result.ID.Hex())
	return result.ID.Hex(), nil
}
