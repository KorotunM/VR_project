package database

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		log.Fatalf("Error connection to MongoDB: %v", err)
	}

	// Проверяем подключение
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB is not responding: %v", err)
	}

	MongoClient = client
}

// Закрытие подключения
func CloseMongoDB() {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := MongoClient.Disconnect(ctx); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
		log.Println("Connection to MongoDB is closed")
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

func GetAllTariffs() ([]TariffTitle, error) {
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

func GetTariff(r *http.Request) (Tariff, error) {
	var (
		tariffId       string
		objectTariffId primitive.ObjectID
		tariff         Tariff
	)
	tariffId = r.URL.Query().Get("id")
	if tariffId == "" {
		return tariff, fmt.Errorf("error taking id")
	}

	// Конвертируем строковый ID в ObjectID
	objectTariffId, err := primitive.ObjectIDFromHex(tariffId)
	if err != nil {
		return tariff, fmt.Errorf("error convertation: %v", err)
	}

	db := MongoClient.Database("Vr")
	collection := db.Collection("Tariffs")

	err = collection.FindOne(context.TODO(), bson.M{"_id": objectTariffId}).Decode(&tariff)
	if err != nil {
		return tariff, fmt.Errorf("error taking a tariff from MongoDB: %v", err)
	}
	return tariff, nil
}
