package handlers

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"io"
	"log"
	"main/internal/order_complete"
	"net/http"
	"net/url"
	"strconv"
	"time"
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

func IdempotentKeyCheckMiddleware(rdb *redis.Client, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idempKey := r.Header.Get("Idempotency-Key")
		if idempKey == "" {
			log.Println("Idempotency-Key not found in request headers")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		stream, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Error to read request body:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		oc := &order_complete.OrderCompleteDto{}
		err = json.Unmarshal(stream, oc)
		if err != nil {
			log.Println("Error unmarshal data from body request:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		lastTimeParams := rdb.Get(context.Background(), idempKey)
		if lastTimeParams.Err() != nil {
			status := rdb.Set(context.Background(), idempKey, oc, 60*60*time.Second)
			log.Println("Redis set Idempotency-Key status:", status)
			next(w, r)
			return
		}

		b, err := lastTimeParams.Bytes()
		if err != nil {
			log.Println("Error convert data from redis to bytes:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ocLast := order_complete.OrderCompleteDto{}
		err = json.Unmarshal(b, &ocLast)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// If params doesn't change
		if ocLast == *oc {
			log.Printf("Request with key %s has already been processed", idempKey)
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			status := rdb.Set(context.Background(), idempKey, oc, 60*60*time.Second)
			log.Println("Redis set Idempotency-Key status:", status)
			next(w, r)
		}
	}
}
