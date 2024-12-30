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

func AddGameDB(r *http.Request) (string, error) {
	var (
		game                  Game
		generalGame           GeneralGame
		name, genre, tariffId string
		objectTariffId        primitive.ObjectID
	)

	if r.Method != http.MethodPost {
		return "", fmt.Errorf("error: invalid method")
	}

	// Получаем данные из формы
	name = r.FormValue("name")
	genre = r.FormValue("genre")

	// Проверяем, заполнены ли обязательные поля
	if name == "" || genre == "" {
		return "", fmt.Errorf("error: all fields are required")
	}

	// Получаем ID тарифа из URL
	tariffId = r.URL.Query().Get("id")
	if tariffId == "" {
		return "", fmt.Errorf("error: missing tariff ID in URL")
	}

	// Конвертация string в ObjectID
	objectTariffId, err := primitive.ObjectIDFromHex(tariffId)
	if err != nil {
		return "", fmt.Errorf("error converting tariff ID to ObjectID: %v", err)
	}

	// Подключение к базе данных
	db := MongoClient.Database("Vr")
	tariffsCollection := db.Collection("Tariffs")
	gamesCollection := db.Collection("Games")

	// Проверяем, существует ли игра с таким именем в тарифе
	tariffFilter := bson.M{"_id": objectTariffId, "games": bson.M{"$elemMatch": bson.M{"name": name}}}
	tariffGameCount, err := tariffsCollection.CountDocuments(context.TODO(), tariffFilter)
	if err != nil {
		return "", fmt.Errorf("error checking existing game in tariff: %v", err)
	}
	if tariffGameCount > 0 {
		return "", fmt.Errorf("game with this name already exists")
	}

	// Извлекаем цену для игры из тарифа
	var tariff struct {
		PriceGame int `bson:"price_game"`
	}
	err = tariffsCollection.FindOne(context.TODO(), bson.M{"_id": objectTariffId}).Decode(&tariff)
	if err != nil {
		return "", fmt.Errorf("error retrieving tariff details: %v", err)
	}

	// Проверяем, существует ли игра с таким именем в общих играх
	generalGameFilter := bson.M{"name": name}
	update := bson.M{"$set": bson.M{"price": tariff.PriceGame, "genre": genre}}
	result, err := gamesCollection.UpdateOne(context.TODO(), generalGameFilter, update)
	if err != nil {
		return "", fmt.Errorf("error updating general game: %v", err)
	}

	// Если игра не была найдена, добавляем новую
	if result.MatchedCount == 0 {
		generalGame = GeneralGame{
			Name:  name,
			Genre: genre,
			Price: tariff.PriceGame,
		}
		_, err = gamesCollection.InsertOne(context.TODO(), generalGame)
		if err != nil {
			return "", fmt.Errorf("error adding game to general games: %v", err)
		}
	}

	// Добавляем игру в коллекцию игр тарифа
	game = Game{Name: name, Genre: genre}
	_, err = tariffsCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectTariffId},
		bson.M{"$push": bson.M{"games": game}},
	)
	if err != nil {
		return "", fmt.Errorf("error adding game to tariff: %v", err)
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

func AddTariffDB(r *http.Request) error {
	var (
		tariff    Tariff
		name      string
		price     int = -1
		priceGame int = -1
	)

	// Получаем параметры формы
	name = r.FormValue("name")
	price, err := strconv.Atoi(r.FormValue("price"))
	if err != nil {
		return fmt.Errorf("error converting string to int (price)")
	}
	priceGame, err = strconv.Atoi(r.FormValue("price_game"))
	if err != nil {
		return fmt.Errorf("error converting string to int (price game)")
	}

	// Проверяем, заполнены ли все поля
	if name == "" {
		return fmt.Errorf("error: name is required")
	}
	if price <= 0 || priceGame <= 0 {
		return fmt.Errorf("wrong price")
	}

	// Проверяем, существует ли тариф с таким названием
	db := MongoClient.Database("Vr")
	collection := db.Collection("Tariffs")

	// Проверка на существование тарифа с таким же названием
	filter := bson.M{"name": name}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error checking existing tariff: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("tariff with this name already exists")
	}

	// Создаем новый тариф
	tariff = Tariff{
		Name:      name,
		Price:     price,
		PriceGame: priceGame,
		Games:     []Game{},   // Пустой массив для игр
		Devices:   []Device{}, // Пустой массив для устройств
	}

	// Вставляем тариф в базу данных
	_, err = collection.InsertOne(context.TODO(), tariff)
	if err != nil {
		return fmt.Errorf("error adding new tariff: %v", err)
	}

	return nil
}

func AddClientDB(r *http.Request) error {
	// Получаем параметры формы
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	email := r.FormValue("email")

	// Проверяем, заполнены ли все поля
	if name == "" || phone == "" || email == "" {
		return fmt.Errorf("all fields are required")
	}

	// Подключение к базе данных MongoDB
	db := MongoClient.Database("Vr")
	collection := db.Collection("Clients")

	// Создаём структуру клиента
	client := Client{
		Name:  name,
		Phone: phone,
		Email: email,
	}

	// Вставляем клиента в базу данных
	_, err := collection.InsertOne(context.TODO(), client)
	if err != nil {
		return fmt.Errorf("error adding new client: %v", err)
	}

	return nil
}

func AddBookingDB(r *http.Request) error {
	// Получаем параметры формы
	clientID := r.FormValue("client")
	tariffID := r.FormValue("tariff")
	date := r.FormValue("date")
	timeSlot := r.FormValue("time")

	// Проверяем, заполнены ли все поля
	if clientID == "" || tariffID == "" || date == "" || timeSlot == "" {
		return fmt.Errorf("all fields are required")
	}

	// Конвертация clientID в ObjectID
	clientObjID, err := primitive.ObjectIDFromHex(clientID)
	if err != nil {
		return fmt.Errorf("error converting client ID to ObjectID: %v", err)
	}

	// Конвертация tariffID в ObjectID
	tariffObjID, err := primitive.ObjectIDFromHex(tariffID)
	if err != nil {
		return fmt.Errorf("error converting tariff ID to ObjectID: %v", err)
	}

	// Парсинг даты
	bookingDate, err := time.Parse("2006-01-02", date)
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

	// Подключение к базе данных MongoDB
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
	if count > 0 {
		return fmt.Errorf("time already exist")
	}

	// Структура бронирования
	booking := map[string]interface{}{
		"client_id":     clientObjID,
		"tariff_id":     tariffObjID,
		"general_games": generalGames,
		"booking_date":  bookingDate,
		"booking_time":  timeSlot,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, booking)
	return err
}

func AddGeneralGameDB(r *http.Request) error {
	// Получаем параметры формы
	name := r.FormValue("name")
	genre := r.FormValue("genre")
	price := r.FormValue("price")

	// Проверяем, заполнены ли все поля
	if name == "" || genre == "" || price == "" {
		return fmt.Errorf("all fields are required")
	}

	// Подключение к базе данных MongoDB
	db := MongoClient.Database("Vr")
	collection := db.Collection("Games")

	// Проверяем, существует ли игра с таким именем
	filter := bson.M{"name": name}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error checking existing game: %v", err)
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

	// Создаём структуру общей игры
	generalGame := GeneralGame{
		Name:  name,
		Genre: genre,
		Price: intPrice,
	}

	// Вставляем общую игру в базу данных
	_, err = collection.InsertOne(context.TODO(), generalGame)
	if err != nil {
		return fmt.Errorf("error adding new general game: %v", err)
	}

	return nil
}
