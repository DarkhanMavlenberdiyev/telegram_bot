package t_bot

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	LayoutISO = "2006-01-02"
	HELP      = "⌨️ Home location 🏠 - Permanent Location\n  ➡My location 📍- see your home location\n" +
		"  ➡Add location ➕ - adding home location if there is not exist\n" +
		"  ➡Delete location ✖️ - deleting current home location\n" +
		"  ➡Check location ✔️ - checking your current home location, if are there any crime events\n\n" +
		"⌨ Find crime 🔎 - finding crime events\n" +
		"  ➡My location 🌍 - check your current location\n" +
		"  ➡Send location by map - send location you want\n\n" +
		"⌨ Story 🗃️ - see the story you once looked for\n" +
		"  ➡All stories - see all history\n" +
		"  ➡Clear - clear the history"
)

var (
	current_time = time.Now().Local()
	Current      = current_time.Format("2006-01-02")
	HomeDel      = tb.InlineButton{Text: "Delete location ✖️", Unique: "h1"}
	HomeMy       = tb.InlineButton{Text: "My location 📍", Unique: "h2"}
	HomeAdd      = tb.InlineButton{Text: "Add location ➕", Unique: "h3"}
	ComeBack     = tb.InlineButton{Text: "Back 🔙", Unique: "h4"}
	CheckHome    = tb.InlineButton{Text: "Check location ✔️", Unique: "h5"}

	HistoryAll        = tb.InlineButton{Text: "All stories", Unique: "hi1"}
	HistoryClear      = tb.InlineButton{Text: "Clear", Unique: "hi2"}
	inlineHistoryKeys = [][]tb.InlineButton{
		[]tb.InlineButton{HistoryAll, HistoryClear},
		[]tb.InlineButton{ComeBack},
	}

	ReplyHome = tb.ReplyButton{Text: "Home location 🏠"}
	ReplyBtn3 = tb.ReplyButton{Text: "Find crime 🔎"}
	ReplyHist = tb.ReplyButton{Text: "Story 🗃️"}
	ReplyHelp = tb.ReplyButton{Text: "Help ❓"}

	homeKeys = [][]tb.InlineButton{
		[]tb.InlineButton{HomeMy, HomeAdd},
		[]tb.InlineButton{HomeDel, CheckHome},
		[]tb.InlineButton{ComeBack},
	}

	ReplyKeys = [][]tb.ReplyButton{
		[]tb.ReplyButton{ReplyHome, ReplyBtn3},
		[]tb.ReplyButton{ReplyHist, ReplyHelp},
	}

	Rad1 = tb.InlineButton{Text: "100 m", Unique: "m1"}
	Rad2 = tb.InlineButton{Text: "500 m", Unique: "m2"}
	Rad3 = tb.InlineButton{Text: "1000 m", Unique: "m3"}
	Rad4 = tb.InlineButton{Text: "2000 m", Unique: "m4"}

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
		[]tb.InlineButton{Rad1, Rad2},
		[]tb.InlineButton{Rad3, Rad4},
	}

	radMarkup = &tb.ReplyMarkup{InlineKeyboard: inlineKeys}

	crimeList = make([]*Crime, 0)
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

func (ef *endpointsFactory) Start(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		b.Send(m.Sender, "Choose", &tb.ReplyMarkup{
			InlineKeyboard:      nil,
			ReplyKeyboard:       ReplyKeys,
			ResizeReplyKeyboard: false,
		})

	}
}

func (ef *endpointsFactory) Hello(b *tb.Bot, endUser *endpointsFactory) func(m *tb.Message) {
	return func(m *tb.Message) {

		getUser, _ := endUser.usersInfo.GetUser(m.Sender.ID)
		photo := &tb.Photo{File: tb.FromDisk("images/crime.jpg"), Caption: "Hi, " + m.Sender.FirstName + ". Welcome to Crime bot.\nChoose to continue"}
		b.Send(m.Sender, ">>", &tb.ReplyMarkup{
			ReplyKeyboard: ReplyKeys,
		})
		b.Send(m.Sender, photo)
		user := &Users{
			ID:        m.Sender.ID,
			FirstName: m.Sender.FirstName,
			LastName:  m.Sender.LastName,
			UserName:  m.Sender.Username,
			Longitude: getUser.Longitude,
			Latitude:  getUser.Latitude,
			Image:     getUser.Image,
			History:   getUser.History,
			IsHome:    getUser.IsHome,
		}
		fmt.Print(user)
		endUser.usersInfo.CreateUser(user)

	}
}

