package database

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Получение данных клиентов
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

func GetBookingStatistics() ([]TariffStats, error) {
	// Подключаемся к базе данных
	db := MongoClient.Database("Vr")
	bookingCollection := db.Collection("Booking")
	tariffCollection := db.Collection("Tariffs")

	// Получаем все бронирования
	cursor, err := bookingCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error fetching bookings: %v", err)
	}
	defer cursor.Close(context.TODO())

	// Создаём структуру для хранения статистики по тарифам
	statsMap := make(map[string]TariffStats)

	// Обрабатываем все бронирования
	for cursor.Next(context.TODO()) {
		var booking BookingDocument
		err := cursor.Decode(&booking)
		if err != nil {
			return nil, fmt.Errorf("error decoding booking document: %v", err)
		}

		// Конвертация bookingID в ObjectID
		objectBookingID, err := primitive.ObjectIDFromHex(booking.TariffID)
		if err != nil {
			return nil, fmt.Errorf("error converting booking ID to ObjectID: %v", err)
		}

		// Получаем тариф по тарифу ID
		var tariff Tariff
		err = tariffCollection.FindOne(context.TODO(), bson.M{"_id": objectBookingID}).Decode(&tariff)
		if err != nil {
			return nil, fmt.Errorf("error fetching tariff: %v", err)
		}

		// Обновляем статистику по тарифу
		if stats, exists := statsMap[tariff.Name]; exists {
			stats.BookingsCount++
			stats.CurrentProfit += tariff.Price
			statsMap[tariff.Name] = stats
		} else {
			statsMap[tariff.Name] = TariffStats{
				TariffName:    tariff.Name,
				CurrentProfit: tariff.Price,
				BookingsCount: 1,
			}
		}
	}

	// Преобразуем карту в срез
	var statsSlice []TariffStats
	for _, stats := range statsMap {
		statsSlice = append(statsSlice, stats)
	}

	return statsSlice, nil
}

func GetDailyBookingStatistics() ([]DailyStats, error) {
	// Подключаемся к базе данных
	db := MongoClient.Database("Vr")
	bookingCollection := db.Collection("Booking")

	// Формируем запрос для выборки по дням
	cursor, err := bookingCollection.Aggregate(context.TODO(), []bson.M{
		{
			"$project": bson.M{
				"date": bson.M{
					"$dateToString": bson.M{
						"format": "%Y-%m-%d", // Форматируем только дату, без времени
						"date":   "$booking_date",
					},
				},
			},
		},
		{
			"$group": bson.M{
				"_id":   "$date",           // Группируем по дате
				"count": bson.M{"$sum": 1}, // Считаем количество бронирований
			},
		},
		{
			"$sort": bson.M{
				"_id": 1, // Сортируем по дате
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("error fetching daily booking statistics: %v", err)
	}
	defer cursor.Close(context.TODO())

	// Создаем срез для статистики по дням
	var dailyStats []DailyStats
	for cursor.Next(context.TODO()) {
		var result struct {
			Date  string `bson:"_id"`
			Count int    `bson:"count"`
		}
		err := cursor.Decode(&result)
		if err != nil {
			return nil, fmt.Errorf("error decoding result: %v", err)
		}

		dailyStats = append(dailyStats, DailyStats{
			Date:          result.Date,
			BookingsCount: result.Count,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over cursor: %v", err)
	}

	return dailyStats, nil
}

func GetAllGeneralGames() ([]GeneralGame, error) {
	var generalGames []GeneralGame

	// для автоматического выключения, если запрос завис
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Работа с коллекцией
	db := MongoClient.Database("Vr")
	collection := db.Collection("Games")

	filter := bson.M{} // Получаем все тарифы
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error search tariffs: %v", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &generalGames); err != nil {
		return nil, fmt.Errorf("error tariffs reading: %v", err)
	}

	return generalGames, nil
}
