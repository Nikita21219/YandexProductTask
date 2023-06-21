package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
	"log"
	"main/internal/config"
	"main/internal/courier"
	"main/internal/order"
	"main/pkg"
	"net/http"
	"os"
)

var courierRepo courier.Repository
var orderRepo order.Repository
var cfg *config.Config

func init() {
	cfg = LoadConfig()

	ctx := context.Background()
	psqlClient, err := pkg.NewClient(ctx, cfg)

	if err != nil {
		log.Fatalln("Error create db client:", err)
	}

	courierRepo = courier.NewRepo(psqlClient)
	orderRepo = order.NewRepo(psqlClient)
}

func LoadConfig() *config.Config {
	confStream, err := os.ReadFile("./config/app.yaml")
	if err != nil {
		log.Fatalln("Error to open read config file:", err)
	}

	conf := config.NewConfig()
	err = yaml.Unmarshal(confStream, conf)
	if err != nil {
		log.Fatalln("Error to unmarshal data from config file:", err)
	}
	return conf
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/couriers", couriers).Methods("GET", "POST")
	r.HandleFunc("/couriers/{id:[0-9]+}", courierId).Methods("GET")
	r.HandleFunc("/orders", orders).Methods("GET", "POST")
	r.HandleFunc("/orders/{id:[0-9]+}", orderId).Methods("GET")

	http.Handle("/", r)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalln("Error launch web server:", err)
	}
}