func (ef *endpointsFactory) GetHome(b *tb.Bot) func(m *tb.Callback) {
	return func(m *tb.Callback) {
		res, err := ef.usersInfo.GetUser(m.Sender.ID)
		if err != nil {
			b.Respond(m, &tb.CallbackResponse{Text: "Please. Try again", ShowAlert: true})
		}
		if res.IsHome == false {
			b.Respond(m, &tb.CallbackResponse{Text: "Home location is not found", ShowAlert: true})
			return
		} else {
			long := fmt.Sprintf("%f", res.Longitude)
			lat := fmt.Sprintf("%f", res.Latitude)
			photo := &tb.Photo{File: tb.FromDisk(fmt.Sprintf("images/%f.jpg", m.Sender.ID)), Caption: "Your Location \nLongitude: " + long + "\nLatitude: " + lat}
			b.Send(m.Sender, photo)
			b.Send(m.Sender, ">>", &tb.ReplyMarkup{InlineKeyboard: homeKeys})
		}
	}
}
func (ef *endpointsFactory) ListHomeKeys(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		b.Send(m.Sender, ">>", &tb.ReplyMarkup{
			InlineKeyboard: homeKeys,
		})
	}
}

func (ef *endpointsFactory) DeleteHome(b *tb.Bot) func(m *tb.Callback) {
	return func(m *tb.Callback) {
		getUser, _ := ef.usersInfo.GetUser(m.Sender.ID)
		if getUser.IsHome == false {
			b.Respond(m, &tb.CallbackResponse{Text: "Home location is not exist", ShowAlert: true})
			return
		}
		user := &Users{
			ID:        m.Sender.ID,
			FirstName: m.Sender.FirstName,
			LastName:  m.Sender.LastName,
			UserName:  m.Sender.Username,
			Longitude: 0,
			Latitude:  0,
			Image:     "",
			History:   getUser.History,
			IsHome:    false,
		}
		_, err := ef.usersInfo.UpdateUser(m.Sender.ID, user)
		if err != nil {
			b.Respond(m, &tb.CallbackResponse{Text: "Can't delete home location. Try again!", ShowAlert: true})

		} else {
			b.Respond(m, &tb.CallbackResponse{Text: "Home location is successfully deleted"})
		}
	}
}

func (ef *endpointsFactory) AddHome(b *tb.Bot, end *endpointsFactory) func(m *tb.Callback) {
	return func(m *tb.Callback) {
		getUser, _ := ef.usersInfo.GetUser(m.Sender.ID)
		if getUser.IsHome == true {
			b.Respond(m, &tb.CallbackResponse{Text: "Home Location is already exist ", ShowAlert: true})
			return
		}

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
			Image:     fmt.Sprintf("images/%f.jpg", m.Sender.ID),
			History:   getUser.History,
			IsHome:    true,
		}
		_, erro := ef.usersInfo.UpdateUser(m.Sender.ID, user)
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
		if erro != nil {
			fmt.Println(err)
			b.Respond(m, &tb.CallbackResponse{Text: "Can't add home location. Try again!", ShowAlert: true})
		} else {
			//b.Send(m.Sender, "Home location is added")
			b.Respond(m, &tb.CallbackResponse{Text: "Home location is added"})
			b.Send(m.Sender, ">>", &tb.ReplyMarkup{InlineKeyboard: homeKeys})

		}

	}
}

