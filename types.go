package main

type ApiServer struct {
	listenAddr string
}

type Account struct {
	ID                int    `json:"id"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Number            string `json:"number"`
	Email             string `json:"email"`
	EncryptedPassword string `json:"-"`
	Balance           int    `json:"balance"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}
