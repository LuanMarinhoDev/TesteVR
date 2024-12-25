package routes

import (
	"TesteVR/handlers"

	"github.com/gorilla/mux"
)

// @title API de Transações
// @version 1.0
// @description Esta API gerencia transações financeiras
// @host localhost:8080
// @BasePath /api/v1
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// @Summary Cria uma nova transação
	// @Description Cria uma nova transação com as informações fornecidas
	// @Tags transactions
	// @Accept json
	// @Produce json
	// @Param transaction body handlers.Transaction true "Dados da transação"
	// @Success 201 {object} handlers.Transaction "Transação criada com sucesso"
	// @Failure 400 {object} handlers.ErrorResponse "Erro de validação"
	// @Router /transactions [post]
	router.HandleFunc("/transactions", handlers.CreateTransactionHandler).Methods("POST")

	// Recupera uma transação pelo ID
	// @Summary Obtém uma transação por ID
	// @Description Retorna os detalhes de uma transação com base no ID fornecido
	// @Tags transactions
	// @Accept json
	// @Produce json
	// @Param id path string true "ID da transação"
	// @Success 200 {object} handlers.Transaction "Detalhes da transação"
	// @Failure 404 {object} handlers.ErrorResponse "Transação não encontrada"
	// @Router /transactions/{id} [get]
	router.HandleFunc("/transactions/{id}", handlers.GetTransactionWithCurrency).Methods("GET")

	return router
}
