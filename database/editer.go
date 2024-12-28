package database

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

	// Проверяем, не существует ли уже новая игра с таким же названием
	newGameFilter := bson.M{
		"_id": objectTariffId,
		"games": bson.M{
			"$elemMatch": bson.M{
				"name": game.Name, // Проверка нового названия
			},
		},
	}
	newGameCount, err := collection.CountDocuments(context.TODO(), newGameFilter)
	if err != nil {
		return "", fmt.Errorf("error checking if new game name exists: %v", err)
	}
	if game.Name == gameName {
		if newGameCount != 1 {
			return "", fmt.Errorf("game with this name already exists or another error")
		}
	} else {
		if newGameCount > 0 {
			return "", fmt.Errorf("game with this name already exists")
		}
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

	// Проверяем, не существует ли уже новое устройство с таким же названием
	newDeviceFilter := bson.M{
		"_id": objectTariffId,
		"devices": bson.M{
			"$elemMatch": bson.M{
				"name": device.Name, // Проверка нового названия
			},
		},
	}
	newDeviceCount, err := collection.CountDocuments(context.TODO(), newDeviceFilter)
	if err != nil {
		return "", fmt.Errorf("error checking if new device name exists: %v", err)
	}
	if device.Name == deviceName {
		if newDeviceCount != 1 {
			return "", fmt.Errorf("device with this name already exists or another error")
		}
	} else {
		if newDeviceCount > 0 {
			return "", fmt.Errorf("device with this name already exists")
		}
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

func EditTariffDB(r *http.Request) (string, error) {
	var (
		tariffId, name string
		objectTariffId primitive.ObjectID
		err            error
		price          int
	)

	tariffId = r.URL.Query().Get("id")
	name = r.FormValue("name")
	price, err = strconv.Atoi(r.FormValue("price"))

	if tariffId == "" || name == "" {
		return "", fmt.Errorf("missing required string parameters from URL")
	}
	if err != nil {
		return "", fmt.Errorf("missing price parameters from URL: %v", err)
	}

	// Конвертация tariffId в ObjectID
	objectTariffId, err = primitive.ObjectIDFromHex(tariffId)
	if err != nil {
		return "", fmt.Errorf("error converting tariff ID to ObjectID: %v", err)
	}

	// Проверка существования тарифа с таким именем
	db := MongoClient.Database("Vr")
	collection := db.Collection("Tariffs")

	filter := bson.M{
		"_id":  bson.M{"$ne": objectTariffId}, // Исключить текущий тариф
		"name": name,                          // Проверить по имени
	}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return "", fmt.Errorf("error checking if tariff name exists: %v", err)
	}
	if count > 0 {
		return "", fmt.Errorf("tariff with this name already exists")
	}

	// Обновление данных тарифа
	update := bson.M{
		"$set": bson.M{
			"name":  name,
			"price": price,
		},
	}

	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectTariffId},
		update,
	)
	if err != nil {
		return "", fmt.Errorf("error updating tariff: %v", err)
	}

	return tariffId, nil
}

func EditClientDB(r *http.Request) error {
	var (
		clientId, name, email, phone string
		objectClientId               primitive.ObjectID
		err                          error
	)

	clientId = r.URL.Query().Get("id")
	name = r.FormValue("name")
	email = r.FormValue("email")
	phone = r.FormValue("phone")

	if clientId == "" || name == "" || email == "" || phone == "" {
		return fmt.Errorf("missing required parameters")
	}

	// Конвертация clientId в ObjectID
	objectClientId, err = primitive.ObjectIDFromHex(clientId)
	if err != nil {
		return fmt.Errorf("error converting client ID to ObjectID: %v", err)
	}

	// Обновление данных клиента
	db := MongoClient.Database("Vr")
	collection := db.Collection("Clients")

	update := bson.M{
		"$set": bson.M{
			"name":         name,
			"email":        email,
			"phone number": phone,
		},
	}

	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectClientId},
		update,
	)
	if err != nil {
		return fmt.Errorf("error updating client: %v", err)
	}

	return nil
}

