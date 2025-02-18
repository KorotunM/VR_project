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
		oldGameName, oldGameGenre string
		newGameName, newGameGenre string
		tariffId                  string
		objectTariffId            primitive.ObjectID
	)

	// Получаем параметры из URL
	tariffId = r.URL.Query().Get("id")
	oldGameName = r.URL.Query().Get("name")   // Старое название игры
	oldGameGenre = r.URL.Query().Get("genre") // Старый жанр игры

	// Получаем новые значения из формы
	newGameName = r.FormValue("name")   // Новое название игры
	newGameGenre = r.FormValue("genre") // Новый жанр игры

	if tariffId == "" || oldGameName == "" || oldGameGenre == "" || newGameName == "" || newGameGenre == "" {
		return "", fmt.Errorf("missing required parameters")
	}

	// Конвертируем tariffId в ObjectID
	objectTariffId, err := primitive.ObjectIDFromHex(tariffId)
	if err != nil {
		return "", fmt.Errorf("error converting tariff ID to ObjectID: %v", err)
	}

	// Подключаемся к базе данных
	db := MongoClient.Database("Vr")
	tariffsCollection := db.Collection("Tariffs")
	gamesCollection := db.Collection("Games")

	// Проверяем, не существует ли уже игра с таким же новым названием в тарифе
	newGameFilter := bson.M{
		"_id": objectTariffId,
		"games": bson.M{
			"$elemMatch": bson.M{
				"name": newGameName,
			},
		},
	}
	newGameCount, err := tariffsCollection.CountDocuments(context.TODO(), newGameFilter)
	if err != nil {
		return "", fmt.Errorf("error checking if new game name exists: %v", err)
	}
	if newGameName == oldGameName {
		if newGameCount != 1 {
			return "", fmt.Errorf("game with this name already exists or another error")
		}
	} else {
		if newGameCount > 0 {
			return "", fmt.Errorf("game with this name already exists")
		}
	}

	// Обновляем игру в тарифе
	update := bson.M{
		"$set": bson.M{
			"games.$.name":  newGameName,
			"games.$.genre": newGameGenre,
		},
	}
	_, err = tariffsCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectTariffId, "games.name": oldGameName},
		update,
	)
	if err != nil {
		return "", fmt.Errorf("error updating game in tariff: %v", err)
	}

	// Извлекаем цену для игры из тарифа
	var tariff struct {
		PriceGame int `bson:"price_game"`
	}
	err = tariffsCollection.FindOne(context.TODO(), bson.M{"_id": objectTariffId}).Decode(&tariff)
	if err != nil {
		return "", fmt.Errorf("error retrieving tariff details: %v", err)
	}

	// Проверяем, существует ли игра с прежним названием в общих играх
	generalGameFilter := bson.M{"name": oldGameName}
	updateGeneralGame := bson.M{
		"$set": bson.M{
			"name":  newGameName,
			"genre": newGameGenre,
			"price": tariff.PriceGame, // Обновляем цену игры
		},
	}
	_, err = gamesCollection.UpdateOne(context.TODO(), generalGameFilter, updateGeneralGame)
	if err != nil {
		return "", fmt.Errorf("error updating general game: %v", err)
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
		tariffId, name   string
		objectTariffId   primitive.ObjectID
		err              error
		price, priceGame int
		oldTariffData    struct {
			PriceGame int `bson:"price_game"`
			Games     []struct {
				Name string `bson:"name"`
			} `bson:"games"`
		}
	)

	tariffId = r.URL.Query().Get("id")
	name = r.FormValue("name")
	price, err = strconv.Atoi(r.FormValue("price"))
	if err != nil {
		return "", fmt.Errorf("missing price parameters from URL: %v", err)
	}
	priceGame, err = strconv.Atoi(r.FormValue("price_game"))
	if err != nil {
		return "", fmt.Errorf("missing price game parameters from URL: %v", err)
	}

	if tariffId == "" || name == "" {
		return "", fmt.Errorf("missing required string parameters from URL")
	}

	if price <= 0 || priceGame <= 0 {
		return "", fmt.Errorf("wrong price")
	}

	// Конвертация tariffId в ObjectID
	objectTariffId, err = primitive.ObjectIDFromHex(tariffId)
	if err != nil {
		return "", fmt.Errorf("error converting tariff ID to ObjectID: %v", err)
	}

	// Подключение к базе данных
	db := MongoClient.Database("Vr")
	collectionTariffs := db.Collection("Tariffs")
	collectionGames := db.Collection("Games")

	// Проверка существования тарифа с таким именем
	filter := bson.M{
		"_id":  bson.M{"$ne": objectTariffId}, // Исключить текущий тариф
		"name": name,                          // Проверить по имени
	}
	count, err := collectionTariffs.CountDocuments(context.TODO(), filter)
	if err != nil {
		return "", fmt.Errorf("error checking if tariff name exists: %v", err)
	}
	if count > 0 {
		return "", fmt.Errorf("tariff with this name already exists")
	}

	// Получение старых данных тарифа
	err = collectionTariffs.FindOne(context.TODO(), bson.M{"_id": objectTariffId}).Decode(&oldTariffData)
	if err != nil {
		return "", fmt.Errorf("error retrieving old tariff data: %v", err)
	}

	// Обновление данных тарифа
	update := bson.M{
		"$set": bson.M{
			"name":       name,
			"price":      price,
			"price_game": priceGame,
		},
	}
	_, err = collectionTariffs.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectTariffId},
		update,
	)
	if err != nil {
		return "", fmt.Errorf("error updating tariff: %v", err)
	}

	// Если цена за игру изменилась, обновляем её в общих играх
	if oldTariffData.PriceGame != priceGame {
		for _, game := range oldTariffData.Games {
			_, err = collectionGames.UpdateOne(
				context.TODO(),
				bson.M{"name": game.Name}, // Обновляем по названию игры
				bson.M{"$set": bson.M{"price": priceGame}},
			)
			if err != nil {
				return "", fmt.Errorf("error updating price for game '%s': %v", game.Name, err)
			}
		}
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
		generalGameId, name, genre, price string
		objectGeneralGameId               primitive.ObjectID
		err                               error
	)

	generalGameId = r.URL.Query().Get("id")
	name = r.FormValue("name")
	genre = r.FormValue("genre")
	price = r.FormValue("price")

	if generalGameId == "" || name == "" || genre == "" || price == "" {
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

	intPrice, err := strconv.Atoi(price)
	if err != nil {
		return fmt.Errorf("error converting game's price: %v", err)
	}

	if intPrice <= 0 {
		return fmt.Errorf("wrong price")
	}

	update := bson.M{
		"$set": bson.M{
			"name":  name,
			"genre": genre,
			"price": intPrice,
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
