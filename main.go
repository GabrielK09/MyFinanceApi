package main

import (
	"log"
	"my_finance/internal/database"
	loggerHelper "my_finance/internal/logger"
	r "my_finance/internal/routes"
	"os"
)

func main() {
	wd, _ := os.Getwd()

	log.Println("Rodando em:", wd)
	loggerHelper.Logger() // Start logger

	dbConn, err := database.Init()

	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	defer dbConn.Close()

	log.Println("Banco de dados conectado com sucesso!")
	r.StartServer(dbConn) // Start server
}
