package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"io"
	"net/http"
	"strconv"

	//"io"
	"log"
	//"net/http"
	"os"
	//"strconv"
	"time"

	"./t_bot"
	"github.com/urfave/cli"
	tb "gopkg.in/tucnak/telebot.v2"
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
		User:     "darkhan",
		Password: "sheha2003",
		Port:     "5432", //5432
		Host:     "127.0.0.1",
	}

	db := t_bot.NewPostgreBot(user)
	dbuser := t_bot.PostgreUser(user)
	users, _ := dbuser.GetAllUser()


	go func() {
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		q, err := ch.QueueDeclare(
			"Crime", // name
			true,    // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		failOnError(err, "Failed to declare a queue")

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		failOnError(err, "Failed to register a consumer")

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				log.Printf("Received a message: %s", d.Body)
				conv, _ := strconv.Atoi(string(d.Body))
				crime, _ := db.GetCrime(conv)
				fmt.Println(crime)
				for _, u := range users {
					distance := t_bot.DistanceBetweenTwoLongLat(crime.Latitude, crime.Longitude, u.Latitude, u.Longitude) * 1000
					datee, _ := time.Parse(t_bot.LayoutISO, crime.Date)
					dat := datee.Format("2006-01-02")
					if distance < 2000.0 && t_bot.Current == dat {
						resp, err := http.Get(fmt.Sprintf("https://static-maps.yandex.ru/1.x/?ll=%f,%f&size=450,450&z=15&l=map&pt=%f,%f,home~%f,%f,flag", u.Longitude, u.Latitude, u.Longitude, u.Latitude, crime.Longitude, crime.Latitude))
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
						b.Send(sendUser, "Location: "+crime.LocationName+"\nDescription: "+crime.Description+"\nDeistance: "+distanceString)
						photo := &tb.Photo{File: tb.FromDisk("filename.png")}
						b.Send(sendUser, photo)
					}

				}
			}
		}()

		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
		<-forever
	}()

	endpoints := t_bot.NewEndpointsFactory(db)

	endpointsUser := t_bot.EndpointsFactoryUser(dbuser)


	b.Handle("/start", endpoints.Hello(b))
	b.Handle(&t_bot.ReplyBtn3, endpoints.Input(b))
	b.Handle(&t_bot.HomeAdd, endpointsUser.AddHome(b, endpoints))
	b.Handle(&t_bot.HomeMy, endpointsUser.GetHome(b))
	b.Handle(&t_bot.HomeDel, endpointsUser.DeleteHome(b))
	b.Handle(&t_bot.ReplyHome,endpoints.ListHomeKeys(b))
	b.Handle(&t_bot.ComeBack,endpoints.BackMenu(b))
	b.Handle(&t_bot.Rad1,endpoints.GetRad1(b))
	b.Handle(&t_bot.Rad2,endpoints.GetRad2(b))
	b.Handle(&t_bot.Rad3,endpoints.GetRad3(b))
	b.Handle(&t_bot.Rad4,endpoints.GetRad4(b))


	b.Start()
	return nil
}
func failOnError(err error, msg string) {
	if err != nil {
		fmt.Errorf("%s: %s", msg, err)
	}
}
