package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertBooking(clientID, tariffID, date, timeSlot string) error {
	collection := MongoClient.Database("Vr").Collection("Booking")

	// Преобразуем строки в ObjectID
	clientObjID, err := primitive.ObjectIDFromHex(clientID)
	if err != nil {
		return err
	}

	tariffObjID, err := primitive.ObjectIDFromHex(tariffID)
	if err != nil {
		return err
	}

	// Парсим дату бронирования
	bookingDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return err
	}

	// Структура бронирования
	booking := map[string]interface{}{
		"client_id":    clientObjID,
		"tariff_id":    tariffObjID,
		"booking_date": bookingDate,
		"booking_time": timeSlot,
	}
	log.Println("Бронирование: ", booking)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, booking)
	return err
}
