package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ExampleDocument struct {
	Name  string `bson:"name"`
	Value int    `bson:"value"`
}

func Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Подключение к MongoDB (без аутентификации)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Ошибка подключения к MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Ошибка отключения: %v", err)
		}
	}()

	// Проверка подключения
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB не отвечает: %v", err)
	}

	fmt.Println("Успешное подключение к MongoDB")

	// Работа с базой данных и коллекцией
	db := client.Database("Vr")          // Замените "exampleDB" на имя вашей базы
	collection := db.Collection("Users") // Замените "example" на имя коллекции

	fmt.Printf("Работаем с коллекцией: %s\n", collection.Name())

	filter := bson.M{}                        // Пустой фильтр для получения всех документов
	projection := bson.M{"_id": 1, "name": 1} // Включить только _id и name

	options := options.Find().SetProjection(projection)

	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		log.Fatalf("Ошибка поиска: %v", err)
	}
	defer cursor.Close(ctx)

	// Чтение документов из курсора
	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		log.Fatalf("Ошибка чтения результатов: %v", err)
	}

	fmt.Println("Результаты в формате JSON:")
	fmt.Println(results)
}

//mongodb://localhost:27017/
