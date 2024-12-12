package database

type AdminPageData struct {
	Clients []Client
	Tariffs []TariffTitle
}

type Client struct {
	Name  string `bson:"name"`
	Phone string `bson:"phone_number"`
	Email string `bson:"email"`
}

type TariffTitle struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}

type Tariff struct {
	ID      string   `bson:"_id"`
	Name    string   `bson:"name"`
	Price   float64  `bson:"price"`
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