func (ef *endpointsFactory) HomeCheck(b *tb.Bot, end *endpointsFactory) func(c *tb.Callback) {
	return func(c *tb.Callback) {
		crimes, _ := end.crimeEvents.GetAllCrimes()
		user, _ := ef.usersInfo.GetUser(c.Sender.ID)
		//fmt.Println(crimes)
		if len(crimes) >= 0 {
			minDistance := math.MaxFloat64
			resCrime := crimes[0]

			for _, crime := range crimes {
				fmt.Println(crime.Date, Current)
				datee, _ := time.Parse(LayoutISO, crime.Date)
				dat := datee.Format(LayoutISO)
				distance := DistanceBetweenTwoLongLat(user.Latitude, user.Longitude, crime.Latitude, crime.Longitude)
				if distance < minDistance && dat == Current {
					minDistance = distance
					resCrime = crime
				}
			}
			datee, _ := time.Parse(LayoutISO, resCrime.Date)
			dat := datee.Format(LayoutISO)
			fmt.Println(dat)
			if dat == Current {
				resp, err := http.Get(fmt.Sprintf("https://static-maps.yandex.ru/1.x/?ll=%f,%f&size=450,450&z=14&l=map&pt=%f,%f,home~%f,%f,flag", user.Longitude, user.Latitude, user.Longitude, user.Latitude, resCrime.Longitude, resCrime.Latitude))
				if err != nil {
					fmt.Println(err)
				}
				defer resp.Body.Close()
				out, err := os.Create(fmt.Sprintf("images/%v%v.jpg", "homeLoc", c.Sender.ID))
				if err != nil {
					fmt.Println(err)
				}
				io.Copy(out, resp.Body)
				defer out.Close()
				distStr := fmt.Sprintf("%f m", minDistance*1000)
				photo := &tb.Photo{File: tb.FromDisk(fmt.Sprintf("images/%v%v.jpg", "homeLoc", c.Sender.ID)), Caption: "Location: " + resCrime.LocationName + "\nDescription: " + resCrime.Description + "\nDistance: " + distStr}
				b.Send(c.Sender, photo)
				b.Send(c.Sender, ">>", &tb.ReplyMarkup{InlineKeyboard: homeKeys})

			} else {
				b.Respond(c, &tb.CallbackResponse{Text: "There are no crime events"})
			}
		} else {
			b.Respond(c, &tb.CallbackResponse{Text: "There are no crime events"})
		}
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
		b.Send(m.Sender, "Chose radius you want", radMarkup)

		LocLng = float64(res.Lng)
		LocLng = float64(res.Lng)

	}
}
func (ef *endpointsFactory) BackMenu(b *tb.Bot) func(c *tb.Callback) {
	return func(c *tb.Callback) {
		b.Delete(c.Message)
		b.Send(c.Sender, "Choose one", &tb.ReplyMarkup{ReplyKeyboard: ReplyKeys})
	}
}

func (ef *endpointsFactory) GetRad1(b *tb.Bot, endUser *endpointsFactory) func(c *tb.Callback) {
	return func(c *tb.Callback) {
		b.Respond(c, &tb.CallbackResponse{
			CallbackID: c.ID,
			Text:       "Ok",
			ShowAlert:  false,
			URL:        "",
		})
		GetCrime(ef, b, c, 100.0, LocLat, LocLng, endUser)
	}
}

func (ef *endpointsFactory) GetRad2(b *tb.Bot, endUser *endpointsFactory) func(c *tb.Callback) {
	return func(c *tb.Callback) {
		GetCrime(ef, b, c, 500.0, LocLat, LocLng, endUser)
	}
}

func (ef *endpointsFactory) GetRad3(b *tb.Bot, endUser *endpointsFactory) func(c *tb.Callback) {
	return func(c *tb.Callback) {
		GetCrime(ef, b, c, 1000.0, LocLat, LocLng, endUser)
	}
}

func (ef *endpointsFactory) GetRad4(b *tb.Bot, endUser *endpointsFactory) func(c *tb.Callback) {
	return func(c *tb.Callback) {
		GetCrime(ef, b, c, 2000.0, LocLat, LocLng, endUser)
	}
}