func EditBookingDB(r *http.Request) error {
	var (
		bookingID, clientID, tariffID, date, timeSlot string
		objectBookingID, clientObjID, tariffObjID     primitive.ObjectID
		bookingDate                                   time.Time
		err                                           error
	)

	// Получаем параметры из URL и формы
	bookingID = r.URL.Query().Get("id")
	oldData := r.URL.Query().Get("date")
	oldTime := r.URL.Query().Get("time")
	clientID = r.FormValue("client")
	tariffID = r.FormValue("tariff")
	date = r.FormValue("date")
	timeSlot = r.FormValue("time")

	// Проверка обязательных параметров
	if bookingID == "" || clientID == "" || tariffID == "" || date == "" || timeSlot == "" {
		return fmt.Errorf("missing required parameters")
	}

	// Конвертация bookingID в ObjectID
	objectBookingID, err = primitive.ObjectIDFromHex(bookingID)
	if err != nil {
		return fmt.Errorf("error converting booking ID to ObjectID: %v", err)
	}

	// Конвертация clientID в ObjectID
	clientObjID, err = primitive.ObjectIDFromHex(clientID)
	if err != nil {
		return fmt.Errorf("error converting client ID to ObjectID: %v", err)
	}

	// Конвертация tariffID в ObjectID
	tariffObjID, err = primitive.ObjectIDFromHex(tariffID)
	if err != nil {
		return fmt.Errorf("error converting tariff ID to ObjectID: %v", err)
	}

	// Парсинг даты
	bookingDate, err = time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("error parsing booking date: %v", err)
	}

	// Обработка поля general_games
	generalGameIDs := r.Form["general-games"] // Получаем массив значений из формы
	var generalGames []bson.M                 // Массив для объектов с _id
	for _, gameID := range generalGameIDs {
		gameObjID, err := primitive.ObjectIDFromHex(gameID)
		if err != nil {
			return fmt.Errorf("error converting general game ID to ObjectID: %v", err)
		}
		generalGames = append(generalGames, bson.M{"_id": gameObjID})
	}

	// Подключаемся к базе данных
	db := MongoClient.Database("Vr")
	collection := db.Collection("Booking")

	// Формируем фильтр для поиска бронирований с такой же датой и временем
	filter := bson.M{
		"booking_date": bookingDate,
		"booking_time": timeSlot,
	}

	// Проверяем, есть ли уже запись с таким же днем и временем
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error checking if time is already booked: %v", err)
	}
	if oldData == date && oldTime == timeSlot {
		if count != 1 {
			return fmt.Errorf("time already exist")
		}
	} else {
		if count > 0 {
			return fmt.Errorf("time already exist")
		}
	}

	// Структура обновления
	update := bson.M{
		"$set": bson.M{
			"client_id":     clientObjID,
			"tariff_id":     tariffObjID,
			"general_games": generalGames,
			"booking_date":  bookingDate,
			"booking_time":  timeSlot,
		},
	}

	// Выполняем обновление
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objectBookingID},
		update,
	)
	if err != nil {
		return fmt.Errorf("error updating booking: %v", err)
	}

	return nil
}

func EditGeneralGameDB(r *http.Request) error {
	var (
		generalGameId, name, genre string
		objectGeneralGameId        primitive.ObjectID
		err                        error
	)

	generalGameId = r.URL.Query().Get("id")
	name = r.FormValue("name")
	genre = r.FormValue("genre")

	if generalGameId == "" || name == "" || genre == "" {
		return fmt.Errorf("missing required parameters")
	}

	// Конвертация generalGameId в ObjectID
	objectGeneralGameId, err = primitive.ObjectIDFromHex(generalGameId)
	if err != nil {
		return fmt.Errorf("error converting general game ID to ObjectID: %v", err)
	}

	// Обновление данных клиента
	db := MongoClient.Database("Vr")
	collection := db.Collection("Games")

	// Проверка существования другой игры с таким же именем
	filter := bson.M{
		"_id":  bson.M{"$ne": objectGeneralGameId}, // Исключаем текущую игру
		"name": name,                               // Проверяем по имени
	}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error checking if game name exists: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("game with this name already exists")
	}

	update := bson.M{
		"$set": bson.M{
			"name":  name,
			"genre": genre,
		},
	}

	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectGeneralGameId},
		update,
	)
	if err != nil {
		return fmt.Errorf("error updating general game: %v", err)
	}

	return nil
}
