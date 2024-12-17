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
	collection := db.Collection("Tariffs")

	// Удаляем тариф с указанным ID
	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": objectTariffId})
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
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid client ID: %v", err)
	}

	db := MongoClient.Database("Vr")
	collection := db.Collection("Clients")

	filter := bson.M{"_id": objectId}
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting client from database: %v", err)
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
