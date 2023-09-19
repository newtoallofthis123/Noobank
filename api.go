package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAPIServer(listenAddr string, store Store) *ApiServer {
	return &ApiServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *ApiServer) handleInit(w http.ResponseWriter, r *http.Request) error {
	err := s.store.Init()
	if err != nil {
		return err
	}
	return nil
}

func (s *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		error := s.handleGetAccount(w, r)
		if error != nil {
			return error
		}
		return nil
	}
	if r.Method == "POST" {
		error := s.handleCreateAccount(w, r)
		if error != nil {
			return error
		}
		return nil
	}
	return nil
}

func (s *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.getAllAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateAccountRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	account, err := NewAccount(req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		return err
	}

	err = s.store.insertAccount(account)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, CreateAccountResponse{Account: *account})
}

func (s *ApiServer) Start() error {
	router := mux.NewRouter()

	router.HandleFunc("/accounts", makeHTTPHandleFunc(s.handleAccount))

	log.Println("JSON API server running on port: ", s.listenAddr)

	err := http.ListenAndServe(s.listenAddr, router)

	if err != nil {
		return err
	}
	return nil
}
