package queue_test

import (
	"TesteVR/models"
	"TesteVR/queue"
	"testing"
	"time"
)

func TestPublishMessage(t *testing.T) {
	queue.SetupRabbitMQ()
	defer queue.CloseRabbitMQ()

	date, err := time.Parse(time.RFC3339, "2024-07-19T00:30:00Z")
	if err != nil {
		t.Fatalf("Erro ao converter data: %v", err)
	}

	transaction := models.Transaction{
		Description: "Teste de transação",
		Date:        date,
		Amount:      100.0,
	}

	err = queue.PublishMessage(transaction)
	if err != nil {
		t.Fatalf("Erro ao publicar mensagem: %v", err)
	}

	consumedTransactions, err := queue.ConsumeAllMessages()
	if err != nil {
		t.Fatalf("Erro ao consumir mensagens: %v", err)
	}

	found := false
	for _, consumedTransaction := range consumedTransactions {
		if consumedTransaction.Description == transaction.Description &&
			consumedTransaction.Amount == transaction.Amount &&
			consumedTransaction.Date.Equal(transaction.Date) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("A transação publicada não foi encontrada na fila")
	}

	if len(consumedTransactions) > 1 {
		t.Errorf("A fila contém mais mensagens do que o esperado. Verifique se outros testes estão limpando a fila corretamente")
	}
}
