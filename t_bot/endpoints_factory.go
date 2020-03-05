package t_bot

import (
	"fmt"
	"io"
	"math"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"sync"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	current_time = time.Now().Local()
	Current      = current_time.Format("2006-01-02")
	ReplyBtn     = tb.ReplyButton{Text: "Delete HOME Location"}
	ReplyBtn1    = tb.ReplyButton{Text: "My HOME location"}
	ReplyBtn2    = tb.ReplyButton{Text: "Add HOME location"}
	ReplyBtn3    = tb.ReplyButton{Text: "Find crime 🔪"}
	ReplyKeys    = [][]tb.ReplyButton{
		[]tb.ReplyButton{ReplyBtn, ReplyBtn3},
		[]tb.ReplyButton{ReplyBtn1, ReplyBtn2},
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
	radius    = ""
	rad       = 0.0
	LayoutISO = "2006-01-02"
)

type Endpoints interface {
	GetCrime() func(m *tb.Message)
}

func NewEndpointsFactory(crimeEvent CrimeEvents) *endpointsFactory {
	return &endpointsFactory{crimeEvents: crimeEvent}
}
func EndpointsFactoryUser(userInfo UserInfo) *endpointsFactory {
	return &endpointsFactory{usersInfo: userInfo}
}

type endpointsFactory struct {
	crimeEvents CrimeEvents
	usersInfo   UserInfo
}

func (ef *endpointsFactory) GetHome(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		res, err := ef.usersInfo.GetUser(m.Sender.ID)
		if err != nil {
			b.Send(m.Sender, "Home location is not found")
		} else {
			long := fmt.Sprintf("%f", res.Longitude)
			lat := fmt.Sprintf("%f", res.Latitude)
			photo := &tb.Photo{File: tb.FromDisk(fmt.Sprintf("%f.jpg", m.Sender.ID))}
			b.Send(m.Sender, "Your Location \nLongitude: "+long+"\nLatitude: "+lat)
			b.Send(m.Sender, photo)
		}
	}
}

func (ef *endpointsFactory) DeleteHome(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		err := ef.usersInfo.DeleteUser(m.Sender.ID)
		if err != nil {
			b.Send(m.Sender, "Home is not exist")
		} else {
			b.Send(m.Sender, "Home is successfully deleted")
		}
	}
}

func (ef *endpointsFactory) AddHome(b *tb.Bot, end *endpointsFactory) func(m *tb.Message) {
	return func(m *tb.Message) {
		b.Send(m.Sender, "Send me your geolocation")
		loc := make(chan *tb.Location)
		b.Handle(tb.OnLocation, func(m *tb.Message) {
			loc <- m.Location
		})

		res := <-loc
		fmt.Println(res)

		user := &Users{
			ID:        m.Sender.ID,
			FirstName: m.Sender.FirstName,
			LastName:  m.Sender.LastName,
			UserName:  m.Sender.Username,
			Longitude: float64(res.Lng),
			Latitude:  float64(res.Lat),
			Image:     fmt.Sprintf("%f.jpg", m.Sender.ID),
		}
		_, err := ef.usersInfo.CreateUser(user)
		resp, err := http.Get(fmt.Sprintf("https://static-maps.yandex.ru/1.x/?ll=%f,%f&size=450,450&z=15&l=map&pt=%f,%f,home", res.Lng, res.Lat, res.Lng, res.Lat))
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(resp)
		defer resp.Body.Close()
		out, err := os.Create(fmt.Sprintf("%f.jpg", m.Sender.ID))
		if err != nil {
			fmt.Println(err)
		}
		io.Copy(out, resp.Body)
		defer out.Close()
		if err != nil {
			fmt.Println(err)
			b.Send(m.Sender, "Home Location is already exist")
		} else {
			b.Send(m.Sender, "Home location is added")
			crimes, _ := end.crimeEvents.GetAllCrimes()
			minDistance := math.MaxFloat64
			resCrime := crimes[0]

			for _, crime := range crimes {
				distance := DistanceBetweenTwoLongLat(float64(res.Lat), float64(res.Lng), crime.Latitude, crime.Longitude)
				fmt.Println(distance, crime)
				if distance < minDistance {
					minDistance = distance
					resCrime = crime
				}
			}
			datee, _ := time.Parse(LayoutISO, resCrime.Date)
			dat := datee.Format(LayoutISO)
			if dat == Current && minDistance < 1 {
				distStr := fmt.Sprintf("%f m", minDistance*1000)
				b.Send(m.Sender, "Location: "+resCrime.LocationName+"\nDescription: "+resCrime.Description+"\nDistance: "+distStr)
				photo := &tb.Photo{File: tb.FromDisk(resCrime.Image)}
				b.Send(m.Sender, photo)
			}
		}

	}
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
			InlineKeyboard:      nil,
			ReplyKeyboard:       ReplyKeys,
			ResizeReplyKeyboard: false,
		})

	}
}

