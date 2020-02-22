package t_bot

import (
	"fmt"

	"github.com/go-pg/pg"
)

type PostgreConfig struct {
	User     string
	Password string
	Port     string
	Host     string
}

type postgreStore struct {
	db *pg.DB
}

func (p postgreStore) AddCrime(crime *Crime) (*Crime, error) {
	panic("implement me")
}

func (p postgreStore) GetCrime(id int) (*Crime, error) {
	crime := &Crime{ID: id}
	err := p.db.Select(crime)
	if err != nil {
		return nil, err
	}

	fmt.Println(crime.ID)
	return crime, nil
}

func (p postgreStore) UpdateCrime(id int, crime *Crime) (*Crime, error) {
	panic("implement me")
}

func (p postgreStore) DeleteCrime(id int) error {
	panic("implement me")
}

func (p postgreStore) GetAllCrimes() ([]*Crime, error) {
	var crimes []*Crime
	err := p.db.Model(&crimes).Select()

	if err != nil {
		return nil, err
	}
	return crimes, nil
}

func NewPostgreBot(config PostgreConfig) CrimeEvents {
	db := pg.Connect(&pg.Options{
		Addr:     config.Host + ":" + config.Port,
		User:     "postgres",
		Password: config.Password,
	})
	return &postgreStore{db: db}

}
