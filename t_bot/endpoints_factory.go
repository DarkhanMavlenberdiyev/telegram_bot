package t_bot

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gospodinzerkalo/crime_city_api/pb"
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

func NewEndpointsFactory(crimeService pb.CrimeServiceClient, ctx context.Context) *endpointsFactory {
	return &endpointsFactory{crimeService: crimeService, ctx: ctx}
}

type endpointsFactory struct {
	crimeService 	pb.CrimeServiceClient
	ctx 			context.Context
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
		photo := &tb.Photo{File: tb.FromDisk("images/crime.jpg"), Caption: "Hi, " + m.Sender.FirstName + ". Welcome to Crime bot.\nChoose to continue"}
		b.Send(m.Sender, ">>", &tb.ReplyMarkup{
			ReplyKeyboard: ReplyKeys,
		})
		b.Send(m.Sender, photo)
	}
}

func (ef *endpointsFactory) GetHome(b *tb.Bot) func(m *tb.Callback) {
	return func(m *tb.Callback) {
		//res, err := ef.usersInfo.GetUser(m.Sender.ID)
		res, err := ef.crimeService.GetHome(ef.ctx, &pb.GetHomeRequest{Id: int64(m.Sender.ID)})
		if err != nil {
			er := fromGRPCErr(err)
			b.Respond(m, &tb.CallbackResponse{Text: er.Error(), ShowAlert: true})
		}else {
			//long := fmt.Sprintf("%f", res.Home.Longitude)
			//lat := fmt.Sprintf("%f", res.Home.Latitude)
			//photo := &tb.Photo{File: tb.FromDisk(fmt.Sprintf("images/%f.jpg", m.Sender.ID)), Caption: "Your Location \nLongitude: " + long + "\nLatitude: " + lat}
			loc := &tb.Location{
				Lat:       float32(res.Home.GetLatitude()),
				Lng:        float32(res.Home.GetLongitude()),
				LivePeriod: 60,
			}
			b.Send(m.Sender, loc)
			b.Send(m.Sender, fmt.Sprintf("Your Location \nLongitude: %f\nLatitude: %f", res.Home.GetLatitude(), res.Home.GetLongitude()))
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
		_, err := ef.crimeService.DeleteHome(ef.ctx, &pb.DeleteHomeRequest{Id: int64(m.Sender.ID)})
		if err != nil {
			er := fromGRPCErr(err)
			b.Respond(m, &tb.CallbackResponse{Text: "Can't delete home location: " + er.Error(), ShowAlert: true})

		} else {
			b.Respond(m, &tb.CallbackResponse{Text: "Home location is successfully deleted"})
		}
	}
}

func (ef *endpointsFactory) AddHome(b *tb.Bot, end *endpointsFactory) func(m *tb.Callback) {
	return func(m *tb.Callback) {
		getUser, err := ef.crimeService.GetHome(ef.ctx, &pb.GetHomeRequest{Id: int64(m.Sender.ID)})
		//if err != nil {
		//	er := fromGRPCErr(err)
		//	b.Respond(m, &tb.CallbackResponse{Text: er.Error(), ShowAlert: true})
		//	return
		//}
		if getUser != nil {
			b.Respond(m, &tb.CallbackResponse{Text: ErrorAlreadyExist.Error(), ShowAlert: true})
			return
		}

		b.Send(m.Sender, "Send me your geolocation")
		loc := make(chan *tb.Location)
		b.Handle(tb.OnLocation, func(m *tb.Message) {
			loc <- m.Location
		})

		res := <-loc

		user := &pb.Home{
			Id:        int64(m.Sender.ID),
			FirstName: m.Sender.FirstName,
			LastName:  m.Sender.LastName,
			UserName:  m.Sender.Username,
			Longitude: float64(res.Lng),
			Latitude:  float64(res.Lat),
			Image:     fmt.Sprintf("images/%f.jpg", m.Sender.ID),
		}
		_, err = ef.crimeService.CreateHome(ef.ctx, &pb.CreateHomeRequest{Home: user})
		if err != nil {
			err = fromGRPCErr(err)
			b.Respond(m, &tb.CallbackResponse{Text: err.Error(), ShowAlert: true})
			return
		}

		b.Respond(m, &tb.CallbackResponse{Text: "Home location is added"})
		b.Send(m.Sender, ">>", &tb.ReplyMarkup{InlineKeyboard: homeKeys})

	}
}

func (ef *endpointsFactory) HomeCheck(b *tb.Bot, end *endpointsFactory) func(c *tb.Callback) {
	return func(c *tb.Callback) {
		user, err := ef.crimeService.CheckHome(ef.ctx, &pb.CheckHomeRequest{Id: int64(c.Sender.ID)})
		if err != nil {
			er := fromGRPCErr(err)
			b.Respond(c, &tb.CallbackResponse{Text: er.Error(), ShowAlert: true})
		}

		if user != nil {
			body := bytes.NewReader(user.MapImage)
			photo := &tb.Photo{File: tb.FromReader(body), Caption: fmt.Sprintf("Location: %s \nDescription: %s \nDistance: %d", user.LocationName, user.Description, user.Distance)}
			b.Send(c.Sender, photo)
			b.Send(c.Sender, ">>", &tb.ReplyMarkup{InlineKeyboard: homeKeys})
		}else {
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
	//b.Delete(m.Message)
	//crimes, err := ef.crimeService.GetCrimes(ef.ctx, &pb.GetCrimesRequest{})
	//if err != nil {
	//	er := fromGRPCErr(err)
	//	b.Respond(m, &tb.CallbackResponse{Text: er.Error(), ShowAlert: true})
	//	return
	//}
	//minDistance := math.MaxFloat64
	//resCrime := crimes.Crimes[0]
	//
	//for _, crime := range crimes.Crimes {
	//	distance := DistanceBetweenTwoLongLat(lat, lng, crime.Latitude, crime.Longitude)
	//
	//	datee, _ := time.Parse(LayoutISO, crime.Date)
	//	dat := datee.Format("2006-01-02")
	//	if distance < minDistance && dat == Current {
	//		minDistance = distance
	//		resCrime = crime
	//	}
	//}
	//if minDistance < r/1000 {
	//
	//	resp, err := http.Get(fmt.Sprintf("https://static-maps.yandex.ru/1.x/?ll=%f,%f&size=450,450&z=14&l=map&pt=%f,%f,vkbkm~%f,%f,flag", lng, lat, lng, lat, resCrime.Longitude, resCrime.Latitude))
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	defer resp.Body.Close()
	//	out, err := os.Create("images/curr.jpg")
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	io.Copy(out, resp.Body)
	//	defer out.Close()
	//	photo := &tb.Photo{File: tb.FromDisk("images/curr.jpg"), Caption: "Location: " + resCrime.LocationName + "\n" + "Description: " + resCrime.Description + "\nDistance: " + fmt.Sprintf("%f m", minDistance*1000)}
	//	b.Send(m.Sender, photo)
	//	hist, _ := endUser.usersInfo.GetUser(m.Sender.ID)
	//	updHist := ""
	//	if len(hist.History) == 0 {
	//		updHist = hist.History + " " + strconv.Itoa(int(resCrime.Id))
	//	} else {
	//		updHist = strconv.Itoa(int(resCrime.Id))
	//	}
	//	updUser := &Users{
	//		ID:        m.Sender.ID,
	//		FirstName: m.Sender.FirstName,
	//		LastName:  m.Sender.LastName,
	//		UserName:  m.Sender.Username,
	//		Longitude: hist.Longitude,
	//		Latitude:  hist.Latitude,
	//		Image:     hist.Image,
	//		History:   updHist,
	//		IsHome:    hist.IsHome,
	//	}
	//	fmt.Println(updUser)
	//	_, errr := endUser.usersInfo.UpdateUser(m.Sender.ID, updUser)
	//	if errr != nil {
	//		fmt.Println(errr.Error())
	//	}
	//} else {
	//	b.Send(m.Sender, "There are no crime events in your radius today")
	//}
	//b.Send(m.Sender, "Choose one", &tb.ReplyMarkup{
	//	ReplyKeyboard: ReplyKeys,
	//})
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
