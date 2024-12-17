package database

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Подключение к MongoDB и получения данных клиентов
func GetClients() ([]Client, error) {

	// для автоматического выключения, если запрос завис
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Работа с коллекцией
	db := MongoClient.Database("Vr")
	collection := db.Collection("Clients")

	// Поиск всех документов
	filter := bson.M{} // Пустой фильтр для получения всех документов
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error search clients: %v", err)
	}
	defer cursor.Close(ctx)

	// Считывание результатов
	var clients []Client
	if err := cursor.All(ctx, &clients); err != nil {
		return nil, fmt.Errorf("error clients reading: %v", err)
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

func GetAllBookings() ([]BookingDocument, error) {
	var bookings []BookingDocument

	// для автоматического выключения, если запрос завис
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Работа с коллекцией
	db := MongoClient.Database("Vr")
	collection := db.Collection("Booking")

	filter := bson.M{} // Получаем все бронирования
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error search bookings: %v", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, fmt.Errorf("error bookings reading: %v", err)
	}

	// Преобразуем дату в нужный формат и сохраняем в поле Date
	for index := range bookings {
		// Заполняем Date как строку в нужном формате "гггг-мм-дд"
		bookings[index].Date = bookings[index].BookingDate.Format("2006-01-02")
	}

	return bookings, nil
}
