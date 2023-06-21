package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/internal/courier"
	"net/http"
	"net/url"
	"strconv"
)

func getLimitAndOffset(query url.Values) (int, int) {
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
		offset = 0
	}
	limit, err := strconv.Atoi(limits[0])
	if err != nil || limit < 0 {
		limit = 1
	}

	return limit, offset
}

func getCouriers(w http.ResponseWriter, r *http.Request) {
	limit, offset := getLimitAndOffset(r.URL.Query())

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

func couriers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getCouriers(w, r)
		return
	} else {
		pushCouriers(w, r)
	}
}
