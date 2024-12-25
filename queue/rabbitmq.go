package queue

import (
	"TesteVR/models"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var channel *amqp.Channel

func SetupRabbitMQ() {
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Erro ao conectar ao RabbitMQ: %v", err)
	}

	channel, err = conn.Channel()
	if err != nil {
		log.Fatalf("Erro ao abrir canal no RabbitMQ: %v", err)
	}

	_, err = channel.QueueDeclare(
		"transactions_queue", // Nome da fila
		true,                 // Durável
		false,                // Auto-delete
		false,                // Exclusiva
		false,                // No-wait
		nil,                  // Argumentos adicionais
	)
	if err != nil {
		log.Fatalf("Erro ao declarar fila: %v", err)
	}

	log.Println("RabbitMQ configurado com sucesso")
}

func PublishMessage(transaction models.Transaction) error {
	message, err := json.Marshal(transaction)
	if err != nil {
		log.Println("Erro ao serializar transação:", err)
		return err
	}

	err = channel.Publish(
		"",                   // Exchange
		"transactions_queue", // Nome da fila
		false,                // Mandatory
		false,                // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		log.Println("Erro ao publicar na fila:", err)
		return err
	}

	return nil
}

func ConsumeAllMessages() ([]*models.Transaction, error) {
	var transactions []*models.Transaction

	for {
		msg, ok, err := channel.Get(
			"transactions_queue", // Nome da fila
			false,                // Auto-acknowledge
		)
		if err != nil {
			log.Printf("Erro ao obter mensagem da fila: %v", err)
			return nil, fmt.Errorf("erro ao obter mensagem da fila")
		}

		if !ok {
			break
		}

		var transaction models.Transaction
		err = json.Unmarshal(msg.Body, &transaction)
		if err != nil {
			log.Printf("Erro ao deserializar a mensagem: %v", err)
			if nackErr := channel.Nack(msg.DeliveryTag, false, false); nackErr != nil {
				log.Printf("Erro ao enviar NACK: %v", nackErr)
			}
			continue
		}

		transactions = append(transactions, &transaction)

		if err := channel.Ack(msg.DeliveryTag, false); err != nil {
			log.Printf("Erro ao enviar ACK: %v", err)
			return nil, err
		}
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("fila vazia")
	}

	return transactions, nil
}

func CloseRabbitMQ() {
	if channel != nil {
		channel.Close()
	}
	if conn != nil {
		conn.Close()
	}
}
