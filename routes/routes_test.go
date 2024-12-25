package routes_test

import (
	"TesteVR/routes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTransactionWithCurrency(t *testing.T) {
	router := routes.SetupRoutes()

	req, err := http.NewRequest("GET", "/transactions/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

}
