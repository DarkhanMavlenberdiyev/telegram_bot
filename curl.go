package main

import (
	"net/http"
	"os"

	"./t_bot"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Commands = cli.Commands{
		&cli.Command{
			Name:   "start",
			Usage:  "start the local server",
			Action: StartServer,
		},
	}
	app.Run(os.Args)

}

func StartServer(c *cli.Context) error {

	user := t_bot.PostgreConfig{
		User:     "postgres",
		Password: "qwerty123",
		Port:     "8080",
		Host:     "0.0.0.0",
	}
	db := t_bot.NewPostgreBot(user)

	endpoints := t_bot.NewEndpointsFactoryCurl(db)

	router := mux.NewRouter()
	router.Methods("GET").Path("/").HandlerFunc(endpoints.GetAllCrime())
	router.Methods("GET").Path("/{id}").HandlerFunc(endpoints.GetCrimeCurl("id"))
	router.Methods("PUT").Path("/{id}").HandlerFunc(endpoints.UpdateCrimeCurl("id"))
	router.Methods("POST").Path("/").HandlerFunc(endpoints.CreateCrimeCurl())
	router.Methods("DELETE").Path("/{id}").HandlerFunc(endpoints.DeleteCrimeCurl("id"))
	http.ListenAndServe("0.0.0.0:8000", router)
	return nil
}
