package services

import (
	"TesteVR/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

func GetTransactionByID(db *sql.DB, id string) (*models.Transaction, error) {
	var transaction models.Transaction

	query := "SELECT id, description, date, amount FROM transactions WHERE id = ?"

	row := db.QueryRow(query, id)

	var dateStr string

	if err := row.Scan(&transaction.ID, &transaction.Description, &dateStr, &transaction.Amount); err != nil {
		log.Printf("Erro ao executar Scan na consulta: %v", err)
		return nil, fmt.Errorf("erro ao buscar transação: %v", err)
	}

	parsedDate, err := time.Parse("2006-01-02 15:04:05-07:00", dateStr)
	if err != nil {
		log.Printf("Erro ao converter data para time.Time: %v", err)
		return nil, fmt.Errorf("erro ao converter data para time.Time: %v", err)
	}

	transaction.Date = parsedDate

	return &transaction, nil
}

func ConvertCurrency(amount float64, currency string, transactionDate time.Time) (string, float64, error) {
	sixMonthsBefore := transactionDate.AddDate(0, -6, 0)

	url := fmt.Sprintf("https://api.fiscaldata.treasury.gov/services/api/fiscal_service/v1/accounting/od/rates_of_exchange?fields=country_currency_desc,currency,exchange_rate,record_date&filter=record_date:gte:%s&page[size]=1000", sixMonthsBefore.Format("2006-01-02"))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			Currency     string `json:"country_currency_desc"`
			ExchangeRate string `json:"exchange_rate"`
			RecordDate   string `json:"record_date"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", 0, err
	}

	if len(result.Data) == 0 {
		return "", 0, fmt.Errorf("não há dados de taxas de câmbio para a data fornecida: %v", transactionDate)
	}

	for _, rate := range result.Data {
		if rate.Currency == currency {
			recordDate, err := time.Parse("2006-01-02", rate.RecordDate)
			if err != nil {
				return "", 0, fmt.Errorf("erro ao converter a data de registro: %v", err)
			}

			if recordDate.After(sixMonthsBefore) && recordDate.Before(transactionDate.Add(24*time.Hour)) {
				exchangeRate, err := strconv.ParseFloat(rate.ExchangeRate, 64)
				if err != nil {
					return "", 0, fmt.Errorf("erro ao converter taxa de câmbio: %v", err)
				}

				convertedAmount := amount * exchangeRate
				convertedAmount = math.Ceil(convertedAmount*10) / 10

				convertedAmountStr := fmt.Sprintf("%.2f", convertedAmount)

				return convertedAmountStr, exchangeRate, nil
			}
		}
	}

	return "", 0, fmt.Errorf("taxa de câmbio não encontrada ou fora do intervalo de 6 meses para a moeda %s", currency)
}
