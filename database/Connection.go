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

var MongoClient *mongo.Client

// Инициализация клиента MongoDB
func InitMongoDB(uri string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Ошибка подключения к MongoDB: %v", err)
	}

	// Проверяем подключение
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB не отвечает: %v", err)
	}

	MongoClient = client
	log.Println("Подключение к MongoDB успешно инициализировано")
}

// Закрытие подключения
func CloseMongoDB() {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := MongoClient.Disconnect(ctx); err != nil {
			log.Fatalf("Ошибка отключения от MongoDB: %v", err)
		}
		log.Println("Подключение к MongoDB закрыто")
	}
}

// Подключение к MongoDB и получения данных клиентов
func GetClients() ([]Client, error) {

	// для автоматического выключения, если запрос завис
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Работа с коллекцией
	db := MongoClient.Database("Vr")
	collection := db.Collection("Users")

	// Поиск всех документов
	filter := bson.M{} // Пустой фильтр для получения всех документов
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error search users: %v", err)
	}
	defer cursor.Close(ctx)

	// Считывание результатов
	var clients []Client
	if err := cursor.All(ctx, &clients); err != nil {
		return nil, fmt.Errorf("error users reading: %v", err)
	}

	return clients, nil
}

func GetTariffs() ([]TariffTitle, error) {
	var tariffs []TariffTitle

	// для автоматического выключения, если запрос завис
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Работа с коллекцией
	db := MongoClient.Database("Vr")
	collection := db.Collection("Tariffs")

	filter := bson.M{} // Получаем все тарифы
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error search tariffs: %v", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &tariffs); err != nil {
		return nil, fmt.Errorf("error tariffs reading: %v", err)
	}

	return tariffs, nil

	// // Обрабатываем курсор вручную, чтобы преобразовать _id в строку
	// for cursor.Next(ctx) {
	//     var tariff TariffTitle
	//     var raw bson.M

	//     // Декодируем документ в bson.M
	//     if err := cursor.Decode(&raw); err != nil {
	//         return nil, fmt.Errorf("error decoding tariff: %v", err)
	//     }

	//     // Преобразуем _id в строку
	//     if id, ok := raw["_id"].(primitive.ObjectID); ok {
	//         tariff.ID = id.Hex() // Сохраняем _id как строку
	//     } else {
	//         return nil, fmt.Errorf("invalid _id type")
	//     }

	//     // Декодируем остальные поля
	//     if name, ok := raw["name"].(string); ok {
	//         tariff.Name = name
	//     } else {
	//         return nil, fmt.Errorf("invalid name type")
	//     }

	//     tariffs = append(tariffs, tariff)
	// }

	// if err := cursor.Err(); err != nil {
	//     return nil, fmt.Errorf("error iterating tariffs: %v", err)
	// }

	// return tariffs, nil
}
