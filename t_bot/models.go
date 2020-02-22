package t_bot

type CrimeEvents interface {
	AddCrime(crime *Crime) (*Crime, error)
	GetCrime(id int) (*Crime, error)
	UpdateCrime(id int, crime *Crime) (*Crime, error)
	DeleteCrime(id int) error
	GetAllCrimes() ([]*Crime, error)
}

type Crime struct {
	ID           int     `json:"id"`
	LocationName string  `json:"location_name"`
	Longitude    float64 `json:"longitude"`
	Latitude     float64 `json:"latitude"`
	Description  string  `json:"description"`
	Image        string  `json:image`
}
