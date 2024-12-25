package services_test

import (
	"TesteVR/services"
	"testing"
	"time"
)

func TestConvertCurrency(t *testing.T) {
	amount := 100.0
	currency := "Zimbabwe-RTGS"
	dateStr := "2024-02-25"
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		t.Errorf("Erro ao converter a data de registro: %v", err)
	}

	convertedAmount, exchangeRate, err := services.ConvertCurrency(amount, currency, date)
	if err != nil {
		t.Fatalf("Erro na conversão de moeda: %v", err)
	}

	if convertedAmount <= "0" {
		t.Errorf("Valor convertido inválido: %v", convertedAmount)
	}

	if exchangeRate <= 0 {
		t.Errorf("Taxa de câmbio inválida: %v", exchangeRate)
	}
}