func GetCrime(ef *endpointsFactory, b *tb.Bot, m *tb.Callback, r float64, lat float64, lng float64, endUser *endpointsFactory) {
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
		photo := &tb.Photo{File: tb.FromDisk("images/curr.jpg"), Caption: "Location: " + resCrime.LocationName + "\n" + "Description: " + resCrime.Description + "\nDistance: " + fmt.Sprintf("%f m", minDistance*1000)}
		b.Send(m.Sender, photo)
		hist, _ := endUser.usersInfo.GetUser(m.Sender.ID)
		updHist := ""
		if len(hist.History) == 0 {
			updHist = hist.History + " " + strconv.Itoa(resCrime.ID)
		} else {
			updHist = strconv.Itoa(resCrime.ID)
		}
		updUser := &Users{
			ID:        m.Sender.ID,
			FirstName: m.Sender.FirstName,
			LastName:  m.Sender.LastName,
			UserName:  m.Sender.Username,
			Longitude: hist.Longitude,
			Latitude:  hist.Latitude,
			Image:     hist.Image,
			History:   updHist,
			IsHome:    hist.IsHome,
		}
		fmt.Println(updUser)
		_, errr := endUser.usersInfo.UpdateUser(m.Sender.ID, updUser)
		if errr != nil {
			fmt.Println(errr.Error())
		}
	} else {
		b.Send(m.Sender, "There are no crime events in your radius today")
	}
	b.Send(m.Sender, "Choose one", &tb.ReplyMarkup{
		ReplyKeyboard: ReplyKeys,
	})

}

func (ef *endpointsFactory) GetCrimeBySend(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		b.Send(m.Sender, "Send me your geolocation")
		loc := make(chan *tb.Location)
		b.Handle(tb.OnLocation, func(m *tb.Message) {
			loc <- m.Location
		})

		res := <-loc
		b.Send(m.Sender, "Chose radius you want", radMarkup)

		LocLat = float64(res.Lat)
		LocLng = float64(res.Lng)
	}
}

func (ef *endpointsFactory) Help(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		b.Send(m.Sender, HELP)
	}
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

func (ef *endpointsFactory) ToHistory(b *tb.Bot, endUsers *endpointsFactory) func(m *tb.Message) {
	return func(m *tb.Message) {
		crimeList = nil
		b.Send(m.Sender, ">>", &tb.ReplyMarkup{InlineKeyboard: inlineHistoryKeys})
		getUser, _ := endUsers.usersInfo.GetUser(m.Sender.ID)
		lis := strings.Split(getUser.History, " ")

		for _, c := range lis {
			toInt, err := strconv.Atoi(c)
			if err == nil {
				crime, _ := ef.crimeEvents.GetCrime(toInt)
				crimeList = append(crimeList, crime)
			}
		}
	}
}

func (ef *endpointsFactory) GetAllHistory(b *tb.Bot) func(c *tb.Callback) {
	return func(c *tb.Callback) {
		if len(crimeList) == 0 {
			b.Respond(c, &tb.CallbackResponse{Text: "The story is empty"})
		} else {
			res := ""
			for i, c := range crimeList {
				res += fmt.Sprintf("%v) Location: %v; Longitude: %v; Latitude: %v; Description: %v; Date: %v\n", i+1, c.LocationName, c.Longitude, c.Latitude, c.Description, c.Date)
			}
			fmt.Println(res, crimeList)
			b.Edit(c.Message, res)
			b.EditReplyMarkup(c.Message, &tb.ReplyMarkup{InlineKeyboard: inlineHistoryKeys})
		}
	}
}
func (ef *endpointsFactory) ClearHistory(b *tb.Bot) func(c *tb.Callback) {
	return func(c *tb.Callback) {
		user, err := ef.usersInfo.GetUser(c.Sender.ID)
		if err == nil {
			updUser := &Users{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				UserName:  user.UserName,
				Longitude: user.Longitude,
				Latitude:  user.Latitude,
				Image:     user.Image,
				History:   "",
				IsHome:    user.IsHome,
			}
			crimeList = nil
			_, err := ef.usersInfo.UpdateUser(user.ID, updUser)
			if err == nil {
				b.Respond(c, &tb.CallbackResponse{Text: "The story has been cleared"})
				b.Edit(c.Message, ">>")
				b.EditReplyMarkup(c.Message, &tb.ReplyMarkup{InlineKeyboard: inlineHistoryKeys})
			}
		}
	}
}
