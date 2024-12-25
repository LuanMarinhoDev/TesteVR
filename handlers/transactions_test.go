package handlers_test

import (
	"TesteVR/config"
	"TesteVR/models"
	"log"
	"strings"
	"testing"
	"time"
)

func TestGetTransactionHandler(t *testing.T) {

	db, err := config.InitDatabaseTest()
	if err != nil {
		t.Fatalf("Erro ao inicializar o banco de dados: %v", err)
	}

	_, err = config.GetDB().Query("SELECT 1")
	if err != nil {
		t.Fatalf("Erro ao conectar ao banco de dados de teste: %v", err)
	}

	err = config.CreateTransactionTable(db)
	if err != nil {
		t.Fatalf("Erro ao criar a tabela de transações: %v", err)
	}

	date, err := time.Parse(time.RFC3339, "2024-06-24T12:30:00Z")
	if err != nil {
		t.Fatalf("Erro ao converter data: %v", err)
	}

	transaction := models.Transaction{
		ID:          "1",
		Description: "Teste",
		Amount:      100.0,
		Date:        date,
	}

	err = config.InsertTransaction(db, &transaction)
	if err != nil {
		t.Fatalf("Erro ao inserir transação: %v", err)
	}
	t.Logf("Inserção da transação %v concluída.", transaction.ID)

	var dbTransaction models.Transaction
	var dateString string
	err = config.GetDB().QueryRow("SELECT id, description, date, amount FROM transactions WHERE id = ?", transaction.ID).
		Scan(&dbTransaction.ID, &dbTransaction.Description, &dateString, &dbTransaction.Amount)
	if err != nil {
		t.Fatalf("Erro ao buscar transação no banco: %v", err)
	}

	parsedDate, err := time.Parse("2006-01-02 15:04:05-07:00", dateString)
	if err != nil {
		log.Printf("Erro ao converter data para time.Time: %v", err)
	}

	transaction.Date = parsedDate

	t.Logf("Transação encontrada: %v", transaction)

	dateString = strings.Replace(dateString, " ", "T", 1)

	dbTransaction.Date, err = time.Parse(time.RFC3339, dateString)
	if err != nil {
		t.Fatalf("Erro ao converter a data: %v", err)
	}

	if dbTransaction.ID != transaction.ID {
		t.Fatalf("A transação inserida não corresponde à transação esperada")
	}
}
