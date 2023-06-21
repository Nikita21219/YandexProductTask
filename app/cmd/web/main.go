package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"main/internal/courier"
	"main/pkg"
	"net/http"
)

var courierRepo courier.Repository

func init() {
	ctx := context.Background()
	psqlClient, err := pkg.NewClient(
		ctx,
		"postgres",
		"password",
		"localhost",
		"5433",
		"yandex",
	)

	if err != nil {
		log.Fatalln("Error create db client:", err)
	}

	courierRepo = courier.NewRepo(psqlClient)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/couriers", couriers).Methods("GET", "POST")
	r.HandleFunc("/couriers/{id:[0-9]+}", courierId).Methods("GET")

	http.Handle("/", r)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln("Error launch web server:", err)
	}
}
