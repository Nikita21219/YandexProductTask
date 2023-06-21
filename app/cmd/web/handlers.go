package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"main/internal/courier"
	"main/internal/order"
	"net/http"
	"net/url"
	"strconv"
)

func getLimitAndOffset(query url.Values) (int, int, error) {
	offsets, ok := query["offset"]
	if !ok || len(offsets) != 1 {
		offsets = []string{"0"}
	}

	limits, ok := query["limit"]
	if !ok || len(limits) != 1 {
		limits = []string{"1"}
	}

	offset, err := strconv.Atoi(offsets[0])
	if err != nil || offset < 0 {
		return -1, -1, err
	}
	limit, err := strconv.Atoi(limits[0])
	if err != nil || limit < 0 {
		return -1, -1, err
	}

	return limit, offset, nil
}

func getCouriers(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := getLimitAndOffset(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error parse query string:", err)
		return
	}

	ctx := context.Background()
	couriersFromDb, err := courierRepo.FindByLimitAndOffset(ctx, limit, offset)
	if err != nil {
		log.Println("Error to get couriers from db", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(couriersFromDb)
	if err != nil {
		log.Println("Error marshal data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println("Error write data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := getLimitAndOffset(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error parse query string:", err)
		return
	}

	ctx := context.Background()
	couriersFromDb, err := orderRepo.FindByLimitAndOffset(ctx, limit, offset)
	if err != nil {
		log.Println("Error to get couriers from db", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(couriersFromDb)
	if err != nil {
		log.Println("Error marshal data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println("Error write data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func pushCouriers(w http.ResponseWriter, r *http.Request) {
	var cours []*courier.CourierDto

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &cours)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	for _, cour := range cours {
		ok, _ := cour.Valid()
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	ctx := context.Background()
	err = courierRepo.CreateAll(ctx, cours)
	if err != nil {
		log.Println("Error to create couriers:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func pushOrders(w http.ResponseWriter, r *http.Request) {
	var ordersFromDb []*order.OrderDto

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &ordersFromDb)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	for _, o := range ordersFromDb {
		if !o.Valid() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	ctx := context.Background()
	err = orderRepo.CreateAll(ctx, ordersFromDb)
	if err != nil {
		log.Println("Error to create orders:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func couriers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getCouriers(w, r)
		return
	} else if r.Method == "POST" {
		pushCouriers(w, r)
	}
}

func orders(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getOrders(w, r)
		return
	} else if r.Method == "POST" {
		pushOrders(w, r)
	}
}

func courierId(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error to convert ascii to int:", err)
		return
	}

	ctx := context.Background()
	c, err := courierRepo.FindOne(ctx, id)
	if err != nil {
		log.Println("Error to get courier from db", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(c)
	if err != nil {
		log.Println("Error marshal data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println("Error write data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func orderId(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error to convert ascii to int:", err)
		return
	}

	ctx := context.Background()
	o, err := orderRepo.FindOne(ctx, id)
	if err != nil {
		log.Println("Error to get order from db:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(o)
	if err != nil {
		log.Println("Error marshal data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println("Error write data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
