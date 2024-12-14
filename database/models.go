package database

type AdminPageData struct {
	Clients []Client
	Tariffs []TariffTitle
}

type Client struct {
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
