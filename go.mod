module github.com/gospodinzerkalo/crime_city_telegram_bot-Golang-

go 1.16

require (
	github.com/go-pg/pg v8.0.7+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/gospodinzerkalo/crime_city_api v0.0.0-20210516080247-43d564a37902
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/onsi/gomega v1.12.0 // indirect
	github.com/urfave/cli v1.22.5
	github.com/urfave/cli/v2 v2.3.0
	gopkg.in/tucnak/telebot.v2 v2.3.5
	mellium.im/sasl v0.2.1 // indirect
)

replace (
	github.com/gospodinzerkalo/crime_city_api => ../crime_city_api
)
