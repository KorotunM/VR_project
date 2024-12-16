package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminPageData struct {
	Clients []Client
	Tariffs []TariffTitle
}

type Client struct {
	Id    string `bson:"_id,omitempty"`
	Name  string `bson:"name"`
	Phone string `bson:"phone number"`
	Email string `bson:"email"`
}

type TariffTitle struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}

type Tariff struct {
	Id      string   `bson:"_id,omitempty"`
	Name    string   `bson:"name"`
	Price   int      `bson:"price"`
	Games   []Game   `bson:"games"`
	Devices []Device `bson:"devices"`
}

type Game struct {
	Name  string `bson:"name"`
	Genre string `bson:"genre"`
}

type Device struct {
	Name     string `bson:"name"`
	Platform string `bson:"platform"`
}

type AjaxDeleteElementTariff struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type AdminFormTariffData struct {
	Validation string
	IdTariff   string
	Action     string
	Name       string
	Genre      string
	Platform   string
}

type AdminFormTariff struct {
	Validation string
	IdTariff   string
	Action     string
	Name       string
	Price      int
}

type AdminFormClient struct {
	Action string
	Name   string
	Phone  string
	Email  string
}

type BookingDocument struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`          // ID документа (генерируется MongoDB)
	ClientID    primitive.ObjectID `bson:"client_id" json:"client_id"`       // ID клиента (ссылка на другой документ)
	TariffID    primitive.ObjectID `bson:"tariff_id" json:"tariff_id"`       // ID тарифа (ссылка на другой документ)
	BookingDate time.Time          `bson:"booking_date" json:"booking_date"` // Дата бронирования
	BookingTime string             `bson:"booking_time" json:"booking_time"` // Время бронирования
}
