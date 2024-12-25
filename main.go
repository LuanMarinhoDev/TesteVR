package main

import (
	"TesteVR/config"
	"TesteVR/queue"
	"TesteVR/routes"
	"log"
	"net/http"
	"time"
)

func main() {
	db, err := config.InitDatabase()
	if err != nil {
		log.Fatalf("Erro ao configurar o banco de dados: %v", err)
	}
	defer db.Close()

	queue.SetupRabbitMQ()
	defer queue.CloseRabbitMQ()

	router := routes.SetupRoutes()

	go func() {
		for {
			log.Println("Tentando consumir transações da fila...")

			transactions, err := queue.ConsumeAllMessages()
			if err != nil {
				log.Println("Erro ao consumir mensagens: ", err)
			} else {
				for _, transaction := range transactions {
					err := config.InsertTransaction(db, transaction)
					if err != nil {
						log.Printf("Erro ao inserir transação no banco: %v", err)
					} else {
						log.Printf("Transação %s inserida com sucesso no banco.", transaction.ID)
					}
				}
			}

			log.Println("Aguardando 5 minutos para a próxima execução...")
			time.Sleep(5 * time.Minute)
		}
	}()

	log.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
