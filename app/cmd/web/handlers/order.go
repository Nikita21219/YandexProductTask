package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"io"
	"log"
	"main/internal/order"
	"net/http"
	"strconv"
)

func getOrders(w http.ResponseWriter, r *http.Request, orderRepo order.Repository) {
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

func pushOrders(w http.ResponseWriter, r *http.Request, orderRepo order.Repository) {
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

func Orders(orderRepo order.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getOrders(w, r, orderRepo)
			return
		} else if r.Method == "POST" {
			pushOrders(w, r, orderRepo)
		}
	}
}

func OrderId(orderRepo order.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func OrderComplete(orderRepo order.Repository, rdb *redis.Client) http.HandlerFunc {
	return IdempotentKeyCheckMiddleware(rdb, func(w http.ResponseWriter, r *http.Request) {

	})
}
