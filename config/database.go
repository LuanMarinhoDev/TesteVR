package config

import (
	"TesteVR/models"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDatabase() (*sql.DB, error) {
	var err error
	db, err = sql.Open("sqlite3", "./transactions.db")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = CreateTransactionTable(db)
	if err != nil {
		return nil, err
	}

	log.Println("Banco de dados configurado com sucesso")
	return db, nil
}

func InitDatabaseTest() (*sql.DB, error) {
	var err error
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = CreateTransactionTable(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateTransactionTable(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS transactions (
		id TEXT PRIMARY KEY,
		description TEXT,
		date TEXT,
		amount REAL
	);
	`
	_, err := db.Exec(createTableQuery)
	return err
}

func InsertTransaction(db *sql.DB, transaction *models.Transaction) error {
	insertQuery := `
	INSERT INTO transactions (id, description, date, amount)
	VALUES (?, ?, ?, ?)
	ON CONFLICT(id) DO NOTHING; -- Evita duplicação caso o ID já exista
	`
	_, err := db.Exec(insertQuery, transaction.ID, transaction.Description, transaction.Date, transaction.Amount)
	return err
}

func GetDB() *sql.DB {
	if db == nil {
		log.Fatal("O banco de dados não foi inicializado. Chame config.InitDatabase() antes de usar GetDB().")
	}
	return db
}
