package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Store interface {
	Init() error
	GetAllAccounts() ([]Account, error)
	InsertAccount(a *Account) error
	AddAmount(to int, amount int) error
	SubtractAmount(from int, amount int) error
}

type PostgresDB struct {
	Db *sql.DB
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

func (s *PostgresDB) createTable() error {
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
	_, err := s.Db.Exec(query)
	return err
}

func (s *PostgresDB) Init() error {
	return s.createTable()
}

func (s *PostgresDB) InsertAccount(a *Account) error {
	query := `
		INSERT INTO accounts (first_name, last_name, number, email, encrypted_password, balance, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := s.Db.Exec(query, a.FirstName, a.LastName, a.Number, a.Email, a.EncryptedPassword, a.Balance, a.CreatedAt)
	return err
}

func (s *PostgresDB) GetAllAccounts() ([]Account, error) {
	query := `
		SELECT * FROM accounts
	`
	rows, err := s.Db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account

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

func (s *PostgresDB) GetAccountByID(id int) *Account {
	query := `
		SELECT * FROM accounts WHERE number = $1
	`
	row := s.Db.QueryRow(query, id)

	var a Account
	err := row.Scan(&a.ID, &a.FirstName, &a.LastName, &a.Number, &a.Email, &a.EncryptedPassword, &a.Balance, &a.CreatedAt)
	if err != nil {
		return nil
	}
	return &a
}

func (s *PostgresDB) AddAmount(to int, amount int) error {
	toAccount := s.GetAccountByID(to)
	if toAccount == nil {
		return fmt.Errorf("account not found")
	}
	query := `
		UPDATE accounts SET balance = $1 WHERE id = $2
	`
	_, err := s.Db.Exec(query, toAccount.Balance+amount, toAccount.ID)
	return err
}

func (s *PostgresDB) SubtractAmount(from int, amount int) error {
	fromAccount := s.GetAccountByID(from)
	if fromAccount == nil {
		return fmt.Errorf("account not found")
	}
	if fromAccount.Balance < amount {
		return fmt.Errorf("not enough balance")
	}
	query := `
		UPDATE accounts SET balance = $1 WHERE id = $2
	`
	_, err := s.Db.Exec(query, fromAccount.Balance-amount, fromAccount.ID)
	return err
}

func (s *PostgresDB) Transfer(from int, to int, amount int) error {
	err := s.SubtractAmount(from, amount)
	if err != nil {
		return err
	}
	err = s.AddAmount(to, amount)
	if err != nil {
		return err
	}
	return nil
}
