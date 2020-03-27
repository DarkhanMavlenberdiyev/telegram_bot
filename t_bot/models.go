package t_bot

type CrimeEvents interface {
	GetCrime(id int) (*Crime, error)
	UpdateCrime(id int, crime *Crime) (*Crime, error)
	DeleteCrime(id int) error
	GetAllCrimes() ([]*Crime, error)
	CreateCrime(crime *Crime) (*Crime, error)
}

type UserInfo interface {
	GetUser(id int) (*Users, error)
	UpdateUser(id int, user *Users) (*Users, error)
	DeleteUser(id int) error
	GetAllUser() ([]*Users, error)
	CreateUser(user *Users) (*Users, error)
}

type Users struct {
	ID        int     `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	UserName  string  `json:"user_name"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Image     string  `json:"image"`
	History   string   `json:"history"`
	IsHome     bool   `json:"is_home"`
}

type Crime struct {
	ID           int     `json:"id"`
	LocationName string  `json:"location_name"`
	Longitude    float64 `json:"longitude"`
	Latitude     float64 `json:"latitude"`
	Description  string  `json:"description"`
	Image        string  `json:"image"`
	Date         string  `json:"date"`
}
