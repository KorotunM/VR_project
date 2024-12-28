package database

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteElementTariffDB(r *http.Request) error {
	var (
		tariffId       string
		requestData    AjaxDeleteElementTariff
		objectTariffId primitive.ObjectID
		filter         bson.M
		update         bson.M
		err            error
	)
	tariffId = r.URL.Query().Get("id")
	if tariffId == "" {
		return fmt.Errorf("error getting id tariff from URL")
	}
	// Конвертация string в ObjectID
	objectTariffId, err = primitive.ObjectIDFromHex(tariffId)
	if err != nil {
		return fmt.Errorf("error converting tariff ID to ObjectID: %v", err)
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		return fmt.Errorf("error getting data from template: %v", err)
	}

	db := MongoClient.Database("Vr")
	collection := db.Collection("Tariffs")

	// Формируем базовый фильтр с использованием _id
	filter = bson.M{"_id": objectTariffId}

	// Проверяем тип элемента (игра или устройство) и дополняем фильтр
	switch requestData.Type {
	case "game":
		filter["games.name"] = requestData.Name
		update = bson.M{"$pull": bson.M{"games": bson.M{"name": requestData.Name}}}
	case "device":
		filter["devices.name"] = requestData.Name
		update = bson.M{"$pull": bson.M{"devices": bson.M{"name": requestData.Name}}}
	default:
		return fmt.Errorf("invalid type in request data: %v", requestData.Type)
	}

	// Обновляем документ: удаляем указанный элемент
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error updating document in MongoDB: %v", err)
	}
	return nil
}

func DeleteTariffDB(tariffId string) error {
	var (
		objectTariffId primitive.ObjectID
		err            error
	)

	// Конвертируем ID тарифа в ObjectID
	objectTariffId, err = primitive.ObjectIDFromHex(tariffId)
	if err != nil {
		return fmt.Errorf("error converting tariff ID to ObjectID: %v", err)
	}

	// Подключаемся к базе данных
	db := MongoClient.Database("Vr")
	collectionTariffs := db.Collection("Tariffs")
	collectionBookings := db.Collection("Booking")

	// Найдём бронирования, связанные с тарифом
	filterBookings := bson.M{"tariff_id": objectTariffId}

	// Удаляем все бронирования, связанные с тарифом
	_, err = collectionBookings.DeleteMany(context.TODO(), filterBookings)
	if err != nil {
		return fmt.Errorf("error deleting bookings linked to the tariff: %v", err)
	}

	// Удаляем сам тариф
	_, err = collectionTariffs.DeleteOne(context.TODO(), bson.M{"_id": objectTariffId})
	if err != nil {
		return fmt.Errorf("error deleting tariff: %v", err)
	}

	return nil
}

func DeleteClientDB(r *http.Request) error {
	id := r.URL.Query().Get("id")
	if id == "" {
		return fmt.Errorf("missing client ID")
	}

	// Конвертируем ID клиента в ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid client ID: %v", err)
	}

	// Подключаемся к базе данных
	db := MongoClient.Database("Vr")
	clientsCollection := db.Collection("Clients")
	bookingsCollection := db.Collection("Booking")

	// Удаляем клиента
	clientFilter := bson.M{"_id": objectId}
	_, err = clientsCollection.DeleteOne(context.TODO(), clientFilter)
	if err != nil {
		return fmt.Errorf("error deleting client from database: %v", err)
	}

	// Удаляем связанные бронирования
	bookingsFilter := bson.M{"client_id": objectId}
	_, err = bookingsCollection.DeleteMany(context.TODO(), bookingsFilter)
	if err != nil {
		return fmt.Errorf("error deleting related bookings: %v", err)
	}

	return nil
}

func DeleteBookingDB(r *http.Request) error {
	id := r.URL.Query().Get("id")
	if id == "" {
		return fmt.Errorf("missing booking ID")
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid booking ID: %v", err)
	}

	db := MongoClient.Database("Vr")
	collection := db.Collection("Booking")

	filter := bson.M{"_id": objectId}
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting booking from database: %v", err)
	}
	return nil
}

func DeleteGeneralGameDB(r *http.Request) error {
	id := r.URL.Query().Get("id")
	if id == "" {
		return fmt.Errorf("missing general game ID")
	}

	// Конвертируем ID в ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid general game ID: %v", err)
	}

	// Подключаемся к базе данных
	db := MongoClient.Database("Vr")
	gamesCollection := db.Collection("Games")
	bookingsCollection := db.Collection("Booking")

	// Удаляем игру из коллекции Games
	filter := bson.M{"_id": objectId}
	_, err = gamesCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting general game from database: %v", err)
	}

	// Удаляем игру из всех бронирований
	updateFilter := bson.M{
		"general_games._id": objectId, // Фильтр ищет документы, где _id совпадает в массиве general_games
	}
	update := bson.M{
		"$pull": bson.M{
			"general_games": bson.M{"_id": objectId}, // Удаляем объект из массива general_games по _id
		},
	}

	_, err = bookingsCollection.UpdateMany(context.TODO(), updateFilter, update)
	if err != nil {
		return fmt.Errorf("error removing general game from bookings: %v", err)
	}

	return nil
}
