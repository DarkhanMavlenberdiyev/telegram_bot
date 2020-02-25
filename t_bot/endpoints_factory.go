package t_bot

import (
	"fmt"
	"math"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	ReplyBtn  = tb.ReplyButton{Text: "Help ❓"}
	ReplyBtn3 = tb.ReplyButton{Text: "Find crime 🔪"}
	ReplyKeys = [][]tb.ReplyButton{
		{ReplyBtn},
		{ReplyBtn3},
	}
	KeyBut = tb.ReplyButton{
		Text:     "My location 🌍",
		Location: true,
	}
	ReplyKeys2 = [][]tb.ReplyButton{
		{KeyBut},
	}

	inlineBtn = tb.InlineButton{
		Unique: "sad_moon",
		Text:   "🌚 Button #2",
	}
	inlineKeys = [][]tb.InlineButton{
		[]tb.InlineButton{inlineBtn},
	}
)

type Endpoints interface {
	GetCrime() func(m *tb.Message)
}

func NewEndpointsFactory(crimeEvent CrimeEvents) *endpointsFactory {
	return &endpointsFactory{crimeEvents: crimeEvent}
}

type endpointsFactory struct {
	crimeEvents CrimeEvents
}

func (ef *endpointsFactory) Hello(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		photo := &tb.Photo{File: tb.FromDisk("crime.jpg")}
		b.Send(m.Sender, "Hi, "+m.Sender.FirstName+". Welcome to Crime bot.\n Choose to continue", &tb.ReplyMarkup{
			ReplyKeyboard: ReplyKeys,
		})
		b.Send(m.Sender, photo)

	}
}
func (ef *endpointsFactory) Start(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		b.Send(m.Sender, "Choose", &tb.ReplyMarkup{
			InlineKeyboard: nil,
			ReplyKeyboard:  ReplyKeys,
		})
	}
}
func (ef *endpointsFactory) Input(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		b.Send(m.Sender, "Enter your location", &tb.ReplyMarkup{
			InlineKeyboard: nil,
			ReplyKeyboard:  ReplyKeys2,
		})
	}
}

func (ef *endpointsFactory) GetCrime(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		long := fmt.Sprintf("%f", m.Location.Lng)
		lat := fmt.Sprintf("%f", m.Location.Lat)
		fmt.Println(long + " " + lat)
		crimes, _ := ef.crimeEvents.GetAllCrimes()

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
	}
}

func (ef *endpointsFactory) ListCrime(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		replyKeys3 := make([][]tb.ReplyButton, 0)
		crimes, _ := ef.crimeEvents.GetAllCrimes()
		for _, crime := range crimes {
			repBtn := tb.ReplyButton{
				Text: crime.LocationName,
			}

			replyKeys3 = append(replyKeys3, []tb.ReplyButton{repBtn})

		}

		b.Send(m.Sender, "Choose one", &tb.ReplyMarkup{
			ReplyKeyboard: replyKeys3,
		})

	}
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
