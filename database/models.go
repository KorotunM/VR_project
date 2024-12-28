package database

import "time"

type AdminPageData struct {
	Clients       []Client
	Tariffs       []TariffTitle
	Bookings      []BookingDocument
	Statistic     []TariffStats
	StatisticDays []DailyStats
	GeneralGames  []GeneralGame
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

type AdminFormTariff struct {
	Validation string
	IdTariff   string
	Action     string
	Name       string
	Genre      string
	Platform   string
	Price      int
}

type AdminFormClient struct {
	Action string
	Name   string
	Phone  string
	Email  string
}

type AdminFormBooking struct {
	Action                   string
	ClientName               string
	Clients                  []Client
	TariffName               string
	Tariffs                  []TariffTitle
	GeneralGames             []GeneralGame
	SelectedGeneralGamesName []string
	BookingDate              string
	BookingTime              string
	AvailableTimes           []string
	Validation               string
}

type BookingDocument struct {
	ID           string `bson:"_id,omitempty" json:"id"`    // ID документа (генерируется MongoDB)
	ClientID     string `bson:"client_id" json:"client_id"` // ID клиента (ссылка на другой документ)
	ClientName   string
	TariffID     string `bson:"tariff_id" json:"tariff_id"` // ID тарифа (ссылка на другой документ)
	TariffName   string
	GeneralGames []GeneralGame `bson:"general_games" json:"general_games"`
	BookingDate  time.Time     `bson:"booking_date" json:"booking_date"` // Дата бронирования
	Date         string
	BookingTime  string `bson:"booking_time" json:"booking_time"` // Время бронирования
}
type BookingRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Tariff      string `json:"tariff"`
	BookingDate string `json:"booking_date"`
	BookingTime string `json:"booking_time"`
}

type TariffStats struct {
	TariffName    string
	CurrentProfit int
	BookingsCount int
}

type DailyStats struct {
	Date          string `json:"date"`
	BookingsCount int    `json:"bookings_count"`
}

type GeneralGame struct {
	Id    string `bson:"_id,omitempty"`
	Name  string `bson:"name"`
	Genre string `bson:"genre"`
}

type AdminFormGeneralGame struct {
	Action     string
	Validation string
	Game       GeneralGame
}
