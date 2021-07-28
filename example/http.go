package main

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func HTTPMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/sleep", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		d, err := time.ParseDuration(values.Get("d"))
		if err != nil {
			zap.L().Error("cannot parse request parameter for duration", zap.Error(err))
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		time.Sleep(d)
	}))
	mux.HandleFunc("/echo/json", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		if err := json.NewEncoder(rw).Encode(values); err != nil {
			zap.L().Error("cannot encode request query parameters as JSON", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}))
	return mux
}
