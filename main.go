package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
	correct = ""
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
	dbuser := t_bot.PostgreUser(user)
	res, _ := db.GetAllCrimes()
	users, _ := dbuser.GetAllUser()
	go func() {
		for true {
			for _, r := range res {
				for _, u := range users {
					distance := t_bot.DistanceBetweenTwoLongLat(r.Latitude, r.Longitude, u.Latitude, u.Longitude) * 1000
					datee, _ := time.Parse(t_bot.LayoutISO, r.Date)
					dat := datee.Format("2006-01-02")
					if distance < 2000.0 && r.IsSend == "false" && t_bot.Current == dat {
						resp, err := http.Get(fmt.Sprintf("https://static-maps.yandex.ru/1.x/?ll=%f,%f&size=450,450&z=15&l=map&pt=%f,%f,home~%f,%f,flag", u.Longitude, u.Latitude, u.Longitude, u.Latitude, r.Longitude, r.Latitude))
						if err != nil {
							fmt.Println(err)
						}
						defer resp.Body.Close()
						out, err := os.Create("filename.png")
						if err != nil {
							fmt.Println(err)
						}
						io.Copy(out, resp.Body)
						defer out.Close()
						sendUser := &tb.User{ID: u.ID}
						distanceString := fmt.Sprintf("%f m", distance)
						b.Send(sendUser, "Location: "+r.LocationName+"\nDescription: "+r.Description+"\nDeistance: "+distanceString)
						photo := &tb.Photo{File: tb.FromDisk("filename.png")}
						b.Send(sendUser, photo)
					}
				}
				db.UpdateCrime(r.ID, &t_bot.Crime{
					ID:           r.ID,
					LocationName: r.LocationName,
					Longitude:    r.Longitude,
					Latitude:     r.Latitude,
					Description:  r.Description,
					Image:        r.Image,
					Date:         r.Date,
					IsSend:       "true",
				})
				res, _ = db.GetAllCrimes()

			}
		}
	}()

	endpoints := t_bot.NewEndpointsFactory(db)

	endpointsUser := t_bot.EndpointsFactoryUser(dbuser)

	b.Handle("/hello", endpoints.Hello(b))
	b.Handle("/start", endpoints.Hello(b))
	b.Handle(&t_bot.ReplyBtn3, endpoints.Input(b))
	b.Handle(&t_bot.ReplyBtn2, endpointsUser.AddHome(b, endpoints))
	b.Handle(&t_bot.ReplyBtn1, endpointsUser.GetHome(b))
	b.Handle(&t_bot.ReplyBtn, endpointsUser.DeleteHome(b))

	b.Start()
	return nil
}
