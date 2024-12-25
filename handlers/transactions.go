package handlers

import (
	"TesteVR/config"
	"TesteVR/models"
	"TesteVR/queue"
	"TesteVR/services"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Requisição inválida", http.StatusBadRequest)
		return
	}

	if transaction.Description == "" {
		http.Error(w, "Descrição não pode ser vazia", http.StatusBadRequest)
		return
	}

	if transaction.Amount == 0 {
		http.Error(w, "Amount não pode ser vazio", http.StatusBadRequest)
		return
	}

	if len(transaction.Description) > 50 {
		http.Error(w, "Descrição não pode exceder 50 caracteres", http.StatusBadRequest)
		return
	}

	if transaction.Amount <= 0 {
		http.Error(w, "Valor da compra deve ser positivo", http.StatusBadRequest)
		return
	}

	transaction.ID = generateID()

	queue.PublishMessage(transaction)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": transaction.ID})
}

func generateID() string {
	return uuid.New().String()
}

func GetTransactionWithCurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["id"]

	currency := r.URL.Query().Get("currency")
	if currency == "" {
		http.Error(w, "Parâmetro de moeda é obrigatório", http.StatusBadRequest)
		return
	}

	db := config.GetDB()

	transaction, err := services.GetTransactionByID(db, transactionID)
	if err != nil {
		http.Error(w, "Transação não encontrada", http.StatusNotFound)
		return
	}

	convertedAmount, exchangeRate, err := services.ConvertCurrency(transaction.Amount, currency, transaction.Date)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao converter moeda: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":               transaction.ID,
		"description":      transaction.Description,
		"date":             transaction.Date.Format(time.RFC3339),
		"amount_usd":       transaction.Amount,
		"exchange_rate":    exchangeRate,
		"converted_amount": convertedAmount,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
