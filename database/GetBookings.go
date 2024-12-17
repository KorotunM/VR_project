package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// GetBookedTimes возвращает список забронированных времён для выбранной даты
func GetBookedTimes(date time.Time) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := MongoClient.Database("Vr").Collection("Booking")

	// Фильтр для поиска записей по дате
	filter := bson.M{
		"booking_date": date, // Передаём date как time.Time напрямую
	}
	log.Println("Фильтр для даты:", date)

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Сохраняем забронированные слоты
	var results []struct {
		BookingTime string `bson:"booking_time"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// Формируем список временных слотов
	bookedTimes := []string{}
	for _, result := range results {
		bookedTimes = append(bookedTimes, result.BookingTime)
	}

	return bookedTimes, nil
}
