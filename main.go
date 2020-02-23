package main

import (
	"log"
	"os"
	"time"

	"./t_bot"
	"github.com/urfave/cli"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
// flags = []cli.Flag{
// 	&cli.StringFlag{
// 		Name:    "config",
// 		Aliases: []string{"c"},
// 	},
// }

)

func main() {

	app := cli.NewApp()
	// app.Flags = flags
	app.Commands = cli.Commands{
		&cli.Command{
			Name:   "start",
			Usage:  "start the bot",
			Action: StartBot,
		},
	}
	app.Run(os.Args)

}

func StartBot(d *cli.Context) error {
	b, err := tb.NewBot(tb.Settings{
		Token:  "1065088890:AAHsp6mSFeTC0mf3sZ5WEi8ODL4ZfxHi1cg",
		URL:    "",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return err
	}

	// geocoder := opencagedata.NewGeocoder("cece8bb38d4a4128a1d97760945dea6c")

	user := t_bot.PostgreConfig{
		User:     "postgres",
		Password: "qwerty123",
		Port:     "8080",
		Host:     "0.0.0.0",
	}
	db := t_bot.NewPostgreBot(user)

	endpoints := t_bot.NewEndpointsFactory(db)

	b.Handle("/hello", endpoints.Hello(b))
	b.Handle("/start", endpoints.Start(b))
	b.Handle("/input", endpoints.Input(b))
	b.Handle(tb.OnLocation, endpoints.GetCrime(b))

	b.Start()
	return nil
}
