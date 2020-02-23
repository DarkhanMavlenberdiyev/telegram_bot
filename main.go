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
	flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
		},
	}
)

func main() {

	app := cli.NewApp()
	app.Flags = flags
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
		Token: "1065088890:AAHsp6mSFeTC0mf3sZ5WEi8ODL4ZfxHi1cg",
		// You can also set custom API URL. If field is empty it equals to "https://api.telegram.org"
		URL:    "",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	// geocoder := opencagedata.NewGeocoder("cece8bb38d4a4128a1d97760945dea6c")

	user := t_bot.PostgreConfig{
		User:     "postgres",
		Password: "qwerty123",
		Port:     "8080",
		Host:     "0.0.0.0",
	}
	db := t_bot.NewPostgreBot(user)

	replyBtn := tb.ReplyButton{Text: "/hello"}
	replyBtn3 := tb.ReplyButton{Text: "/input"}
	replyKeys := [][]tb.ReplyButton{
		{replyBtn},
		{replyBtn3},
	}

	//inlineBtn := tb.InlineButton{
	//	Unique: "sad_moon",
	//
	//	Text: "🌚 Button #2",
	//}
	//inlineKeys := [][]tb.InlineButton{
	//	[]tb.InlineButton{inlineBtn},
	//}

	keyBut := tb.ReplyButton{
		Text:     "My location",
		Location: true,
	}
	replyKeys2 := [][]tb.ReplyButton{
		{keyBut},
	}
	if err != nil {
		log.Fatal(err)
		return err
	}
	endpoints := t_bot.NewEndpointsFactory(db)

	b.Handle("/hello", func(m *tb.Message) {
		//b.Send(m.Sender, "Hi, "+m.Sender.FirstName)
		photo := &tb.Photo{File: tb.FromDisk("crime.jpg")}
		b.Send(m.Sender, "Hi, "+m.Sender.FirstName+".Welcome to Crime bot")
		b.Send(m.Sender, photo)

	})

	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, "Choose", &tb.ReplyMarkup{
			InlineKeyboard: nil,
			ReplyKeyboard:  replyKeys,
		})

	})

	b.Handle("/input", func(m *tb.Message) {
		b.Send(m.Sender, "Enter your location", &tb.ReplyMarkup{
			InlineKeyboard: nil,
			ReplyKeyboard:  replyKeys2,
		})
	})
	b.Handle(tb.OnLocation, endpoints.GetCrime(b))

	b.Start()
	return nil
}
