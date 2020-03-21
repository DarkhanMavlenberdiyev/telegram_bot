package t_bot

import (
	"fmt"
	"io"
	"math"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	current_time = time.Now().Local()
	Current      = current_time.Format("2006-01-02")
	HomeDel     = tb.InlineButton{Text: "Delete HOME Location",Unique:"h1"}
	HomeMy   = tb.InlineButton{Text: "My HOME location",Unique:"h2"}
	HomeAdd    = tb.InlineButton{Text: "Add HOME location",Unique:"h3"}
	ComeBack = tb.InlineButton{Text:"Back",Unique:"h4"}

	ReplyHome 	= tb.ReplyButton{Text:"Home location"}
	ReplyBtn3    = tb.ReplyButton{Text: "Find crime 🔪",}

	homeKeys = [][]tb.InlineButton{
		[]tb.InlineButton{HomeMy,HomeAdd},
		[]tb.InlineButton{HomeDel,ComeBack},

	}

	ReplyKeys    = [][]tb.ReplyButton{
		[]tb.ReplyButton{ReplyHome, ReplyBtn3},
	}

	Rad1 = tb.InlineButton{Text:"100 m",Unique:"m1"}
	Rad2 = tb.InlineButton{Text:"500 m",Unique:"m2"}
	Rad3 = tb.InlineButton{Text:"1000 m",Unique:"m3"}
	Rad4 = tb.InlineButton{Text:"2000 m",Unique:"m4"}

	LocLat = 0.0
	LocLng = 0.0

	KeyBut = tb.ReplyButton{
		Text:     "My location 🌍",
		Location: true,
	}
	ReplyKeys2 = [][]tb.ReplyButton{
		{KeyBut},
	}

	inlineKeys = [][]tb.InlineButton{
		[]tb.InlineButton{Rad1,Rad2},
		[]tb.InlineButton{Rad3,Rad4},
	}

	radMarkup = &tb.ReplyMarkup{InlineKeyboard:      inlineKeys,}

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

func (ef *endpointsFactory) GetHome(b *tb.Bot) func(m *tb.Callback) {
	return func(m *tb.Callback) {
		res, err := ef.usersInfo.GetUser(m.Sender.ID)
		if err != nil {
			//b.Send(m.Sender, "Home location is not found")
			b.Respond(m,&tb.CallbackResponse{
				CallbackID: "",
				Text:       "Home location is not found",
				ShowAlert:  true,
				URL:        "",
			})
		} else {
			long := fmt.Sprintf("%f", res.Longitude)
			lat := fmt.Sprintf("%f", res.Latitude)
			photo := &tb.Photo{File: tb.FromDisk(fmt.Sprintf("images/%f.jpg", m.Sender.ID))}
			b.Send(m.Sender, "Your Location \nLongitude: "+long+"\nLatitude: "+lat)
			b.Send(m.Sender, photo)
			b.Send(m.Sender,">>",&tb.ReplyMarkup{InlineKeyboard:homeKeys})
		}
	}
}
func (ef *endpointsFactory) ListHomeKeys(b *tb.Bot) func(m *tb.Message){
	return func(m *tb.Message) {
		b.Send(m.Sender,">>",&tb.ReplyMarkup{
			InlineKeyboard:      homeKeys,
		})
	}
}

func (ef *endpointsFactory) DeleteHome(b *tb.Bot) func(m *tb.Callback) {
	return func(m *tb.Callback) {
		err := ef.usersInfo.DeleteUser(m.Sender.ID)
		if err != nil {
			b.Respond(m, &tb.CallbackResponse{Text:"Home location is not exist",ShowAlert:true})

		} else {
			b.Respond(m, &tb.CallbackResponse{Text:"Home location is successfully deleted",})
		}
	}
}

func (ef *endpointsFactory) AddHome(b *tb.Bot, end *endpointsFactory) func(m *tb.Callback) {
	return func(m *tb.Callback) {
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
		defer resp.Body.Close()
		out, err := os.Create(fmt.Sprintf("images/%f.jpg", m.Sender.ID))
		if err != nil {
			fmt.Println(err)
		}
		io.Copy(out, resp.Body)
		defer out.Close()
		if err != nil {
			fmt.Println(err)
			//b.Send(m.Sender, "Home Location is already exist")
			b.Respond(m,&tb.CallbackResponse{Text:"Home Location is already exist ",ShowAlert:true})
		} else {
			//b.Send(m.Sender, "Home location is added")
			b.Respond(m,&tb.CallbackResponse{Text:"Home location is added"})
			b.Send(m.Sender,">>",&tb.ReplyMarkup{InlineKeyboard:homeKeys})
			//crimes, _ := end.crimeEvents.GetAllCrimes()
			//if len(crimes)>=0{
			//	minDistance := math.MaxFloat64
			//	resCrime := crimes[0]
			//
			//	for _, crime := range crimes {
			//		distance := DistanceBetweenTwoLongLat(float64(res.Lat), float64(res.Lng), crime.Latitude, crime.Longitude)
			//		fmt.Println(distance, crime)
			//		if distance < minDistance {
			//			minDistance = distance
			//			resCrime = crime
			//		}
			//	}
			//	datee, _ := time.Parse(LayoutISO, resCrime.Date)
			//	dat := datee.Format(LayoutISO)
			//	if dat == Current && minDistance < 1 {
			//		distStr := fmt.Sprintf("%f m", minDistance*1000)
			//		b.Send(m.Sender, "Location: "+resCrime.LocationName+"\nDescription: "+resCrime.Description+"\nDistance: "+distStr)
			//		photo := &tb.Photo{File: tb.FromDisk("images/" + resCrime.Image)}
			//		b.Send(m.Sender, photo)
			//	}
			//}
		}

	}
}

func (ef *endpointsFactory) Hello(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		photo := &tb.Photo{File: tb.FromDisk("images/crime.jpg")}
		b.Send(m.Sender, "Hi, "+m.Sender.FirstName+". Welcome to Crime bot.\nChoose to continue", &tb.ReplyMarkup{
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
		b.Send(m.Sender,"Chose radius you want",radMarkup)

		LocLat = float64(res.Lat)
		LocLng = float64(res.Lng)



	}
}
func (ef *endpointsFactory) BackMenu(b *tb.Bot) func(c *tb.Callback){
	return func(c *tb.Callback) {
		//b.EditReplyMarkup(c.Message,&tb.ReplyMarkup{})
		b.Delete(c.Message)

		b.Send(c.Sender,"Choose one",&tb.ReplyMarkup{ReplyKeyboard:ReplyKeys})
	}
}

func (ef *endpointsFactory) GetRad1(b *tb.Bot) func(c *tb.Callback){
	return func(c *tb.Callback) {
		b.Respond(c,&tb.CallbackResponse{
			CallbackID:c.ID,
			Text:       "Ok",
			ShowAlert:  false,
			URL:        "",
		})

		GetCrime(ef,b,c,100.0,LocLat,LocLng)
	}
}
func (ef *endpointsFactory) GetRad2(b *tb.Bot) func(c *tb.Callback){
	return func(c *tb.Callback) {
		GetCrime(ef,b,c,500.0,LocLat,LocLng)
	}
}
func (ef *endpointsFactory) GetRad3(b *tb.Bot) func(c *tb.Callback){
	return func(c *tb.Callback) {
		GetCrime(ef,b,c,1000.0,LocLat,LocLng)
	}
}
func (ef *endpointsFactory) GetRad4(b *tb.Bot) func(c *tb.Callback){
	return func(c *tb.Callback) {
		GetCrime(ef,b,c,2000.0,LocLat,LocLng)
	}
}
func GetCrime(ef *endpointsFactory,b *tb.Bot,m *tb.Callback, r float64,lat float64,lng float64) {
	b.Delete(m.Message)
	crimes, _ := ef.crimeEvents.GetAllCrimes()
	minDistance := math.MaxFloat64
	resCrime := crimes[0]

	for _, crime := range crimes {
		distance := DistanceBetweenTwoLongLat(lat, lng, crime.Latitude, crime.Longitude)

		datee, _ := time.Parse(LayoutISO, crime.Date)
		dat := datee.Format("2006-01-02")
		if distance < minDistance && dat == Current {
			minDistance = distance
			resCrime = crime
		}
	}
	if minDistance < r/1000 {

		resp, err := http.Get(fmt.Sprintf("https://static-maps.yandex.ru/1.x/?ll=%f,%f&size=450,450&z=14&l=map&pt=%f,%f,vkbkm~%f,%f,flag", lng, lat, lng, lat, resCrime.Longitude, resCrime.Latitude))
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		out, err := os.Create("images/curr.jpg")
		if err != nil {
			fmt.Println(err)
		}
		io.Copy(out, resp.Body)
		defer out.Close()
		photo := &tb.Photo{File: tb.FromDisk("images/curr.jpg")}
		b.Send(m.Sender, "Location: "+resCrime.LocationName+"\n"+"Description: "+resCrime.Description+"\nDistance: "+fmt.Sprintf("%f m", minDistance*1000))
		b.Send(m.Sender, photo)
	} else {
		b.Send(m.Sender, "There are no crime events in your radius today")
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