func (ef *endpointsFactory) Input(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		b.Send(m.Sender, "Enter your location", &tb.ReplyMarkup{
			InlineKeyboard: nil,
			ReplyKeyboard:  ReplyKeys2,
		})
		loc := make(chan *tb.Location)

		b.Handle(tb.OnLocation, func(ms *tb.Message) {
			loc <- ms.Location
		})

		res := <-loc
		fmt.Println(res)
		g := &sync.WaitGroup{}
		g.Add(1)
		getCrime(ef, b, m, res.Lat, res.Lng, g)
		g.Wait()

	}
}

func getCrime(ef *endpointsFactory, b *tb.Bot, m *tb.Message, lat float32, lng float32, g *sync.WaitGroup) {
	defer g.Done()
	b.Send(m.Sender, "Enter the radius you want find (m)")
	messages := make(chan string)

	b.Handle(tb.OnText, func(m *tb.Message) {
		messages <- m.Text
	})

	radius = <-messages
	res, cor := check(radius)
	if cor == true {
		rad = res
		radius = ""
	} else if cor == false && radius != "" {
		b.Send(m.Sender, "Incorrect input. Try again...")
		radius = ""
		getCrime(ef, b, m, lat, lng, g)
		return
	}

	crimes, _ := ef.crimeEvents.GetAllCrimes()

	minDistance := math.MaxFloat64
	resCrime := crimes[0]

	for _, crime := range crimes {
		distance := DistanceBetweenTwoLongLat(float64(lat), float64(lng), crime.Latitude, crime.Longitude)
		fmt.Println(distance, crime)
		datee, _ := time.Parse(LayoutISO, crime.Date)
		dat := datee.Format("2006-01-02")
		if distance < minDistance && dat == Current {
			minDistance = distance
			resCrime = crime
		}
	}
	fmt.Println(resCrime)
	if rad == 0.0 || minDistance < rad/1000 {
		// b.Send(m.Sender, photo)
		resp, err := http.Get(fmt.Sprintf("https://static-maps.yandex.ru/1.x/?ll=%f,%f&size=450,450&z=14&l=map&pt=%f,%f,vkbkm~%f,%f,flag", lng, lat, lng, lat, resCrime.Longitude, resCrime.Latitude))
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		out, err := os.Create("curr.jpg")
		if err != nil {
			fmt.Println(err)
		}
		io.Copy(out, resp.Body)
		defer out.Close()
		photo := &tb.Photo{File: tb.FromDisk("curr.jpg")}
		b.Send(m.Sender, "Location: "+resCrime.LocationName+"\n"+"Description: "+resCrime.Description+"\nDistance: "+fmt.Sprintf("%f m", minDistance*1000))
		b.Send(m.Sender, photo)
	} else {
		b.Send(m.Sender, "Crime location not found")
	}
	b.Send(m.Sender, "Choose one", &tb.ReplyMarkup{
		ReplyKeyboard: ReplyKeys,
	})

}

func DistanceBetweenTwoLongLat(lat1 float64, long1 float64, lat2 float64, long2 float64) float64 {
	r := 6371.0090667
	lat1 = lat1 * math.Pi / 180.0
	long1 = long1 * math.Pi / 180.0
	lat2 = lat2 * math.Pi / 180.0
	long2 = long2 * math.Pi / 180.0
	dlon := long1 - long2
	d := math.Acos(math.Sin(lat1)*math.Sin(lat2)+math.Cos(lat1)*math.Cos(lat2)*math.Cos(dlon)) * r
	return d
}

func check(r string) (float64, bool) {
	if res, err := strconv.ParseFloat(r, 64); err == nil {
		return res, true
	}
	return 0.0, false
}
