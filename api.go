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
		err := s.handleGetAccount(w)
		if err != nil {
			return err
		}
		return nil
	}
	if r.Method == "POST" {
		err := s.handleCreateAccount(w, r)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (s *ApiServer) handleGetAccount(w http.ResponseWriter) error {
	accounts, err := s.store.GetAllAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return WriteJSON(w, http.StatusMethodNotAllowed, ApiError{Error: "Method not allowed"})
	}
	req := new(CreateAccountRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	account, err := NewAccount(req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		return err
	}

	err = s.store.InsertAccount(account)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, CreateAccountResponse{Account: *account})
}

func (s *ApiServer) handleWithdraw(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return WriteJSON(w, http.StatusMethodNotAllowed, ApiError{Error: "Method not allowed"})
	}
	req := new(WithdrawRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	err = s.store.SubtractAmount(req.FromID, req.Amount)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, "Withdrawn successfully")
}

func (s *ApiServer) handleDeposit(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return WriteJSON(w, http.StatusMethodNotAllowed, ApiError{Error: "Method not allowed"})
	}
	req := new(DepositRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	err = s.store.AddAmount(req.ToID, req.Amount)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, "Deposited successfully")
}

func (s *ApiServer) Start() error {
	router := mux.NewRouter()

	router.HandleFunc("/accounts", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/withdraw", makeHTTPHandleFunc(s.handleWithdraw))
	router.HandleFunc("/deposit", makeHTTPHandleFunc(s.handleDeposit))

	log.Println("JSON API server running on port: ", s.listenAddr)

	err := http.ListenAndServe(s.listenAddr, router)

	if err != nil {
		return err
	}
	return nil
}
