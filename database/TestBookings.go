package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// PrintAllBookings выводит все записи из коллекции Bookings
func PrintAllBookings() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Получаем коллекцию Bookings
	collection := MongoClient.Database("Vr").Collection("Booking")

	// Пустой фильтр, чтобы получить все записи
	filter := bson.M{}

	// Выполняем запрос
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatalf("Ошибка при получении данных из MongoDB: %v", err)
	}
	defer cursor.Close(ctx)

	// Итерируем по результатам
	fmt.Println("Все записи в коллекции 'Bookings':")
	for cursor.Next(ctx) {
		var booking BookingDocument
		if err := cursor.Decode(&booking); err != nil {
			log.Printf("Ошибка декодирования записи: %v", err)
			continue
		}

		// Выводим запись в консоль
		fmt.Printf("ID: %v\n", booking.ID)
		fmt.Printf("Client ID: %s\n", booking.ClientID)
		fmt.Printf("Tariff ID: %s\n", booking.TariffID)
		fmt.Printf("Booking Date: %s\n", booking.BookingDate)
		fmt.Printf("Booking Time: %s\n", booking.BookingTime)
		fmt.Println("--------------------------")
	}

	// Проверка на ошибки после итерации
	if err := cursor.Err(); err != nil {
		log.Fatalf("Ошибка при обработке данных: %v", err)
	}
	log.Println("Вывод всех записей завершен")
}
