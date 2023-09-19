package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func GetEnv() *Env {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	env := &Env{
		DB:       os.Getenv("DB"),
		Password: os.Getenv("PASSWORD"),
		User:     os.Getenv("USER"),
	}
	return env
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			err := WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
			if err != nil {
				return
			}
		}
	}
}
