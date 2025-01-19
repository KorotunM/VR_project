package database

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Game структура для описания игры (определите поля согласно вашей модели)
type Games struct {
	ID    primitive.ObjectID `bson:"_id"`
	Name  string             `bson:"name"`
	Price float64            `bson:"price"`
}

func GetGameByID(id string) (Games, error) {
	collection := MongoClient.Database("Vr").Collection("Games")

	var result Games

	// Преобразуем строковый ID в ObjectID
	gameID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Games{}, errors.New("invalid game ID format")
	}

	// Фильтр поиска
	filter := bson.M{"_id": gameID}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Выполняем поиск в коллекции
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return Games{}, errors.New("game not found")
	}

	log.Println("Найденная игра:", result.Name)
	return result, nil
}
