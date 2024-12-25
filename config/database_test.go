package config_test

import (
	"TesteVR/config"
	"testing"
)

func TestInitDatabase(t *testing.T) {
	db, err := config.InitDatabase()

	if db == nil {
		t.Fatalf("Falha ao inicializar o banco de dados: db é nil")
	}
	err_ping := db.Ping()
	if err_ping != nil {
		t.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
}
