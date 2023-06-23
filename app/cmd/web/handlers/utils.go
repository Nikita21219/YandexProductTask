package handlers

import (
	"fmt"
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

func IdempotentKeyCheckMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello from middleware")
		next(w, r)
	}
}
