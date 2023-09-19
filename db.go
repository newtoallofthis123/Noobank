package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Store interface {
	Init() error
	getAllAccounts() ([]Account, error)
	insertAccount(a *Account) error
}

type PostgresDB struct {
	db *sql.DB
}

// NewStore TODO: Read Credentials from ENV
func NewStore() (*PostgresDB, error) {
	env := GetEnv()
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=localhost port=5432 sslmode=disable", env.User, env.DB, env.Password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &PostgresDB{db}, nil
}

func (s *PostgresDB) CreateTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS accounts (
		    id SERIAL PRIMARY KEY NOT NULL,
		    first_name TEXT NOT NULL,
		    last_name TEXT,
		    number TEXT NOT NULL UNIQUE,
		    email TEXT,
		    encrypted_password TEXT NOT NULL,
		    balance INT NOT NULL DEFAULT 0,
		    created_at TIMESTAMP
		);
	`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresDB) Init() error {
	return s.CreateTable()
}

func (s *PostgresDB) insertAccount(a *Account) error {
	query := `
		INSERT INTO accounts (first_name, last_name, number, email, encrypted_password, balance, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := s.db.Exec(query, a.FirstName, a.LastName, a.Number, a.Email, a.EncryptedPassword, a.Balance, a.CreatedAt)
	return err
}

func (s *PostgresDB) getAllAccounts() ([]Account, error) {
	query := `
		SELECT * FROM accounts
	`
	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []Account{}

	for rows.Next() {
		var a Account
		err := rows.Scan(&a.ID, &a.FirstName, &a.LastName, &a.Number, &a.Email, &a.EncryptedPassword, &a.Balance, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, nil
}
