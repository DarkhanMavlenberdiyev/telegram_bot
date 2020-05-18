package t_bot

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type EndpointsCurl interface {
	GetCrimeCurl(idParam string) func(w http.ResponseWriter, r *http.Request)
	CreateCrimeCurl() func(w http.ResponseWriter, r *http.Request)
	UpdateCurl(idParam string) func(w http.ResponseWriter, r *http.Request)
	DeleteCurl(idParam string) func(w http.ResponseWriter, r *http.Request)
}

func NewEndpointsFactoryCurl(crimeEvent CrimeEvents) *endpointsFactory {
	return &endpointsFactory{crimeEvents: crimeEvent}
}

type endpointsFactoryCurl struct {
	crimeEvents CrimeEvents
}

func (ef *endpointsFactory) GetCrimeCurl(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			writeResponse(w, http.StatusBadRequest, []byte("Crime ID not found"))
			return
		}
		idd, _ := strconv.Atoi(id)
		res, err := ef.crimeEvents.GetCrime(idd)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, []byte("Error: "+err.Error()))
			return
		}
		data, err := json.Marshal(res)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, []byte("Error: "+err.Error()))
			return
		}
		writeResponse(w, http.StatusOK, data)
	}
}

type Street struct {
	City      string `json:"city"`
	Staddress string `json:"staddress"`
}

func (ef *endpointsFactory) CreateCrimeCurl() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, []byte("Error: "+err.Error()))
			return
		}
		crime := &Crime{}
		if err := json.Unmarshal(data, crime); err != nil {
			writeResponse(w, http.StatusBadRequest, []byte("Error: "+err.Error()))
			return
		}

		getStreet, _ := http.Get(fmt.Sprintf(fmt.Sprintf("https://geocode.xyz/%v,%v?geoit=json&lang=en", crime.Latitude, crime.Longitude)))
		//fmt.Println(getStreet)
		stdata, _ := ioutil.ReadAll(getStreet.Body)
		street := &Street{}
		json.Unmarshal(stdata, street)
		//fmt.Println(street,stdata)
		crime.LocationName = fmt.Sprintf("City: %v; Street: %v", street.City, street.Staddress)
		createMap(crime.Image, crime.Longitude, crime.Latitude)
		result, err := ef.crimeEvents.CreateCrime(crime)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, []byte("Error: "+err.Error()))
			return
		}
		response, err := json.Marshal(result)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, []byte("Error: "+err.Error()))
			return
		}
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

		msg := amqp.Publishing{
			Body: []byte(strconv.Itoa(crime.ID)),
		}

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			msg)
		failOnError(err, "Failed to publish a message")

		writeResponse(w, http.StatusCreated, response)

	}
}
func createMap(name string, long float64, lat float64) {
	lng := fmt.Sprintf("%f", long)
	lt := fmt.Sprintf("%f", lat)
	resp, err := http.Get("https://static-maps.yandex.ru/1.x/?ll=" + lng + "," + lt + "&size=450,450&z=16&l=map&pt=" + lng + "," + lt + ",flag")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(resp)
	defer resp.Body.Close()
	out, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
	}
	io.Copy(out, resp.Body)
	defer out.Close()
}
func (ef *endpointsFactory) UpdateCrimeCurl(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			writeResponse(w, http.StatusBadRequest, []byte("Crime ID not found"))
			return
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, []byte("Error: "+err.Error()))
			return
		}
		crime := &Crime{}
		if err := json.Unmarshal(data, crime); err != nil {
			writeResponse(w, http.StatusBadRequest, []byte("Error: "+err.Error()))
			return
		}
		idd, _ := strconv.Atoi(id)
		res, err := ef.crimeEvents.UpdateCrime(idd, crime)
		if err != nil {
			writeResponse(w, http.StatusBadRequest, []byte("Error: "+err.Error()))
			return
		}
		response, err := json.Marshal(res)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, []byte("Error: "+err.Error()))
			return
		}
		writeResponse(w, http.StatusCreated, response)

	}
}
func (ef *endpointsFactory) DeleteCrimeCurl(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			writeResponse(w, http.StatusBadRequest, []byte("Error: Not Found"))
			return
		}
		idd, _ := strconv.Atoi(id)
		err := ef.crimeEvents.DeleteCrime(idd)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, []byte("Error: "+err.Error()))
			return
		}
		writeResponse(w, http.StatusOK, []byte("Crime is deleted successfully"))

	}
}
func writeResponse(w http.ResponseWriter, status int, msg []byte) {
	w.WriteHeader(status)
	w.Write(msg)
}
func failOnError(err error, msg string) {
	if err != nil {
		fmt.Errorf("%s: %s", msg, err)
	}
}
