package database

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	// Получаем тариф ID из URL
	tariffId = r.URL.Query().Get("id")
	if tariffId == "" {
		return fmt.Errorf("error getting id tariff from URL")
	}

	// Конвертация string в ObjectID
	objectTariffId, err = primitive.ObjectIDFromHex(tariffId)
	if err != nil {
		return fmt.Errorf("error converting tariff ID to ObjectID: %v", err)
	}

	// Декодируем данные из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		return fmt.Errorf("error decoding request body: %v", err)
	}

	// Подключаемся к базе данных
	db := MongoClient.Database("Vr")
	tariffsCollection := db.Collection("Tariffs")
	gamesCollection := db.Collection("Games")
	bookingsCollection := db.Collection("Booking")

	// Формируем базовый фильтр с использованием _id
	filter = bson.M{"_id": objectTariffId}

	// Удаляем элемент из тарифа в зависимости от типа (игра или устройство)
	switch requestData.Type {
	case "game":
		filter["games.name"] = requestData.Name
		update = bson.M{"$pull": bson.M{"games": bson.M{"name": requestData.Name}}}

		// Удаляем игру из тарифа
		_, err = tariffsCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return fmt.Errorf("error updating document in MongoDB: %v", err)
		}

		// Пытаемся удалить игру из таблицы общих игр, если она там есть
		gameFilter := bson.M{"name": requestData.Name}
		gameToDelete := gamesCollection.FindOne(context.TODO(), gameFilter)

		// Проверяем, существует ли игра в коллекции
		var gameData struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		err = gameToDelete.Decode(&gameData)
		if err == nil {
			// Если игра существует, удаляем её из коллекции общих игр
			_, err = gamesCollection.DeleteOne(context.TODO(), gameFilter)
			if err != nil {
				return fmt.Errorf("error deleting game from general games collection: %v", err)
			}

			// Удаляем игру из всех бронирований по `_id`
			updateFilter := bson.M{
				"general_games._id": gameData.ID,
			}
			updateBookings := bson.M{
				"$pull": bson.M{
					"general_games": bson.M{"_id": gameData.ID},
				},
			}

			_, err = bookingsCollection.UpdateMany(context.TODO(), updateFilter, updateBookings)
			if err != nil {
				return fmt.Errorf("error removing general game from bookings: %v", err)
			}
		}

	case "device":
		filter["devices.name"] = requestData.Name
		update = bson.M{"$pull": bson.M{"devices": bson.M{"name": requestData.Name}}}

		// Обновляем документ: удаляем указанный элемент
		_, err = tariffsCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return fmt.Errorf("error updating document in MongoDB: %v", err)
		}

	default:
		return fmt.Errorf("invalid type in request data: %v", requestData.Type)
	}

	return nil
}

func DeleteTariffDB(tariffId string) error {
	var (
		objectTariffId primitive.ObjectID
		err            error
		tariffData     struct {
			Games []struct {
				Name string `bson:"name"`
			} `bson:"games"`
		}
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
	collectionGames := db.Collection("Games")

	// Получаем данные тарифа перед удалением
	err = collectionTariffs.FindOne(context.TODO(), bson.M{"_id": objectTariffId}).Decode(&tariffData)
	if err != nil {
		return fmt.Errorf("error finding tariff data: %v", err)
	}

	// Удаляем игры, связанные с тарифом
	for _, game := range tariffData.Games {
		// Находим `_id` игры в таблице общих игр по имени
		var gameData struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		err = collectionGames.FindOne(context.TODO(), bson.M{"name": game.Name}).Decode(&gameData)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				continue // Если игра уже удалена, пропускаем
			}
			return fmt.Errorf("error finding game '%s': %v", game.Name, err)
		}

		// Удаляем игру из таблицы общих игр
		_, err = collectionGames.DeleteOne(context.TODO(), bson.M{"_id": gameData.ID})
		if err != nil {
			return fmt.Errorf("error deleting game '%s' from general games collection: %v", game.Name, err)
		}

		// Удаляем игру из всех бронирований
		updateFilter := bson.M{
			"general_games._id": gameData.ID,
		}
		update := bson.M{
			"$pull": bson.M{
				"general_games": bson.M{"_id": gameData.ID},
			},
		}
		_, err = collectionBookings.UpdateMany(context.TODO(), updateFilter, update)
		if err != nil {
			return fmt.Errorf("error removing game '%s' from bookings: %v", game.Name, err)
		}
	}

	// Удаляем бронирования, связанные с тарифом
	filterBookings := bson.M{"tariff_id": objectTariffId}
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
