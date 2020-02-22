package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"./t_bot"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	command = ""
)

func main() {
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
		return
	}

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
		command = m.Text
	})

	b.Handle(tb.OnLocation, func(m *tb.Message) {
		long := fmt.Sprintf("%f", m.Location.Lng)
		lat := fmt.Sprintf("%f", m.Location.Lat)
		fmt.Println(long + " " + lat)
		crimes, _ := db.GetAllCrimes()

		minDistance := math.MaxFloat64
		resCrime := crimes[0]

		for _, crime := range crimes {
			distance := distanceBetweenTwoLongLat(float64(m.Location.Lat), float64(m.Location.Lng), crime.Latitude, crime.Longitude)
			fmt.Println(distance, crime)
			if distance < minDistance {
				minDistance = distance
				resCrime = crime
			}
		}
		fmt.Println(resCrime)
		photo := &tb.Photo{File: tb.FromDisk(resCrime.Image)}
		// res, _ := geocoder.Geocode(lat+", "+long, nil)
		b.Send(m.Sender, "Location: "+resCrime.LocationName+"\n"+"Description: "+resCrime.Description)
		b.Send(m.Sender, photo)
	})

	b.Start()
}
func distanceBetweenTwoLongLat(lat1 float64, long1 float64, lat2 float64, long2 float64) float64 {
	r := 6371.0090667
	lat1 = lat1 * math.Pi / 180.0
	long1 = long1 * math.Pi / 180.0
	lat2 = lat2 * math.Pi / 180.0
	long2 = long2 * math.Pi / 180.0
	dlon := long1 - long2
	d := math.Acos(math.Sin(lat1)*math.Sin(lat2)+math.Cos(lat1)*math.Cos(lat2)*math.Cos(dlon)) * r
	return d
}
