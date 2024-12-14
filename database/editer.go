package database

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func EditGameDB(r *http.Request) (string, error) {
	var (
		game                          Game
		gameName, gameGenre, tariffId string
		objectTariffId                primitive.ObjectID
	)

	// Получаем параметры из URL
	tariffId = r.URL.Query().Get("id")
	gameName = r.URL.Query().Get("name")
	gameGenre = r.URL.Query().Get("genre")

	if tariffId == "" || gameName == "" || gameGenre == "" {
		return "", fmt.Errorf("missing required parameters in URL")
	}

	// Конвертируем tariffId в ObjectID
	objectTariffId, err := primitive.ObjectIDFromHex(tariffId)
	if err != nil {
		return "", fmt.Errorf("error converting tariff ID to ObjectID: %v", err)
	}

	// Создаем объект с новыми данными
	game = Game{
		Name:  r.FormValue("name"),  // Новое название игры
		Genre: r.FormValue("genre"), // Новый жанр игры
	}

	// Подключаемся к базе данных
	db := MongoClient.Database("Vr")
	collection := db.Collection("Tariffs")

	// Проверяем, существует ли игра с таким же названием в тарифе
	filter := bson.M{
		"_id": objectTariffId,
		"games": bson.M{
			"$elemMatch": bson.M{
				"name":  gameName,  // Ищем игру по старому имени
				"genre": gameGenre, // и старому жанру
			},
		},
	}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return "", fmt.Errorf("error checking existing game: %v", err)
	}
	if count == 0 {
		return "", fmt.Errorf("game not found")
	}

	// Обновляем запись игры
	update := bson.M{
		"$set": bson.M{
			"games.$.name":  game.Name,
			"games.$.genre": game.Genre,
		},
	}

	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectTariffId, "games.name": gameName, "games.genre": gameGenre},
		update,
	)
	if err != nil {
		return "", fmt.Errorf("error updating game: %v", err)
	}

	// Возвращаем tariffId для редиректа
	return tariffId, nil
}

func EditDeviceDB(r *http.Request) (string, error) {
	var (
		device                               Device
		deviceName, devicePlatform, tariffId string
		objectTariffId                       primitive.ObjectID
	)

	// Получаем параметры из URL
	tariffId = r.URL.Query().Get("id")
	deviceName = r.URL.Query().Get("name")
	devicePlatform = r.URL.Query().Get("platform")

	if tariffId == "" || deviceName == "" || devicePlatform == "" {
		return "", fmt.Errorf("missing required parameters in URL")
	}

	// Конвертируем tariffId в ObjectID
	objectTariffId, err := primitive.ObjectIDFromHex(tariffId)
	if err != nil {
		return "", fmt.Errorf("error converting tariff ID to ObjectID: %v", err)
	}

	// Создаем объект с новыми данными
	device = Device{
		Name:     r.FormValue("name"),     // Новое название устройства
		Platform: r.FormValue("platform"), // Новая платформа устройства
	}

	// Подключаемся к базе данных
	db := MongoClient.Database("Vr")
	collection := db.Collection("Tariffs")

	// Проверяем, существует ли устройство с таким же названием в тарифе
	filter := bson.M{
		"_id": objectTariffId,
		"devices": bson.M{
			"$elemMatch": bson.M{
				"name":     deviceName,     // Ищем устройство по старому имени
				"platform": devicePlatform, // и старой платформе
			},
		},
	}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return "", fmt.Errorf("error checking existing device: %v", err)
	}
	if count == 0 {
		return "", fmt.Errorf("device not found")
	}

	// Обновляем запись устройства
	update := bson.M{
		"$set": bson.M{
			"devices.$.name":     device.Name,
			"devices.$.platform": device.Platform,
		},
	}

	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectTariffId, "devices.name": deviceName, "devices.platform": devicePlatform},
		update,
	)
	if err != nil {
		return "", fmt.Errorf("error updating device: %v", err)
	}

	// Возвращаем tariffId для редиректа
	return tariffId, nil
}
