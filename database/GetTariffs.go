package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// GetTariffs получает все тарифы из коллекции "Tariffs"
func GetTariffs() ([]Tariff, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := MongoClient.Database("Vr").Collection("Tariffs")

	var tariffs []Tariff
	cursor, err := collection.Find(ctx, bson.M{}) // Пустой фильтр - получаем все документы
	if err != nil {
		log.Println("Ошибка при получении тарифов из MongoDB:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	// Декодируем документы в слайс тарифов
	if err := cursor.All(ctx, &tariffs); err != nil {
		log.Println("Ошибка при декодировании данных тарифов:", err)
		return nil, err
	}

	return tariffs, nil
}
