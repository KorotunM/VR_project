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
	if r.Method != http.MethodPost {
		return fmt.Errorf("error method")
	}
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
