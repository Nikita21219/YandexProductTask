package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"main/internal/courier"
	"net/http"
	"strconv"
)

func getCouriers(w http.ResponseWriter, r *http.Request, courierRepo courier.Repository) {
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

func pushCouriers(w http.ResponseWriter, r *http.Request, courierRepo courier.Repository) {
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

func Couriers(courierRepo courier.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCouriers(w, r, courierRepo)
			return
		} else if r.Method == "POST" {
			pushCouriers(w, r, courierRepo)
		}
	}
}

func CourierId(courierRepo courier.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}
