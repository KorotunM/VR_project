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
