package database

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddGameDB(r *http.Request) (string, error) {
	var (
		game                  Game
		name, genre, tariffId string
		objectTariffId        primitive.ObjectID
	)
	if r.Method != http.MethodPost {
		return "", fmt.Errorf("error method")
	}

	name = r.FormValue("name")
	genre = r.FormValue("genre")

	tariffId = r.URL.Query().Get("id")
	if tariffId == "" {
		return "", fmt.Errorf("error getting id tariff from URL")
	}

	game = Game{name, genre}

	db := MongoClient.Database("Vr")
	collection := db.Collection("Tariffs")

	// Конвертация string в ObjectID
	objectTariffId, err := primitive.ObjectIDFromHex(tariffId)
	if err != nil {
		return "", fmt.Errorf("error converting tariff ID to ObjectID: %v", err)
	}

	// Проверяем, существует ли игра с таким именем
	filter := bson.M{"_id": objectTariffId, "games": bson.M{"$elemMatch": bson.M{"name": name}}}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return "", fmt.Errorf("error checking existing game: %v", err)
	}
	if count > 0 {
		return "", fmt.Errorf("game with this name already exists")
	}

	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectTariffId},
		bson.M{"$push": bson.M{"games": game}},
	)
	if err != nil {
		return "", fmt.Errorf("error adding new game: %v", err)
	}
	return tariffId, nil
}

func AddDeviceDB(r *http.Request) (string, error) {
	var (
		device                   Device
		name, platform, tariffId string
		objectTariffId           primitive.ObjectID
	)
	if r.Method != http.MethodPost {
		return "", fmt.Errorf("error method")
	}
	name = r.FormValue("name")
	platform = r.FormValue("platform")

	tariffId = r.URL.Query().Get("id")
	if tariffId == "" {
		return "", fmt.Errorf("error getting id tariff from URL")
	}

	device = Device{name, platform}

	db := MongoClient.Database("Vr")
	collection := db.Collection("Tariffs")

	// Конвертация string в ObjectID
	objectTariffId, err := primitive.ObjectIDFromHex(tariffId)
	if err != nil {
		return "", fmt.Errorf("error converting tariff ID to ObjectID: %v", err)
	}

	// Проверяем, существует ли игра с таким именем
	filter := bson.M{"_id": objectTariffId, "devices": bson.M{"$elemMatch": bson.M{"name": name}}}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return "", fmt.Errorf("error checking existing device: %v", err)
	}
	if count > 0 {
		return "", fmt.Errorf("device with this name already exists")
	}

	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectTariffId},
		bson.M{"$push": bson.M{"devices": device}},
	)
	if err != nil {
		return "", fmt.Errorf("error adding new device: %v", err)
	}
	return tariffId, nil
}
