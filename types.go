package main

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"time"
)

type Env struct {
	DB       string
	Password string
	User     string
}

type ApiServer struct {
	listenAddr string
	store      Store
}

type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Number            int64     `json:"number"`
	Email             string    `json:"email"`
	EncryptedPassword string    `json:"-"`
	Balance           int       `json:"balance"`
	CreatedAt         time.Time `json:"created_at"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type CreateAccountResponse struct {
	Account Account `json:"account"`
}

type WithdrawRequest struct {
	FromID int `json:"from_id"`
	Amount int `json:"amount"`
}

type DepositRequest struct {
	ToID   int `json:"to_id"`
	Amount int `json:"amount"`
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func (a *Account) ValidPassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw)) == nil
}

func NewAccount(firstName, lastName, email, password string) (*Account, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		Email:             email,
		EncryptedPassword: string(encryptedPassword),
		Number:            int64(rand.Intn(1000000)),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
