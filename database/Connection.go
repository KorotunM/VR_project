package database

import (
	"context"
	"encoding/json"
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

func DeleteElementTariffDB(r *http.Request) error {
	var (
		requestData AjaxDeleteElementTariff
		filter      bson.M
		update      bson.M
		err         error
	)
	if r.Method != http.MethodPost {
		return fmt.Errorf("error method: %v", err)
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		return fmt.Errorf("error getting data from template: %v", err)
	}

	db := MongoClient.Database("Vr")
	collection := db.Collection("Tariffs")

	// Проверяем тип элемента (игра или устройство)
	switch requestData.Type {
	case "game":
		filter = bson.M{"games.name": requestData.Name}
		update = bson.M{"$pull": bson.M{"games": bson.M{"name": requestData.Name}}}
	case "device":
		filter = bson.M{"devices.name": requestData.Name}
		update = bson.M{"$pull": bson.M{"devices": bson.M{"name": requestData.Name}}}
	default:
		return fmt.Errorf("error data: %v", err)
	}

	// Обновляем документ: удаляем указанный элемент
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error deleting data: %v", err)
	}
	return nil
}
